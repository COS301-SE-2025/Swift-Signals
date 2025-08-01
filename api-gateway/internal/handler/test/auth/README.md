To isolate the api-gateway handler/auth/ tests:
```bash
go test -v ./internal/handler/test/auth
```

To get coverage for client/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/handler ./internal/handler/test/auth
go tool cover -html=coverage.out
```
