To isolate the api-gateway handler/profile/ tests:
```bash
go test -v ./internal/handler/test/profile
```

To get coverage for profile/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/handler ./internal/handler/test/profile
go tool cover -html=coverage.out
```
