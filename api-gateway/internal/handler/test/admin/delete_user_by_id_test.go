package admin

import (
	"net/http"
	"net/http/httptest"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestDeleteUserByID_Success() {
	suite.service.On("DeleteUserByID", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_EmptyPathValue() {
	req := httptest.NewRequest(http.MethodDelete, "/admin/users/", nil)
	req.SetPathValue("id", "")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "User ID is required")

	suite.service.AssertNotCalled(suite.T(), "DeleteUserByID")
}

func (suite *TestSuite) TestDeleteUserByID_NotFound() {
	suite.service.On("DeleteUserByID", mock.Anything, "nonexistent-id").
		Return(errs.NewNotFoundError("user not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/nonexistent-id", nil)
	req.SetPathValue("id", "nonexistent-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_ServiceForbiddenError() {
	suite.service.On("DeleteUserByID", mock.Anything, "test-user-id").
		Return(errs.NewForbiddenError("only admins can access this endpoint", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "only admins can access this endpoint")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_ServiceInternalError() {
	suite.service.On("DeleteUserByID", mock.Anything, "test-user-id").
		Return(errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_ServiceValidationError() {
	suite.service.On("DeleteUserByID", mock.Anything, "invalid-id").
		Return(errs.NewValidationError("invalid user ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid user ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_ServiceUnauthorizedError() {
	suite.service.On("DeleteUserByID", mock.Anything, "test-user-id").
		Return(errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_ConstraintError() {
	// Test case where user cannot be deleted due to foreign key constraints or business rules
	suite.service.On("DeleteUserByID", mock.Anything, "constrained-user-id").
		Return(errs.NewAlreadyExistsError("user has dependent resources and cannot be deleted", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/constrained-user-id", nil)
	req.SetPathValue("id", "constrained-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Contains(w.Body.String(), "user has dependent resources and cannot be deleted")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_LongUserID() {
	longId := "very-long-user-id-that-might-cause-issues-with-some-systems-but-should-still-be-handled-properly-by-the-delete-endpoint-implementation"

	suite.service.On("DeleteUserByID", mock.Anything, longId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+longId, nil)
	req.SetPathValue("id", longId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_SpecialCharactersInID() {
	specialId := "user-id-with-special-chars-123_456"

	suite.service.On("DeleteUserByID", mock.Anything, specialId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+specialId, nil)
	req.SetPathValue("id", specialId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_MultipleSuccessfulDeletes() {
	// Test that multiple delete operations work independently
	userIds := []string{"user1", "user2", "user3"}

	for _, id := range userIds {
		suite.service.On("DeleteUserByID", mock.Anything, id).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+id, nil)
		req.SetPathValue("id", id)
		req = req.WithContext(suite.ctx)
		w := httptest.NewRecorder()

		suite.handler.DeleteUserByID(w, req)

		suite.Equal(http.StatusNoContent, w.Code)
		suite.Empty(w.Body.String())
	}

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_AdminUser() {
	// Test deleting an admin user - should work if service allows it
	suite.service.On("DeleteUserByID", mock.Anything, "admin-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/admin-user-id", nil)
	req.SetPathValue("id", "admin-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_UserWithIntersections() {
	// Test deleting a user who has intersections - should be handled by the service layer
	suite.service.On("DeleteUserByID", mock.Anything, "user-with-intersections").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/user-with-intersections", nil)
	req.SetPathValue("id", "user-with-intersections")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_NumericUserID() {
	numericId := "123456789"

	suite.service.On("DeleteUserByID", mock.Anything, numericId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+numericId, nil)
	req.SetPathValue("id", numericId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_UUIDUserID() {
	uuidId := "550e8400-e29b-41d4-a716-446655440000"

	suite.service.On("DeleteUserByID", mock.Anything, uuidId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+uuidId, nil)
	req.SetPathValue("id", uuidId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteUserByID_NoResponseHeaders() {
	// Test that successful deletes return minimal response with correct status
	suite.service.On("DeleteUserByID", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteUserByID(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())
	// 204 No Content should not have Content-Type header
	suite.Empty(w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}
