package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestLogin_Success() {
	expectedRequest := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "8characters",
	}
	expectedResponse := model.LoginResponse{
		Message: "Login successful",
		Token:   "very.real.token",
	}

	suite.service.On("LoginUser", mock.Anything, expectedRequest).Return(expectedResponse, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedResponse, response)
	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogin_InvalidJSON() {
	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBufferString(`{"invalid": json}`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func (suite *TestSuite) TestLogin_MissingEmail() {
	requestBody := model.LoginRequest{
		Email:    "",
		Password: "8characters",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Email and password are required")
}

func (suite *TestSuite) TestLogin_MissingPassword() {
	requestBody := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Email and password are required")
}

func (suite *TestSuite) TestLogin_AllFieldsMissing() {
	requestBody := model.LoginRequest{
		Email:    "",
		Password: "",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Email and password are required")
}

func (suite *TestSuite) TestLogin_ServiceError() {
	expectedRequest := model.LoginRequest{
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	emptyResponse := model.LoginResponse{}
	suite.service.On("LoginUser", mock.Anything, expectedRequest).
		Return(emptyResponse, errs.NewUnauthorizedError("invalid credentials", map[string]any{}))

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "invalid credentials")
	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestLogin_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func (suite *TestSuite) TestLogin_NilBody() {
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Login(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func TestHandlerLogin(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
