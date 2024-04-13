# Go Rest Starter
A Go REST API starter template with authentication, permissions, email, and more 🚀

This starter template provides the following features to build on:
* 👋 Authentication with stateful tokens
* 💂‍♂️ Authorization with roles (by default `admin` and `superadmin`)
* 👀 Middleware for authentication, permission checks, and panic recovery
* 📧 Emails for welcoming new users and resetting passwords
* 🔁 Graceful shutdown that waits for background tasks to finish
* 🧪 Testing package makes it easy to write integration tests (see [TestAuthE2E](https://github.com/jtbergman/go-rest-starter/blob/main/internal/routes/auth/auth_test.go))
* ⏰ Centralized error handling – always return `ServerError` or `ClientError` and send it with `rest.Error(err)`

## Get Started

Do a bulk find and replace for `go-rest-starter.jtbergman.me` and replace it with your desired module name. 

### Env

Create a `.env` file with this format. Both Docker and Make rely on these values. [SMTP](http://mailtrap.io) optional for `local`. 

```env
# Env: local | dev | prod
# 
# local will log email data to STDOUT, SMTP not required
ENV="local"

# Server
PORT=4000

# Docker
#
# PROJECT_NAME is the group name used by Docker
PROJECT_NAME="go-rest-starter"

# Postgres
#
# Used to create DSN for `make run`. This data persists after stop.
DB_NAME="postgres"
DB_USER="postgres"
DB_PASSWORD="password"

# Postgres Tests
#
# Used to create DSN for `make tests`. This data is not persisted after stop.
TEST_DB_NAME="tests"
TEST_DB_USER="tests"
TEST_DB_PASSWORD="password"

# SMTP
#
# These values are provided by SMTP service, can skip while using local
SMTP_HOST=""
SMTP_PORT=25
SMTP_USERNAME=""
SMTP_PASSWORD=""
SMTP_SENDER="Go Rest Starter <no-reply@go-rest-starter.com>"
```

### Make

To run the application, just run `make run`. Alternatively, run `make` to see all the commands.

```
Usage:
  # These automatically start the database and apply migrations
  run                    run API
  tests                  run tests
  tests/short            run tests skipping integration
  tests/cover            run tests with code coverage

  # Manually start and stop the database
  db/start               start the API database
  db/start/tests         start the Tests database
  db/stop                stop the API database
  db/stop/tests          stop the Tests database

  # Connect to the database to inspect with SQL
  sql                    connect to the API database with psql
  sql/tests              connect to the Tests database with psql

  # Manage databae migrations (requires go-migrate)
  mig/new name=$1        create a new database migration
  mig/up                 migrate to a specific version, or apply all migrations
  mig/down               apply all down database migrations
  mig/force version=$1   force the database to a migration version

  # Count the lines of code in your application
  util/loc               lists the total lines of code

  # Build the app, check the version (date+git hash)
  build                  build the API
  version                Output version of current binary
```

## Routes
The supported routes are demonstrated using [HTTPie](https://httpie.io) syntax.

`/v1/auth/register` Create a user with an email and password.

```
http POST localhost:4000/v1/auth/register \
	email="test@example.com" \
	password="password"
```

`/v1/auth/activate` Activate a user using the activation token (see previous logs)

```
http PUT localhost:4000/v1/auth/activate \
	token="<Activation Token (See Server Logs)>"
```

`/v1/auth/login` Login to get an authentication token.

```
http POST localhost:4000/v1/auth/login \
	email="test@example.com" \
	password="password"
```

`/v1/auth/logout` Logout your user (authentication required)

```
http POST localhost:4000/v1/auth/logout \
	"Authorization: Bearer <Authentication Token>
```

`/v1/auth/rest` request and create new passwords

```
# Request Reset
http POST localhost:4000/v1/auth/reset \
	email="test@example.com"

# Rest with token (see server logs)
http PUT localhost:4000/v1/auth/reset \
	token="<Reset Token (See Server Logs)" \
	password="pa55word"
```

`/v1/auth/delete` Delete your account (authentication required)

```
http POST localhost:4000/v1/auth/delete \
	email="test@example.com" \
	password="pa55word" \
	"Authorization: Bearer <New Authentication Token>
```

`/v1/debug/vars` Check server metrics (admin user required)

```bash
# Create a user, activate, and login
$ http localhost:4000/v1/auth/register email="test@example.com" password="password"
$ http PUT localhost:4000/v1/auth/activate token=<Activation Token (See Server Logs)>
$ http POST localhost:4000/v1/auth/login email="test@example.com" password="password"

# You cannot see /v1/debug/vars
$ http localhost:4000/v1/debug/vars "Authorization: Bearer <Login Token>"

# Connect to database, grant admin 
$ make sql
> SELECT * FROM users;
> SELECT * FROM permissions;
> INSERT INTO user_permissions (user_id, permission_id) VALUES (1, 1);

# Now you can
$ http localhost:4000/v1/debug/vars "Authorization: Bearer <Login Token>"
```

## Adding Routes

To define new routes, create a new package in `internal/routes`. 

A route package should specify its dependencies with a struct using interfaces where possible.

```go
// Encapsulates the Application dependencies required by routes
type Auth struct {
	bg     app.Backgrounder
	logger xlogger.Logger
	mailer mailer.Mailer
	rest   *rest.Rest
	tokens tokens.TokensRepository
	users  users.UsersRepository
}
```

Create a `New` function that takes the `App` dependencies type and initializes itself.

```go
func New(app *app.App) *Auth {
	return &Auth{
		bg:     app.BG,
		logger: app.Logger,
		mailer: app.Mailer,
		rest:   app.Rest,
		tokens: app.Models.Tokens,
		users:  app.Models.Users,
	}
}
```

Define a `Route` function with the following signature and register your routes.

```go
func (auth *Auth) Route(mux *http.ServeMux, mw *middleware.Middleware) {
	mux.HandleFunc(ActivateRoute, auth.Activate)

	mux.HandleFunc(DeleteRoute, mw.Authenticated(auth.Delete))

	mux.HandleFunc(LoginRoute, auth.Login)

	mux.HandleFunc(LogoutRoute, mw.Authenticated(auth.Logout))

	mux.HandleFunc(RegisterRoute, auth.Register)

	mux.HandleFunc(ResetRoute, auth.Reset)
}
```

I prefer to use the `switch r.Method` approach when defining my routes.

```go
const ResetRoute = "/v1/auth/reset"

func (app *Auth) Reset(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/reset.html")

	case "POST":
		app.resetPost(w, r)

	case "PUT":
		app.resetPut(w, r)

	default:
		app.rest.MethodNotAllowed(w, r, "GET, POST, PUT")
	}
}
```

## Writing Route Handlers

Route handlers are defined on the dependencies struct (i.e. `Auth`). 

The `Rest` dependency makes it easy to read JSON, write JSON, and handle errors.

```go
func (auth *Auth) registerPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request
	if err := auth.rest.ReadJSON(w, r, "auth.registerPost", &input); err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Create user
	user, err := auth.users.New(input.Email, input.Password)
	if err != nil {
		auth.rest.Error(w, err)
		return
	}

	// Insert user
	if err := auth.users.Insert(user); err != nil {
		if err.Matches(xerrors.ErrUniqueViolation) {
			err.Data = "That email is already taken"
		}
		auth.rest.Error(w, err)
		return
	}
	
	// ...
}
```

## Accessing the Database

To interact with the database, create a new package in `internal/models`

Create a file named `service.go` that defines and implements an interface.

```go
// Defines a mockable interface for user operations
type UsersRepository interface {
	Delete(user *User) (int64, *xerrors.AppError)
	GetByEmail(email string) (*User, *xerrors.AppError)
	GetByToken(plaintext string) (*User, *xerrors.AppError)
	Insert(user *User) *xerrors.AppError
	New(email, plaintext string) (*User, *xerrors.AppError)
	Update(user *User) *xerrors.AppError
}
```

Create a concrete instance that depends on `core.Queryable`. This allows the same service to use transactions and `*sql.DB` without additional code.

```go
// Provides access to User database methods
type Users struct {
	DB core.Queryable
}
```

Provide a `Repository` method to make service initialization consistent.

```go
func Repository(db core.Queryable) UsersRepository {
	return &Users{DB: db}
}
```

Use `core.RowsAffected` and `xerrors.DatabaseError` to simplify error handling.

```go
// Deletes a user
func (m Users) Delete(user *User) (int64, *xerrors.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, user.ID)
	if err != nil {
		return 0, xerrors.DatabaseError(err, "users.Delete")
	}

	return core.RowsAffected(result, "users.Delete")
}
```

## Writing Tests

This template includes helpers for writing integration test. To create an `App` with mocked dependencies, just call `mocks.App()`. You can then easily create a test handler using the `Routes` method from your package. For an example, see `auth_test.go`. Use functions from the `assert` package to easily write integration tests. Example:

```go
func TestDelete(t *testing.T) {
	// Mark an integration test (skipped with make tests/short)
	assert.Integration(t)

	// Easily create dependencies
	app := mocks.App(t)
	handler := authHandler(app)
	credentials := `{"email": "test@example.com", "password": "password"}`

	// Seed – create user, activate user, login user
	assert.Check(t, registerUser(handler, credentials))
	assert.Check(t, activateUser(handler, app))
	token := loginUser(handler, credentials)
	assert.Check(t, len(token) > 0)

	// Auth Required
	assert.RunHandlerTestCase[failures](t, handler, "POST", DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/AuthRequired",
		Body:   credentials,
		Status: http.StatusUnauthorized,
	})

	// User Not Found
	assert.RunHandlerTestCase[failures](t, handler, "POST", DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/UserNotFound",
		Body:   `{"email": "test2@example.com", "password": "password"}`,
		Auth:   token,
		Status: http.StatusNotFound,
	})

	// Credentials Invalid
	assert.RunHandlerTestCase[failures](t, handler, "POST", DeleteRoute, assert.HandlerTestCase[failures]{
		Name:   "Delete/CredentialsInvalid",
		Body:   `{"email": "test@example.com", "password": "pa55word"}`,
		Auth:   token,
		Status: http.StatusUnauthorized,
	})

	// Success
	assert.RunHandlerTestCase[message](t, handler, "POST", DeleteRoute, assert.HandlerTestCase[message]{
		Name:   "Delete/CredentialsInvalid",
		Body:   credentials,
		Auth:   token,
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your account has been deleted")
		},
	})
}
```

If you have a failing test, use the following to inspect server logs from the test
```go
mocks.Logger(app).Begin()

assert.RunHandlerTestCase[message](t, handler, "POST", DeleteRoute, assert.HandlerTestCase[message]{
		Name:   "Delete/CredentialsInvalid",
		Body:   credentials,
		Auth:   token,
		Status: http.StatusOK,
		FN: func(t *testing.T, result message) {
			assert.Equal(t, result.Message, "Your account has been deleted")
		},
	})

mocks.Logger(app).End()
```
