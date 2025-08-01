To isolate the api-gateway client/user/ tests:
```bash
go test -v ./internal/client/test/user
```

To get coverage for client/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/client ./internal/client/test/user
go tool cover -html=coverage.out
```
