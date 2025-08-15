# User Client Tests

This directory contains comprehensive tests for the UserClient in the Swift-Signals API Gateway.

## Quick Start

To run all user client tests:
```bash
go test -v ./internal/client/test/user
```

To get coverage for client tests:
```bash
go test --coverprofile=coverage.out --coverpkg=./internal/client ./internal/client/test/user
go tool cover -html=coverage.out
```

## Test Coverage Summary

- **Total UserClient Coverage**: ~95% (excluding streaming methods)
- **Total Test Cases**: 80+ comprehensive test scenarios
- **Test Files**: 6 comprehensive test files covering all major functionality

## Test Files

1. **`register_user_test.go`** - User registration functionality
2. **`login_user_test.go`** - Authentication and login tests
3. **`user_crud_test.go`** - CRUD operations (Create, Read, Update, Delete)
4. **`auth_admin_test.go`** - Authentication, password management, and admin operations
5. **`intersection_management_test.go`** - User-intersection relationship management
6. **`constructor_error_test.go`** - Constructor and comprehensive error handling

## What's Tested

### ✅ Fully Covered Methods (100%)
- User registration and authentication
- User CRUD operations (get, update, delete)
- Password management (change, reset)
- Admin operations (make/remove admin)
- Intersection management (add/remove intersections)
- Error handling and validation
- Context timeout management (5-second timeouts)

### ⚠️ Partially Covered (75%)
- `RemoveIntersectionID()` - Single intersection removal

### ❌ Not Tested (Streaming Methods)
- `GetAllUsers()` - Requires streaming client mock
- `GetUserIntersectionIDs()` - Requires streaming client mock
- `NewUserClientFromConn()` - Requires integration testing

## Running Specific Tests

```bash
# Registration tests
go test -v ./internal/client/test/user -run TestClientRegisterUser

# Login tests  
go test -v ./internal/client/test/user -run TestClientLoginUser

# CRUD tests
go test -v ./internal/client/test/user -run TestClientUserCRUD

# Admin operations
go test -v ./internal/client/test/user -run TestClientAuthAdmin

# Intersection management
go test -v ./internal/client/test/user -run TestClientIntersectionManagement

# Error handling
go test -v ./internal/client/test/user -run TestConstructorAndErrorHandling
```

## Error Scenarios Tested

- **Authentication Errors**: Invalid credentials, user not found
- **Validation Errors**: Empty fields, weak passwords, invalid emails
- **Permission Errors**: Insufficient privileges, unauthorized operations
- **System Errors**: Service unavailable, internal errors, rate limiting
- **Context Errors**: Cancellation, timeouts, deadlines

## Key Features

- **Request Validation**: Ensures correct request building
- **Response Handling**: Validates proper response processing
- **Context Management**: Tests timeout and cancellation scenarios
- **Error Propagation**: Verifies gRPC error conversion
- **Mock Integration**: Uses testify/mock for reliable testing

For detailed coverage information, see `COVERAGE_SUMMARY.md`.
