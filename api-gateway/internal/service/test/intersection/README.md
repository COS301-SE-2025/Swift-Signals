To isolate the api-gateway service/intersection/ tests:
```bash
go test -v ./internal/service/test/intersection
```

To get coverage for intersection service tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/service ./internal/service/test/intersection
go tool cover -html=coverage.out
```
