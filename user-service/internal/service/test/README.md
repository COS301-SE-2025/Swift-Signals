To isolate the user-service service/ tests:
```bash
go test -v ./internal/service/test
```

To get coverage for db/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/service ./internal/service/test
go tool cover -html=coverage.out
```
