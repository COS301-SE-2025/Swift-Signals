package admin

import (
	"bytes"
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

func (suite *TestSuite) TestGetAllUsers_Success() {
	expectedRequest := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
	}
	expectedUsers := []model.User{
		{
			ID:              "user1",
			Username:        "user1",
			Email:           "user1@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{"int1", "int2"},
		},
		{
			ID:              "user2",
			Username:        "user2",
			Email:           "user2@example.com",
			IsAdmin:         true,
			IntersectionIDs: []string{},
		},
	}

	suite.service.On("GetAllUsers", mock.Anything, expectedRequest.Page, expectedRequest.PageSize).
		Return(expectedUsers, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response, 2)
	suite.Equal(expectedUsers[0].ID, response[0].ID)
	suite.Equal(expectedUsers[0].Username, response[0].Username)
	suite.Equal(expectedUsers[1].ID, response[1].ID)
	suite.Equal(expectedUsers[1].Username, response[1].Username)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidJSON() {
	req := httptest.NewRequest(
		http.MethodGet,
		"/admin/users",
		bytes.NewBufferString(`{"invalid": json}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_MissingPage() {
	requestBody := model.GetAllUsersRequest{
		Page:     0, // Invalid - required and min=1
		PageSize: 10,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request parameters")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_MissingPageSize() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 0, // Invalid - required and min=1
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request parameters")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_PageSizeTooLarge() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 101, // Invalid - max=100
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request parameters")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ServiceForbiddenError() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return([]model.User{}, errs.NewForbiddenError("only admins can access this endpoint", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "only admins can access this endpoint")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ServiceInternalError() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return([]model.User{}, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyUserList() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return([]model.User{}, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Empty(response)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyBody() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_NilBody() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_MinimalValidRequest() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 1,
	}

	expectedUsers := []model.User{
		{
			ID:              "single-user",
			Username:        "singleuser",
			Email:           "single@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{},
		},
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return(expectedUsers, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response, 1)
	suite.Equal("single-user", response[0].ID)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_MaxPageSize() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 100, // Maximum allowed
	}

	expectedUsers := make([]model.User, 100)
	for i := 0; i < 100; i++ {
		expectedUsers[i] = model.User{
			ID:       fmt.Sprintf("user%d", i+1),
			Username: fmt.Sprintf("user%d", i+1),
			Email:    fmt.Sprintf("user%d@example.com", i+1),
			IsAdmin:  false,
		}
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return(expectedUsers, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response, 100)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ResponseHeadersSet() {
	requestBody := model.GetAllUsersRequest{
		Page:     1,
		PageSize: 10,
	}

	suite.service.On("GetAllUsers", mock.Anything, requestBody.Page, requestBody.PageSize).
		Return([]model.User{}, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodGet, "/admin/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerGetAllUsers(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
