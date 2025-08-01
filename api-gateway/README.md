To generate mocks/ type `mockery` in /user-service directory in the terminal

To run all of the unit tests with clean output run:
```bash
go test ./internal/client/test/user ./internal/service/test/auth ./internal/handler/test/auth
```

With verbose:
```bash
go test -v ./internal/client/test/user ./internal/service/test/auth ./internal/handler/test/auth 
```

With coverage:
```bash
go test -v -coverprofile=coverage.out -coverpkg=./internal/client/,./internal/service/,./internal/handler/ ./internal/client/test/user ./internal/service/test/auth ./internal/handler/test/auth 
```

