To generate mocks/ type `mockery` in /user-service directory in the terminal

To run all of the unit tests with clean output run:
```bash
go test ./internal/db/test ./internal/service/test ./internal/handler/test
```

With verbose:
```bash
go test -v ./internal/db/test ./internal/service/test ./internal/handler/test
```

With coverage:
```bash
go test -v -coverprofile=coverage.out -coverpkg=./internal/db/,./internal/service/,./internal/handler/ ./internal/db/test ./internal/service/test ./internal/handler/test
```

