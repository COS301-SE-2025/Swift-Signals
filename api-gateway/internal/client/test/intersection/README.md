# Intersection Client Tests

This directory contains comprehensive unit tests for the intersection client that handles gRPC communication with the intersection service.

## Test Structure

The tests are organized using the testify/suite pattern with the following files:

### Test Files

- **`TestSuite.go`** - Base test suite setup and helper functions
- **`create_intersection_test.go`** - Tests for CreateIntersection method
- **`get_intersection_test.go`** - Tests for GetIntersection method  
- **`get_all_intersections_test.go`** - Tests for GetAllIntersections streaming method
- **`update_intersection_test.go`** - Tests for UpdateIntersection method
- **`delete_intersection_test.go`** - Tests for DeleteIntersection method
- **`put_optimisation_test.go`** - Tests for PutOptimisation method

## Test Coverage

### CreateIntersection Tests
- ✅ Successful intersection creation
- ✅ gRPC error handling
- ✅ Context timeout scenarios
- ✅ String to enum conversions (traffic density, optimization type, intersection type)
- ✅ Empty field handling
- ✅ Request structure validation
- ✅ Timeout context verification

### GetIntersection Tests
- ✅ Successful intersection retrieval
- ✅ Not found errors
- ✅ Empty ID validation
- ✅ Internal server errors
- ✅ Context timeout handling
- ✅ Unauthorized access scenarios
- ✅ Request structure validation

### GetAllIntersections Tests
- ✅ Successful stream creation
- ✅ gRPC error handling
- ✅ Authorization failures
- ✅ Service unavailable scenarios
- ✅ Context cancellation
- ✅ No timeout verification (streaming method)
- ✅ Request structure validation

### UpdateIntersection Tests
- ✅ Successful intersection updates
- ✅ Not found handling
- ✅ Empty ID/name validation
- ✅ Empty details handling
- ✅ Authorization errors
- ✅ Internal server errors
- ✅ Context timeout scenarios
- ✅ Long values testing
- ✅ Special characters handling

### DeleteIntersection Tests
- ✅ Successful intersection deletion
- ✅ Not found errors
- ✅ Empty ID validation
- ✅ Authorization failures
- ✅ Internal server errors
- ✅ Context timeout handling
- ✅ Conflict scenarios (dependencies)
- ✅ Service unavailable errors
- ✅ Long/special character IDs
- ✅ UUID format validation

### PutOptimisation Tests
- ✅ Successful optimization parameter updates
- ✅ Not found handling
- ✅ Empty ID validation
- ✅ Authorization errors
- ✅ Internal server errors
- ✅ String to enum conversions
- ✅ No timeout verification (long-running operation)
- ✅ Zero/large value handling
- ✅ Conflict scenarios
- ✅ Request structure validation

## Key Testing Patterns

### Mock Setup
```go
suite.grpcClient.On("MethodName",
    mock.MatchedBy(func(ctx context.Context) bool {
        // Context validation (timeouts, cancellation)
    }),
    mock.MatchedBy(func(req *RequestType) bool {
        // Request validation
    })).Return(expectedResponse, nil)
```

### Timeout Verification
```go
mock.MatchedBy(func(ctx context.Context) bool {
    deadline, hasDeadline := ctx.Deadline()
    if !hasDeadline {
        return false
    }
    timeUntilDeadline := time.Until(deadline)
    return timeUntilDeadline > 4*time.Second && timeUntilDeadline <= 5*time.Second
})
```

### Error Testing
- gRPC status codes (NotFound, InvalidArgument, Internal, etc.)
- Context timeouts and cancellation
- Network errors and service unavailability
- Authorization and permission errors

### Data Validation
- String to protobuf enum conversions
- Parameter boundary testing (zero, large values)
- Special character and Unicode handling
- Request structure integrity

## Running Tests

```bash
# Run all intersection client tests
go test ./internal/client/test/intersection/... -v

# Run specific test file
go test ./internal/client/test/intersection/create_intersection_test.go -v

# Run with coverage
go test ./internal/client/test/intersection/... -cover

# Run specific test case
go test ./internal/client/test/intersection/... -run TestCreateIntersection_Success -v
```

## Dependencies

- `testify/suite` - Test suite framework
- `testify/mock` - Mocking framework
- `google.golang.org/grpc` - gRPC status codes and errors
- `google.golang.org/protobuf` - Protocol buffer utilities

## Mock Generation

The tests use generated mocks for the gRPC client interfaces. Mocks are located in:
- `internal/mocks/grpc_client/MockIntersectionServiceClient.go`

## Notes

- All methods except `GetAllIntersections` and `PutOptimisation` set a 5-second timeout
- Streaming methods (`GetAllIntersections`) don't set timeouts to allow for long-running operations
- String to enum conversions have default fallback values for invalid inputs
- Tests validate both successful operations and comprehensive error scenarios