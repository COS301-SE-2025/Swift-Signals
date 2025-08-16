To isolate the api-gateway service/auth/ tests:
```bash
go test -v ./internal/service/test/auth
```

To get coverage for client/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/service ./internal/service/test/auth
go tool cover -html=coverage.out
```
