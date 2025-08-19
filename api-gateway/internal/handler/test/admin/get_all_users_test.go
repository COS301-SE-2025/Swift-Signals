package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestGetAllUsers_Success() {
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

	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10", nil)
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

func (suite *TestSuite) TestGetAllUsers_DefaultValues() {
	expectedUsers := []model.User{
		{
			ID:              "user1",
			Username:        "user1",
			Email:           "user1@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{"int1", "int2"},
		},
	}

	// When no query params are provided, defaults should be page=1, page_size=100
	suite.service.On("GetAllUsers", mock.Anything, 1, 100).
		Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response []model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response, 1)
	suite.Equal(expectedUsers[0].ID, response[0].ID)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_InvalidQueryParam_NonNumericPage() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=invalid&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page number")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_InvalidQueryParam_NonNumericPageSize() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=invalid", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page size")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_PageZero() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=0&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page number")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_PageSizeZero() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=0", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page size")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_PageSizeTooLarge() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=101", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page size")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_NegativePage() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=-1&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page number")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ValidationError_NegativePageSize() {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=-5", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page size")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_ServiceForbiddenError() {
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, errs.NewForbiddenError("only admins can access this endpoint", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "only admins can access this endpoint")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_ServiceInternalError() {
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyUserList() {
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10", nil)
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

func (suite *TestSuite) TestGetAllUsers_MinimalValidRequest() {
	expectedUsers := []model.User{
		{
			ID:              "single-user",
			Username:        "singleuser",
			Email:           "single@example.com",
			IsAdmin:         false,
			IntersectionIDs: []string{},
		},
	}

	suite.service.On("GetAllUsers", mock.Anything, 1, 1).
		Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=1", nil)
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
	expectedUsers := make([]model.User, 100)
	for i := 0; i < 100; i++ {
		expectedUsers[i] = model.User{
			ID:       fmt.Sprintf("user%d", i+1),
			Username: fmt.Sprintf("user%d", i+1),
			Email:    fmt.Sprintf("user%d@example.com", i+1),
			IsAdmin:  false,
		}
	}

	suite.service.On("GetAllUsers", mock.Anything, 1, 100).
		Return(expectedUsers, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=100", nil)
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
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&page_size=10", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_LargePageNumber() {
	suite.service.On("GetAllUsers", mock.Anything, 999999, 10).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=999999&page_size=10", nil)
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

func (suite *TestSuite) TestGetAllUsers_ExtraQueryParameters() {
	// Test that extra query parameters are ignored
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin/users?page=1&page_size=10&extra=ignored&another=param",
		nil,
	)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_DuplicateQueryParameters() {
	// Test behavior when duplicate query parameters are provided (should use first value)
	suite.service.On("GetAllUsers", mock.Anything, 1, 10).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/admin/users?page=1&page=2&page_size=10&page_size=20",
		nil,
	)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_EmptyQueryParameterValues() {
	// Empty page parameter should use default (page=1)
	suite.service.On("GetAllUsers", mock.Anything, 1, 100).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=&page_size=", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_WhitespaceQueryParameterValues() {
	// Whitespace values should be treated as invalid and return error
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=%20&page_size=%20", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid page number")

	suite.service.AssertNotCalled(suite.T(), "GetAllUsers")
}

func (suite *TestSuite) TestGetAllUsers_OnlyPageProvided() {
	// Only page provided, page_size should use default (100)
	suite.service.On("GetAllUsers", mock.Anything, 5, 100).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=5", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllUsers_OnlyPageSizeProvided() {
	// Only page_size provided, page should use default (1)
	suite.service.On("GetAllUsers", mock.Anything, 1, 25).
		Return([]model.User{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users?page_size=25", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllUsers(w, req)

	suite.Equal(http.StatusOK, w.Code)

	suite.service.AssertExpectations(suite.T())
}
