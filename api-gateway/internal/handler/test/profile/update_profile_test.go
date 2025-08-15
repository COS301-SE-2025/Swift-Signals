package profile

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestUpdateProfile_Success() {
	expectedRequest := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "updateduser",
		Email:           "updated@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"int1", "int2"},
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", expectedRequest).
		Return(expectedUser, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

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

func (suite *TestSuite) TestUpdateProfile_MissingUserID() {
	requestBody := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_InvalidJSON() {
	req := httptest.NewRequest(
		http.MethodPatch,
		"/me",
		bytes.NewBufferString(`{"invalid": json}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_ValidationError_UsernameTooShort() {
	requestBody := model.UpdateUserRequest{
		Username: "ab", // Invalid - min=3
		Email:    "valid@example.com",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid input data")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_ValidationError_UsernameTooLong() {
	requestBody := model.UpdateUserRequest{
		Username: "thisusernameiswaytoolongandexceedsthemaximumlengthallowed", // Invalid - max=32
		Email:    "valid@example.com",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid input data")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_ValidationError_InvalidEmail() {
	requestBody := model.UpdateUserRequest{
		Username: "validuser",
		Email:    "invalid-email-format", // Invalid email format
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid input data")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_PartialUpdate_UsernameOnly() {
	requestBody := model.UpdateUserRequest{
		Username: "newusername",
		Email:    "", // Empty email should be valid for partial update with omitempty
	}
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "newusername",
		Email:           "original@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{"int1"},
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(expectedUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("newusername", response.Username)
	suite.Equal("original@example.com", response.Email)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_PartialUpdate_EmailOnly() {
	requestBody := model.UpdateUserRequest{
		Username: "", // Empty username should be valid for partial update with omitempty
		Email:    "newemail@example.com",
	}
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "originaluser",
		Email:           "newemail@example.com",
		IsAdmin:         true,
		IntersectionIDs: []string{},
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(expectedUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("originaluser", response.Username)
	suite.Equal("newemail@example.com", response.Email)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceNotFound() {
	requestBody := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewNotFoundError("user not found", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceForbiddenError() {
	requestBody := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewForbiddenError("access denied", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "access denied")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceInternalError() {
	requestBody := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceValidationError() {
	requestBody := model.UpdateUserRequest{
		Username: "duplicateuser",
		Email:    "duplicate@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewValidationError("username already exists", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "username already exists")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceAlreadyExistsError() {
	requestBody := model.UpdateUserRequest{
		Username: "existinguser",
		Email:    "existing@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewAlreadyExistsError("email already in use", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Contains(w.Body.String(), "email already in use")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ServiceUnauthorizedError() {
	requestBody := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	emptyResponse := model.User{}
	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewUnauthorizedError("token expired", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_EmptyBody() {
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_NilBody() {
	req := httptest.NewRequest(http.MethodPatch, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateProfile")
}

func (suite *TestSuite) TestUpdateProfile_MinimalValidUsername() {
	requestBody := model.UpdateUserRequest{
		Username: "abc", // Minimum valid length
		Email:    "valid@example.com",
	}
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "abc",
		Email:           "valid@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{},
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(expectedUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("abc", response.Username)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_MaximalValidUsername() {
	requestBody := model.UpdateUserRequest{
		Username: "12345678901234567890123456789012", // Maximum valid length (32 chars)
		Email:    "valid@example.com",
	}
	expectedUser := model.User{
		ID:              "test-user-id",
		Username:        "12345678901234567890123456789012",
		Email:           "valid@example.com",
		IsAdmin:         false,
		IntersectionIDs: []string{},
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(expectedUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("12345678901234567890123456789012", response.Username)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateProfile_ResponseHeadersSet() {
	requestBody := model.UpdateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
	}
	expectedUser := model.User{
		ID:       "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	suite.service.On("UpdateProfile", mock.Anything, "test-user-id", requestBody).
		Return(expectedUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateProfile(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}
