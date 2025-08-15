package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestGetProfile_Success() {
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "testuser",
		Email:           "test@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"int1", "int2", "int3"},
	}

	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

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

func (suite *TestSuite) TestGetProfile_AdminUser() {
	expectedUser := model.User{
		ID:              "admin-user-id",
		Username:        "adminuser",
		Email:           "admin@example.com",
		IsAdmin:         true,
		IntersectionIDs: []string{},
	}

	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

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

func (suite *TestSuite) TestGetProfile_MissingUserID() {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "GetProfile")
}

func (suite *TestSuite) TestGetProfile_ServiceNotFound() {
	emptyResponse := model.User{}
	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewNotFoundError("user not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_ServiceForbiddenError() {
	emptyResponse := model.User{}
	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewForbiddenError("access denied", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "access denied")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_ServiceInternalError() {
	emptyResponse := model.User{}
	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_ServiceValidationError() {
	emptyResponse := model.User{}
	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewValidationError("invalid user ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid user ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_ServiceUnauthorizedError() {
	emptyResponse := model.User{}
	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_UserWithManyIntersections() {
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

	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

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

func (suite *TestSuite) TestGetProfile_ResponseHeadersSet() {
	expectedUser := model.User{
		ID:       "test-id",
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetProfile_UserWithNoIntersections() {
	expectedUser := model.User{
		ID:              "new-user-id",
		Username:        "newuser",
		Email:           "new@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{},
	}

	suite.service.On("GetProfile", mock.Anything, "test-user-id").
		Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("new-user-id", response.ID)
	suite.Equal("newuser", response.Username)
	suite.Empty(response.IntersectionIDs)

	suite.service.AssertExpectations(suite.T())
}
