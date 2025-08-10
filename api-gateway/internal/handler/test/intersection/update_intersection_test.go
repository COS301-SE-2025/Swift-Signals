package intersection

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/model"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestUpdateIntersection_Success() {
	expectedRequest := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}
	expectedRequest.Details.Address = "456 Updated St"
	expectedRequest.Details.City = "Cape Town"
	expectedRequest.Details.Province = "Western Cape"

	expectedResponse := model.Intersection{
		ID:             "test-intersection-id",
		Name:           "Updated Intersection",
		CreatedAt:      time.Now(),
		LastRunAt:      time.Now(),
		Status:         "unoptimised",
		RunCount:       5,
		TrafficDensity: "medium",
		Details: model.Details{
			Address:  "456 Updated St",
			City:     "Cape Town",
			Province: "Western Cape",
		},
	}

	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", expectedRequest).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(expectedRequest)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersection
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedResponse.ID, response.ID)
	suite.Equal(expectedResponse.Name, response.Name)
	suite.Equal(expectedResponse.Details.Address, response.Details.Address)
	suite.Equal(expectedResponse.Details.City, response.Details.City)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_InvalidJSON() {
	req := httptest.NewRequest(
		http.MethodPatch,
		"/intersections/test-intersection-id",
		bytes.NewBufferString(`{"invalid": json}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateIntersectionByID")
}

func (suite *TestSuite) TestUpdateIntersection_MissingUserID() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}
	requestBody.Details.Address = "456 Updated St"
	requestBody.Details.City = "Cape Town"
	requestBody.Details.Province = "Western Cape"

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "UpdateIntersectionByID")
}

func (suite *TestSuite) TestUpdateIntersection_NotFound() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}
	requestBody.Details.Address = "456 Updated St"
	requestBody.Details.City = "Cape Town"
	requestBody.Details.Province = "Western Cape"

	emptyResponse := model.Intersection{}
	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "nonexistent-id", requestBody).
		Return(emptyResponse, errs.NewNotFoundError("intersection not found", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/nonexistent-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "nonexistent-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "intersection not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_Forbidden() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}
	requestBody.Details.Address = "456 Updated St"
	requestBody.Details.City = "Cape Town"
	requestBody.Details.Province = "Western Cape"

	emptyResponse := model.Intersection{}
	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "forbidden-id", requestBody).
		Return(emptyResponse, errs.NewForbiddenError("intersection not in user's intersection list", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/forbidden-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "forbidden-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "intersection not in user's intersection list")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_ValidationError() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Invalid Name That Is Way Too Long For The System To Handle And Should Trigger A Validation Error Because It Exceeds The Maximum Length Allowed For Intersection Names In The Database Schema And Application Logic",
	}
	requestBody.Details.Address = "456 Updated St"
	requestBody.Details.City = "Cape Town"
	requestBody.Details.Province = "Western Cape"

	emptyResponse := model.Intersection{}
	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", requestBody).
		Return(emptyResponse, errs.NewValidationError("intersection name too long", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "intersection name too long")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_InternalError() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}
	requestBody.Details.Address = "456 Updated St"
	requestBody.Details.City = "Cape Town"
	requestBody.Details.Province = "Western Cape"

	emptyResponse := model.Intersection{}
	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", requestBody).
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_EmptyBody() {
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateIntersectionByID")
}

func (suite *TestSuite) TestUpdateIntersection_NilBody() {
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "Invalid request payload")

	suite.service.AssertNotCalled(suite.T(), "UpdateIntersectionByID")
}

func (suite *TestSuite) TestUpdateIntersection_PartialUpdate_NameOnly() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Only Name Updated",
	}

	expectedResponse := model.Intersection{
		ID:             "test-intersection-id",
		Name:           "Only Name Updated",
		CreatedAt:      time.Now(),
		LastRunAt:      time.Now(),
		Status:         "unoptimised",
		RunCount:       3,
		TrafficDensity: "low",
		Details: model.Details{
			Address:  "Original Address",
			City:     "Original City",
			Province: "Original Province",
		},
	}

	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", requestBody).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersection
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("Only Name Updated", response.Name)
	suite.Equal("Original Address", response.Details.Address)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_PartialUpdate_DetailsOnly() {
	requestBody := model.UpdateIntersectionRequest{}
	requestBody.Details.Address = "New Address Only"
	requestBody.Details.City = "New City Only"
	requestBody.Details.Province = "New Province Only"

	expectedResponse := model.Intersection{
		ID:             "test-intersection-id",
		Name:           "Original Name",
		CreatedAt:      time.Now(),
		LastRunAt:      time.Now(),
		Status:         "optimised",
		RunCount:       8,
		TrafficDensity: "high",
		Details: model.Details{
			Address:  "New Address Only",
			City:     "New City Only",
			Province: "New Province Only",
		},
	}

	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", requestBody).
		Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersection
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("Original Name", response.Name)
	suite.Equal("New Address Only", response.Details.Address)
	suite.Equal("New City Only", response.Details.City)
	suite.Equal("New Province Only", response.Details.Province)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestUpdateIntersection_UnauthorizedError() {
	requestBody := model.UpdateIntersectionRequest{
		Name: "Updated Intersection",
	}

	emptyResponse := model.Intersection{}
	suite.service.On("UpdateIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id", requestBody).
		Return(emptyResponse, errs.NewUnauthorizedError("token expired", map[string]any{}))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/intersections/test-intersection-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.UpdateIntersection(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerUpdateIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
