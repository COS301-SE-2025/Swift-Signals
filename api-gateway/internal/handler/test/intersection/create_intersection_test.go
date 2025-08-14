package intersection

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
)

func (suite *TestSuite) TestCreateIntersection_Success() {
	expectedRequest := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	expectedRequest.Details.Address = "123 Test St"
	expectedRequest.Details.City = "Pretoria"
	expectedRequest.Details.Province = "Gauteng"

	expectedResponse := model.CreateIntersectionResponse{
		Id: "created-intersection-id",
	}

	suite.service.On("CreateIntersection", mock.Anything, "test-user-id", expectedRequest).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.CreateIntersectionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedResponse.Id, response.Id)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_InvalidJSON() {
	req := httptest.NewRequest(
		http.MethodPost,
		"/intersections",
		bytes.NewBufferString(`{"invalid": json}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_MissingName() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "", // Missing name
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "name, density, and parameters are required")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_MissingParameters() {
	// Create a request with missing default_parameters field
	requestBody := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		// DefaultParameters is missing - this will be zero value which should fail validation
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "name, density, and parameters are required")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_MissingUserID() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_ServiceError() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	emptyResponse := model.CreateIntersectionResponse{}
	suite.service.On("CreateIntersection", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_ValidationError() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	emptyResponse := model.CreateIntersectionResponse{}
	suite.service.On("CreateIntersection", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewValidationError("invalid traffic density", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid traffic density")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_AlreadyExistsError() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "Test Intersection",
		TrafficDensity: "high",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "t-junction",
			Green:            10,
			Yellow:           2,
			Red:              6,
			Speed:            60,
			Seed:             123456,
		},
	}
	requestBody.Details.Address = "123 Test St"
	requestBody.Details.City = "Pretoria"
	requestBody.Details.Province = "Gauteng"

	emptyResponse := model.CreateIntersectionResponse{}
	suite.service.On("CreateIntersection", mock.Anything, "test-user-id", requestBody).
		Return(emptyResponse, errs.NewAlreadyExistsError("intersection already exists at this location", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Contains(w.Body.String(), "intersection already exists at this location")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestCreateIntersection_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_NilBody() {
	req := httptest.NewRequest(http.MethodPost, "/intersections", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "CreateIntersection")
}

func (suite *TestSuite) TestCreateIntersection_MinimalValidRequest() {
	requestBody := model.CreateIntersectionRequest{
		Name:           "Minimal Intersection",
		TrafficDensity: "low",
		DefaultParameters: model.SimulationParameters{
			IntersectionType: "traffic-light",
			Green:            5,
			Yellow:           1,
			Red:              3,
			Speed:            30,
			Seed:             999,
		},
	}

	expectedResponse := model.CreateIntersectionResponse{
		Id: "minimal-intersection-id",
	}

	suite.service.On("CreateIntersection", mock.Anything, "test-user-id", requestBody).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/intersections", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.CreateIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.CreateIntersectionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedResponse.Id, response.Id)

	suite.service.AssertExpectations(suite.T())
}
