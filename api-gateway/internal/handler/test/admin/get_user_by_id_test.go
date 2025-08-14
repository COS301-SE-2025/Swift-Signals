package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestGetUserByID_Success() {
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "testuser",
		Email:           "test@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"int1", "int2", "int3"},
	}

	suite.service.On("GetUserByID", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, response.ID)
	suite.Equal(expectedUser.Username, response.Username)
	suite.Equal(expectedUser.Email, response.Email)
	suite.Equal(expectedUser.IsAdmin, response.IsAdmin)
	suite.Equal(expectedUser.IntersectionIDs, response.IntersectionIDs)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_AdminUser() {
	expectedUser := model.User{
		ID:              "admin-user-id",
		Username:        "adminuser",
		Email:           "admin@example.com",
		IsAdmin:         true,
		IntersectionIDs: []string{},
	}

	suite.service.On("GetUserByID", mock.Anything, "admin-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/admin-user-id", nil)
	req.SetPathValue("id", "admin-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedUser.ID, response.ID)
	suite.Equal(expectedUser.Username, response.Username)
	suite.Equal(expectedUser.Email, response.Email)
	suite.True(response.IsAdmin)
	suite.Empty(response.IntersectionIDs)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_EmptyPathValue() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users/", nil)
	req.SetPathValue("id", "")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "User ID is required")

	suite.service.AssertNotCalled(suite.T(), "GetUserByID")
}

func (suite *TestSuite) TestGetUserByID_NotFound() {
	emptyResponse := model.User{}
	suite.service.On("GetUserByID", mock.Anything, "nonexistent-id").
		Return(emptyResponse, errs.NewNotFoundError("user not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users/nonexistent-id", nil)
	req.SetPathValue("id", "nonexistent-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_ServiceForbiddenError() {
	emptyResponse := model.User{}
	suite.service.On("GetUserByID", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewForbiddenError("only admins can access this endpoint", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "only admins can access this endpoint")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_ServiceInternalError() {
	emptyResponse := model.User{}
	suite.service.On("GetUserByID", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_ServiceValidationError() {
	emptyResponse := model.User{}
	suite.service.On("GetUserByID", mock.Anything, "invalid-id-format").
		Return(emptyResponse, errs.NewValidationError("invalid user ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users/invalid-id-format", nil)
	req.SetPathValue("id", "invalid-id-format")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid user ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_ServiceUnauthorizedError() {
	emptyResponse := model.User{}
	suite.service.On("GetUserByID", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users/test-user-id", nil)
	req.SetPathValue("id", "test-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_LongUserID() {
	longId := "very-long-user-id-that-might-cause-issues-with-some-systems-but-should-still-be-handled-properly-by-the-get-endpoint-implementation"

	expectedUser := model.User{
		ID:              longId,
		Username:        "longuserid",
		Email:           "long@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{},
	}

	suite.service.On("GetUserByID", mock.Anything, longId).
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/"+longId, nil)
	req.SetPathValue("id", longId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(longId, response.ID)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_SpecialCharactersInID() {
	specialId := "user-id-with-special-chars-123_456"

	expectedUser := model.User{
		ID:              specialId,
		Username:        "specialuser",
		Email:           "special@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"special1", "special2"},
	}

	suite.service.On("GetUserByID", mock.Anything, specialId).
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/"+specialId, nil)
	req.SetPathValue("id", specialId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(specialId, response.ID)
	suite.Equal("specialuser", response.Username)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_ResponseHeadersSet() {
	expectedUser := model.User{
		ID:       "test-id",
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	suite.service.On("GetUserByID", mock.Anything, "test-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/test-id", nil)
	req.SetPathValue("id", "test-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetUserByID_UserWithManyIntersections() {
	manyIntersections := make([]string, 50)
	for i := 0; i < 50; i++ {
		manyIntersections[i] = fmt.Sprintf("intersection-%d", i+1)
	}

	expectedUser := model.User{
		ID:              "busy-user-id",
		Username:        "busyuser",
		Email:           "busy@example.com",
		IsAdmin:         false,
		IntersectionIDs: manyIntersections,
	}

	suite.service.On("GetUserByID", mock.Anything, "busy-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/busy-user-id", nil)
	req.SetPathValue("id", "busy-user-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetUserByID(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("busy-user-id", response.ID)
	suite.Len(response.IntersectionIDs, 50)
	suite.Equal("intersection-1", response.IntersectionIDs[0])
	suite.Equal("intersection-50", response.IntersectionIDs[49])

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerGetUserByID(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
