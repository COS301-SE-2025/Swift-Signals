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

func (suite *TestSuite) TestGetIntersection_Success() {
	expectedIntersection := model.Intersection{
		ID:             "test-intersection-id",
		Name:           "Test Intersection",
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
		DefaultParameters: model.OptimisationParameters{
			OptimisationType: "grid_search",
			SimulationParameters: model.SimulationParameters{
				IntersectionType: "t-junction",
				Green:            10,
				Yellow:           2,
				Red:              6,
				Speed:            60,
				Seed:             123456,
			},
		},
	}

	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(expectedIntersection, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersection
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal(expectedIntersection.ID, response.ID)
	suite.Equal(expectedIntersection.Name, response.Name)
	suite.Equal(expectedIntersection.Status, response.Status)
	suite.Equal(expectedIntersection.TrafficDensity, response.TrafficDensity)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_MissingUserID() {
	req := httptest.NewRequest(http.MethodGet, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	// Service should not be called when user ID is missing
	suite.service.AssertNotCalled(suite.T(), "GetIntersectionByID")
}

func (suite *TestSuite) TestGetIntersection_NotFound() {
	emptyResponse := model.Intersection{}
	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "nonexistent-id").
		Return(emptyResponse, errs.NewNotFoundError("intersection not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections/nonexistent-id", nil)
	req.SetPathValue("id", "nonexistent-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "intersection not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_Forbidden() {
	emptyResponse := model.Intersection{}
	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "forbidden-id").
		Return(emptyResponse, errs.NewForbiddenError("intersection not in user's intersection list", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections/forbidden-id", nil)
	req.SetPathValue("id", "forbidden-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "intersection not in user's intersection list")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_InternalError() {
	emptyResponse := model.Intersection{}
	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(emptyResponse, errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_ValidationError() {
	emptyResponse := model.Intersection{}
	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "invalid-id").
		Return(emptyResponse, errs.NewValidationError("invalid intersection ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid intersection ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_EmptyPathValue() {
	expectedIntersection := model.Intersection{
		ID:             "",
		Name:           "Empty ID Test",
		Status:         "unoptimised",
		TrafficDensity: "medium",
	}

	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "").
		Return(expectedIntersection, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections/", nil)
	req.SetPathValue("id", "")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var response model.Intersection
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	suite.Equal("", response.ID)
	suite.Equal("Empty ID Test", response.Name)

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_UnauthorizedError() {
	emptyResponse := model.Intersection{}
	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(emptyResponse, errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodGet, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestGetIntersection_ResponseHeadersSet() {
	expectedIntersection := model.Intersection{
		ID:             "test-id",
		Name:           "Test Intersection",
		Status:         "unoptimised",
		TrafficDensity: "low",
	}

	suite.service.On("GetIntersectionByID", mock.Anything, "test-user-id", "test-id").
		Return(expectedIntersection, nil)

	req := httptest.NewRequest(http.MethodGet, "/intersections/test-id", nil)
	req.SetPathValue("id", "test-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.GetIntersection(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerGetIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
