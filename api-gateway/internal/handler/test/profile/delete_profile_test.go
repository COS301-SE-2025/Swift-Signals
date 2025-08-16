package profile

import (
	"net/http"
	"net/http/httptest"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestDeleteProfile_Success() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_MissingUserID() {
	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "DeleteProfile")
}

func (suite *TestSuite) TestDeleteProfile_NotFound() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewNotFoundError("user not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_Forbidden() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewForbiddenError("cannot delete admin user", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "cannot delete admin user")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_InternalError() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ValidationError() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewValidationError("invalid user ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid user ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_UnauthorizedError() {
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ConstraintError() {
	// Test case where user cannot be deleted due to foreign key constraints or business rules
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewAlreadyExistsError("user has dependent resources and cannot be deleted", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Contains(w.Body.String(), "user has dependent resources and cannot be deleted")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_UserWithIntersections() {
	// Test deleting a user who has intersections - should be handled by the service layer
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_NoResponseHeaders() {
	// Test that successful deletes return minimal response with correct status
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())
	// 204 No Content should not have Content-Type header
	suite.Empty(w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_AdminUser() {
	// Test deleting an admin user - should work if service allows it
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_MultipleOperations() {
	// Test multiple delete attempts (should all fail after first success, but testing idempotency)
	userIds := []string{"user1", "user2", "user3"}

	for _, id := range userIds {
		// Create new context for each user
		ctx := suite.ctx
		if id != "test-user-id" {
			// Override the user ID in context for this test
			suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
				Return(nil)
		} else {
			suite.service.On("DeleteProfile", mock.Anything, id).
				Return(nil)
		}

		req := httptest.NewRequest(http.MethodDelete, "/me", nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		suite.handler.DeleteProfile(w, req)

		suite.Equal(http.StatusNoContent, w.Code)
		suite.Empty(w.Body.String())
	}

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteProfile_ServiceTimeout() {
	// Test case where service takes too long to respond
	suite.service.On("DeleteProfile", mock.Anything, "test-user-id").
		Return(errs.NewInternalError("operation timeout", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}
