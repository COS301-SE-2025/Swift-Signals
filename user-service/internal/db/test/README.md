To isolate the user-service db/ tests:
```bash
go test -v ./internal/db/test`
```

To get coverage for db/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/db ./internal/db/test
go tool cover -html=coverage.out
```
