To isolate the api-gateway service/admin/ tests:
```bash
go test -v ./internal/service/test/admin
```

To get coverage for admin service tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/service ./internal/service/test/admin
go tool cover -html=coverage.out
```
