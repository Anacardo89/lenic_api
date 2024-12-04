# lenic_api

A gRPC API for [lenic](https://github.com/Anacardo89/lenic)
- it brings most functionality from the social network to a gRPC API you can interact with via Postman

## Setup:
- install [go](https://go.dev/doc/install)
- setup the yaml config files `config`
- run `go mod tidy` to fetch dependencies
- make sure [lenic](https://github.com/Anacardo89/lenic) is running, or at least the DB
- inside `/cmd` run `gp build` to compile, or `go run .` to run with out compiling
- if you built it, run the executable
- you can now send requests to the API via Postman, use the `lenic.proto` file so Postman can get the definition of the service
