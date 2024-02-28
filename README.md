# JWT auth service

Authorization service working through http

# Table of contents

- [Startup](#startup)
    - [Direct startup](#direct-startup)
    - [Docker startup](#docker-startup)
    - [Prepare env](#prepare-env)
- [Http request examples](#http-request-examples)
    - [Register](#Register)
    - [Create token pair](#create-token-pair)
    - [Refresh token](#refresh-token)

## Startup

You can run the microservice through docker or directly.

### Direct startup

You can use `make` command to get bin file or by yourself compile service by using:</br>
`CGO_ENABLED=0 GOOS=linux go build -o <output path> <path to main.go>`

Than you need to start file via command:</br>
`file -config=<path_to_config>`
Examples:</br>
`./serverMain -config=./configs/config.yaml`</br>
`go run ./cmd/server/main.go -config=./configs/config.yaml`

### Docker startup

### Prepare env
If you will use docker-startup firstly you need to convert config to .env file.
There are two ways to do that:
- Use `make prepareEnv`
- Use `go run config_to_env.go <path_to_config>`

</br>
Then just run docker-compose: `docker-compose up`

## Http request examples

### Register

Request

```curl
curl --location '0.0.0.0:9999/api/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "new_test@test.test"
}'
```

Response

```json
{
    "uuid": "1c5d19c0-79e1-4cb2-8fda-e02f3ce20554"
}
```

### Create token pair

Request

```curl
curl --location '0.0.0.0:9999/api/createTokenPair/1c5d19c0-79e1-4cb2-8fda-e02f3ce20554'
```

Response

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5ld190ZXN0QHRlc3QudGVzdCIsImV4cCI6MTcwOTExNDQ5OSwidXVpZCI6IjFjNWQxOWMwLTc5ZTEtNGNiMi04ZmRhLWUwMmYzY2UyMDU1NCJ9.wKyAothxPjJH-nAvGqDa0VYq_uO6QXBc34GD2fn-ytbP57sNrg6ZUl7iTz4oXf7I5UAF0j5mis7ecK5NIv89bA",
    "refresh_token": "ODE3YTlmNTMzM2ZiYzNmYTQ3ZGU3OWU5NDZmOWI4MTUxYzU0ODhiZWE3ZThhY2ZlNzJhYmRlNDMzMjg5ZTRlNw=="
}
```

### Refresh token

Request

```curl
curl --location '0.0.0.0:9999/api/refreshToken' \
--header 'Content-Type: application/json' \
--data '{
    "token_pair": {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5ld190ZXN0QHRlc3QudGVzdCIsImV4cCI6MTcwOTExNDQ5OSwidXVpZCI6IjFjNWQxOWMwLTc5ZTEtNGNiMi04ZmRhLWUwMmYzY2UyMDU1NCJ9.wKyAothxPjJH-nAvGqDa0VYq_uO6QXBc34GD2fn-ytbP57sNrg6ZUl7iTz4oXf7I5UAF0j5mis7ecK5NIv89bA",
        "refresh_token": "ODE3YTlmNTMzM2ZiYzNmYTQ3ZGU3OWU5NDZmOWI4MTUxYzU0ODhiZWE3ZThhY2ZlNzJhYmRlNDMzMjg5ZTRlNw=="
    }
}'
```

Response

```json
{
    "token_pair": {
        "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im5ld190ZXN0QHRlc3QudGVzdCIsImV4cCI6MTcwOTExNDU0OSwidXVpZCI6IjFjNWQxOWMwLTc5ZTEtNGNiMi04ZmRhLWUwMmYzY2UyMDU1NCJ9.RU6Bznxp93lEE3fWO9bGjG5pCemrVg3QetyZUD_lvga5IgfihgqnLEQndaP_lH_n5gjUzZk5cpPwxifO326x8w",
        "refresh_token": "YWM4ZmIxZDlkZDBlNzkxZTljYjk1ODkxZGZkZWIxOWI5MWJkYWQ5YWI3NzM3ZTczZDlhZWVlODI4ZmQ0OWZlYw=="
    }
}
```