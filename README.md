# GCOM

## Running the Project

To run the project for development, run:
`go run main.go`

To build a binary, run:
`go build`

The project is hosted at `localhost:1323`

We highly recommend testing endpoints using the Insomnia Client and will be sharing testing endpoints there as well. To join the organization, message the channel on Discord with the email you signed up for Insomnia with.

## Accessing the Docs
To access the automatically generated documentation for the API,
navigate to the [Swagger Docs](localhost:1323/swagger/index.html)

## Major Dependencies

- [Echo (webserver framework)](https://echo.labstack.com/docs)
- [GORM (ORM database handler)](https://gorm.io/docs/)

## Authorization
TBD

## File Structure

### Models

This is where struct definitions go, using the naming convention `structname_model.go`. Remember to add your models to `migrate.go` as well so
their tables are automatically created.

### Controllers

This is where controllers for the structs go, which handle all the operations which will be performed with that struct (CRUD for example). Use the naming convension `structname_controller.go`

### Responses

This is where we will store reusable response objects to ensure a consistent messaging format. We encourage using
`error_response.go` to send error responses to standardize them

### Util

This is where utility classes go.

### Configs

This is where configurations go and is also where the db code is stored.

### Tests

This is where tests for every model go, using the naming convention `structname_test.go`

## Where to find examples?

The waypoint set of files `waypoint_models.go` and `waypoint_controller.go` are heavily annotated to provide context as to what most lines are doing and why they are there. Other help can be found in the documentation for both major dependencies, linked above.
