To isolate the user-handler handler/ tests:
```bash
go test -v ./internal/handler/test
```

To get coverage for db/ tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/handler ./internal/handler/test
go tool cover -html=coverage.out
```
