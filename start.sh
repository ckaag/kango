#!/bin/sh

sqlc generate
templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
