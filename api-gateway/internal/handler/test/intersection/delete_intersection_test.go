package intersection

import (
	"net/http"
	"net/http/httptest"
	"testing"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func (suite *TestSuite) TestDeleteIntersection_Success() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_MissingUserID() {
	req := httptest.NewRequest(http.MethodDelete, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	// Don't set context with user ID
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertNotCalled(suite.T(), "DeleteIntersectionByID")
}

func (suite *TestSuite) TestDeleteIntersection_NotFound() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "nonexistent-id").
		Return(errs.NewNotFoundError("intersection not found", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/nonexistent-id", nil)
	req.SetPathValue("id", "nonexistent-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "intersection not found")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_Forbidden() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "forbidden-id").
		Return(errs.NewForbiddenError("intersection not in user's intersection list", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/forbidden-id", nil)
	req.SetPathValue("id", "forbidden-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
	suite.Contains(w.Body.String(), "intersection not in user's intersection list")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_InternalError() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(errs.NewInternalError("database connection failed", nil, map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "something went wrong")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_ValidationError() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "invalid-id").
		Return(errs.NewValidationError("invalid intersection ID format", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/invalid-id", nil)
	req.SetPathValue("id", "invalid-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid intersection ID format")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_UnauthorizedError() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "test-intersection-id").
		Return(errs.NewUnauthorizedError("token expired", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/test-intersection-id", nil)
	req.SetPathValue("id", "test-intersection-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "token expired")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_EmptyPathValue() {
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "").
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/intersections/", nil)
	req.SetPathValue("id", "")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_AlreadyExistsError() {
	// This is an edge case - intersection exists but cannot be deleted due to constraints
	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", "constrained-id").
		Return(errs.NewAlreadyExistsError("intersection has dependent resources and cannot be deleted", map[string]any{}))

	req := httptest.NewRequest(http.MethodDelete, "/intersections/constrained-id", nil)
	req.SetPathValue("id", "constrained-id")
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Contains(w.Body.String(), "intersection has dependent resources and cannot be deleted")

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_MultipleSuccessfulDeletes() {
	// Test that multiple delete operations work independently
	intersectionIds := []string{"id1", "id2", "id3"}

	for _, id := range intersectionIds {
		suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", id).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/intersections/"+id, nil)
		req.SetPathValue("id", id)
		req = req.WithContext(suite.ctx)
		w := httptest.NewRecorder()

		suite.handler.DeleteIntersection(w, req)

		suite.Equal(http.StatusNoContent, w.Code)
		suite.Empty(w.Body.String())
	}

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_LongIntersectionID() {
	longId := "very-long-intersection-id-that-might-cause-issues-with-some-systems-but-should-still-be-handled-properly-by-the-delete-endpoint-implementation"

	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", longId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/intersections/"+longId, nil)
	req.SetPathValue("id", longId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func (suite *TestSuite) TestDeleteIntersection_SpecialCharactersInID() {
	specialId := "test-id-with-special-chars-123_456"

	suite.service.On("DeleteIntersectionByID", mock.Anything, "test-user-id", specialId).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/intersections/"+specialId, nil)
	req.SetPathValue("id", specialId)
	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.handler.DeleteIntersection(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
	suite.Empty(w.Body.String())

	suite.service.AssertExpectations(suite.T())
}

func TestHandlerDeleteIntersection(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
