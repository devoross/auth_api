# AUTH API 

This is the authentication layer which makes up a small part of a larger application.

## Prerequisites

Required applications:

* Docker
* Docker compose, if you don't already have docker compose see [here](https://docs.docker.com/compose/install/)

The above is suited to run Redis and Jaeger, which are both useful for the execution of the application. Redis is our storage layer (temp, as this will be a caching layer to a postgres DB), which is required and Jaeger is used to log our traces to so that we can visualise them. See screenshot below.

**NOTE: Make sure you have added yourself to the docker group so you can run docker as non-sudo**

## Starting the application ins Docker

Build the docker image for auth_api application

```bash
./build.sh
```

Once the build is complete, you should now have a docker image locally installed and ready to run, therefore you can now run the below command.

**Sometimes the command may be `docker-compose`**

```bash
docker compose up -d
```

What the above command does, is start all the containers configured in the docker-compose yml, which should be everything required to get the application up and running locally.

## Run the application locally

Download dependencies

```bash
go mod download
```

Run from source

```bash
go run main.go
```

Build from source

```bash
go build -o auth_api
```

Run the built binary

```bash
./auth_api
```

The default port is **8080**

## Sessions

We have implemented sessions into the application so that logins can persist for a period of time to allow the user to stay logged in across our applications using a single session ID. 

The Session ID as seen below is returned when the user has successfully logged into their account, this can then be stored in a cookie at the frontend to be provided to all subsequent requests to the website, providing immediate access.

Example Session ID:

```json
{
    "sid": "6a92daba-09c0-4c20-9edf-9d537dce7409"
}
```

TODO : HOW DO WE AUTHENTICATE REQUESTS TO THE API, DO WE EVALUATE THE EXISTENCE OF THIS SESSION ID AT THE SERVER SIDE, OR GENERATE AND STORE AN API KEY SPECIFIC TO THAT SESSION ID?

# Endpoints

**NOTE: Ensure credentials are sent over HTTPs only**

## Register

The register endpoint is responsible for registering users, and creating them a session from that point onwards. **Passwords are stored in the database hashed**

### Request

```bash
curl -XPOST -H "Content-Type: application/json" -d '{"email": "test@test.com", "username": "test_username", "password": "test_password"}' http://localhost:8080/api/auth/register
```

### Response

If registration was successful, you'll receive the below response, and this will have created the user a session with that ID.

```json
{
    "sid": "6a92daba-09c0-4c20-9edf-9d537dce7409"
}
```

## Login with username and password

This will use the passed in password, hash it and compare it with the stored value. If the login is successful, all meta information associated with that user will be returned as seen in the response below

### Request

```bash
curl -XPOST -H "Content-Type: application/json" -d '{"username": "test_username", "password": "test_password"}' http://localhost:8080/api/auth/login
```

### Response

The session id, for subsequent requests, the email that belongs directly to that person and a username

```json
{
    "sid": "f08ff456-3926-4dbf-935e-2b9dee3ffd32",
    "email":"test@test.com",
    "username":"test_username"
}
```

## Login with session

The session ID can be used to log the user in, but without passing credentials around and having the user re-enter them on a frontend. **Worth noting that expiries are not currently supported on sessions, and they last indefinitely**

### Request

```bash
curl -XPOST -H 'Content-Type: application/json' -H "x-session-id: 7fc2639b-2eef-4a15-96ad-1a3cefbc6bf7" http://localhost:8080/api/auth/login  -v
```

### Response

```json
{
    "sid": "f08ff456-3926-4dbf-935e-2b9dee3ffd32",
    "email":"test@test.com",
    "username":"test_username"
}
```

The idea is, that if say the session no longer existed, and couldn't be found, the API would return a 401 allowing the frontend to respond by showing the login page
