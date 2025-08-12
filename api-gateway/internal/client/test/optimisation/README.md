# Optimization Client Tests

This directory contains comprehensive unit tests for the optimization client that handles gRPC communication with the optimization service.

## Test Structure

The tests are organized using the testify/suite pattern with the following files:

### Test Files

- **`TestSuite.go`** - Base test suite setup and helper functions
- **`run_optimisation_test.go`** - Tests for RunOptimisation method

## Test Coverage

### RunOptimisation Tests
- ✅ Successful optimization execution
- ✅ Internal server errors
- ✅ Invalid parameter handling
- ✅ Context timeout scenarios
- ✅ Service unavailable errors
- ✅ Zero value parameter testing
- ✅ Large value parameter testing
- ✅ Optimization type mapping validation
- ✅ Resource exhausted scenarios
- ✅ Request structure validation
- ✅ Timeout context verification

## Key Testing Patterns

### Mock Setup
```go
suite.grpcClient.On("RunOptimisation",
    mock.MatchedBy(func(ctx context.Context) bool {
        // Context validation (timeouts)
    }),
    mock.MatchedBy(func(req *optimisationpb.OptimisationParameters) bool {
        // Request parameter validation
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

### Parameter Validation
```go
mock.MatchedBy(func(req *optimisationpb.OptimisationParameters) bool {
    return req.Parameters.Green == expectedGreen &&
        req.Parameters.Yellow == expectedYellow &&
        req.Parameters.Red == expectedRed &&
        req.Parameters.Speed == expectedSpeed &&
        req.Parameters.Seed == expectedSeed
})
```

## Error Scenarios Tested

### gRPC Status Codes
- `codes.Internal` - Internal server errors
- `codes.InvalidArgument` - Invalid optimization parameters
- `codes.Unavailable` - Service unavailable
- `codes.ResourceExhausted` - Optimization queue full
- `context.DeadlineExceeded` - Context timeout

### Parameter Edge Cases
- Zero values in all simulation parameters
- Large values (999999+) in parameters
- Invalid optimization type strings
- Empty/null parameter validation

### Optimization Type Mapping
Tests validation of optimization type string to enum conversion:
- `OPTIMISATION_TYPE_GRIDSEARCH`
- `OPTIMISATION_TYPE_GENETIC_EVALUATION`
- `OPTIMISATION_TYPE_NONE`

## RunOptimisation Method Details

The `RunOptimisation` method:
- Accepts `model.OptimisationParameters` input
- Returns `*optimisationpb.OptimisationParameters` response
- Sets a 5-second context timeout
- Converts model parameters to protobuf format
- Handles gRPC errors with proper error conversion

### Input Parameter Structure
```go
type OptimisationParameters struct {
    OptimisationType     string
    SimulationParameters SimulationParameters
}

type SimulationParameters struct {
    IntersectionType string
    Green            int
    Yellow           int
    Red              int
    Speed            int
    Seed             int
}
```

### Output Response Structure
```protobuf
message OptimisationParameters {
    OptimisationType optimisation_type = 1;
    SimulationParameters parameters = 2;
}

message SimulationParameters {
    IntersectionType intersection_type = 1;
    int32 green = 2;
    int32 yellow = 3;
    int32 red = 4;
    int32 speed = 5;
    int32 seed = 6;
}
```

## Running Tests

```bash
# Run all optimization client tests
go test ./internal/client/test/optimisation/... -v

# Run specific test file
go test ./internal/client/test/optimisation/run_optimisation_test.go -v

# Run with coverage
go test ./internal/client/test/optimisation/... -cover

# Run specific test case
go test ./internal/client/test/optimisation/... -run TestRunOptimisation_Success -v
```

## Dependencies

- `testify/suite` - Test suite framework
- `testify/mock` - Mocking framework
- `google.golang.org/grpc` - gRPC status codes and errors
- `google.golang.org/protobuf` - Protocol buffer utilities

## Mock Generation

The tests use generated mocks for the gRPC client interfaces. Mocks are located in:
- `internal/mocks/grpc_client/MockOptimisationServiceClient.go`

## Notes

### Context Handling
- The `RunOptimisation` method sets a 5-second timeout context
- Tests verify proper timeout propagation
- Context cancellation scenarios are tested

### Parameter Conversion
- Model parameters are converted to protobuf format
- Type mapping is handled with enum value lookups
- Integer parameters are converted from `int` to `int32`

### Error Handling
- gRPC errors are converted using `util.GrpcErrorToErr`
- All error scenarios return `nil` response with proper error
- Error propagation is tested comprehensively

### Performance Considerations
- Optimization operations can be long-running
- Timeout handling is critical for client responsiveness
- Resource exhaustion scenarios are tested for queue management