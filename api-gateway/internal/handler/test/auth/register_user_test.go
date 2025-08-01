package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestRegisterUser_Success() {
	expectedRequest := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}
	expectedResponse := model.RegisterResponse{
		UserID: "generated_id",
	}

	suite.service.On("RegisterUser", mock.Anything, expectedRequest).Return(expectedResponse, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	var response model.RegisterResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(expectedResponse, response)
	suite.service.AssertExpectations(suite.T())
}
func (suite *TestSuite) TestRegister_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"invalid": json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func (suite *TestSuite) TestRegister_MissingUsername() {
	requestBody := model.RegisterRequest{
		Username: "",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Username, email, and password are required")
}

func (suite *TestSuite) TestRegister_MissingEmail() {
	requestBody := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "",
		Password: "8characters",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Username, email, and password are required")
}

func (suite *TestSuite) TestRegister_MissingPassword() {
	requestBody := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Username, email, and password are required")
}

func (suite *TestSuite) TestRegister_AllFieldsMissing() {
	requestBody := model.RegisterRequest{
		Username: "",
		Email:    "",
		Password: "",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Username, email, and password are required")
}

func (suite *TestSuite) TestRegister_ServiceError() {
	expectedRequest := model.RegisterRequest{
		Username: "Valid Name",
		Email:    "valid@gmail.com",
		Password: "8characters",
	}

	emptyResponse := model.RegisterResponse{}
	suite.service.On("RegisterUser", mock.Anything, expectedRequest).Return(emptyResponse, errors.New("user already exists"))

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "user already exists")
	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestRegister_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func (suite *TestSuite) TestRegister_NilBody() {
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.Register(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")
}

func TestHandlerRegisterUser(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
