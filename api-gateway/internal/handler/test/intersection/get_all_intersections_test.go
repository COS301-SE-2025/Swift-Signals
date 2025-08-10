package intersection

import (
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

func (suite *TestSuite) TestGetAllIntersections_Success() {
	expectedIntersections := model.Intersections{
		Intersections: []model.Intersection{
			{
				ID:             "1",
				Name:           "Test Intersection 1",
				CreatedAt:      time.Now(),
				LastRunAt:      time.Now(),
				Status:         "unoptimised",
				RunCount:       5,
				TrafficDensity: "high",
				Details: model.Details{
					Address:  "123 Test St",
					City:     "Pretoria",
					Province: "Gauteng",
				},
			},
			{
				ID:             "2",
				Name:           "Test Intersection 2",
				CreatedAt:      time.Now(),
				LastRunAt:      time.Now(),
				Status:         "optimised",
				RunCount:       3,
				TrafficDensity: "medium",
				Details: model.Details{
					Address:  "456 Main Ave",
					City:     "Cape Town",
					Province: "Western Cape",
				},
			},
		},
	}

	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(expectedIntersections, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersections
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedIntersections.Intersections[0].ID, response.Intersections[0].ID)
	suite.Equal(expectedIntersections.Intersections[0].Name, response.Intersections[0].Name)
	suite.Equal(expectedIntersections.Intersections[1].ID, response.Intersections[1].ID)
	suite.Equal(expectedIntersections.Intersections[1].Name, response.Intersections[1].Name)
	suite.Len(response.Intersections, 2)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_EmptyList() {
	expectedIntersections := model.Intersections{
		Intersections: []model.Intersection{},
	}

	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(expectedIntersections, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersections
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Empty(response.Intersections)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_MissingUserID() {
	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	// Service should not be called when user ID is missing
	suite.service.AssertNotCalled(suite.T(), "GetAllIntersections")
}

func (suite *TestSuite) TestGetAllIntersections_ServiceInternalError() {
	emptyResponse := model.Intersections{}
	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ServiceNotFoundError() {
	emptyResponse := model.Intersections{}
	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewNotFoundError("user not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ServiceUnauthorizedError() {
	emptyResponse := model.Intersections{}
	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewUnauthorizedError("invalid token", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "invalid token")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ServiceForbiddenError() {
	emptyResponse := model.Intersections{}
	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewForbiddenError("access denied", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "access denied")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ServiceValidationError() {
	emptyResponse := model.Intersections{}
	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(emptyResponse, errs.NewValidationError("invalid user ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid user ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_SingleIntersection() {
	expectedIntersections := model.Intersections{
		Intersections: []model.Intersection{
			{
				ID:             "single-id",
				Name:           "Single Test Intersection",
				CreatedAt:      time.Now(),
				LastRunAt:      time.Now(),
				Status:         "optimising",
				RunCount:       1,
				TrafficDensity: "low",
				Details: model.Details{
					Address:  "789 Single Rd",
					City:     "Durban",
					Province: "KwaZulu-Natal",
				},
			},
		},
	}

	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(expectedIntersections, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersections
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Len(response.Intersections, 1)
	suite.Equal("single-id", response.Intersections[0].ID)
	suite.Equal("Single Test Intersection", response.Intersections[0].Name)
	suite.Equal("optimising", response.Intersections[0].Status)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetAllIntersections_ResponseHeadersSet() {
	expectedIntersections := model.Intersections{
		Intersections: []model.Intersection{},
	}

	suite.service.On("GetAllIntersections", mock.Anything, "test-user-id").
		Return(expectedIntersections, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections", nil)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetAllIntersections(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerGetAllIntersections(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
