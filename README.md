# Auth

Ths is a passwordless auth system written in the style of clean architecture.

# example app:

1. `brew install sqlite`
0. `go run ./cmd/auth/main.go`
0. In another terminal
0. `curl --data email="myemail@example.com http://localhost:8888/login`
0. Look back at the first terminal and copy the URL
0. `curl -i '<paste URL between single quotes>'`
0. Inspect headers
