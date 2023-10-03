#!/usr/bin/env bash

set -x

go get .
go install github.com/swaggo/swag/cmd/swag@latest
(
  cd development || exit
  go run .
)
swag init
go run .
