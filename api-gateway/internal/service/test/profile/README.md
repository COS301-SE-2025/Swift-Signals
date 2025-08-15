To isolate the api-gateway service/profile/ tests:
```bash
go test -v ./internal/service/test/profile
```

To get coverage for profile service tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/service ./internal/service/test/profile
go tool cover -html=coverage.out
```
