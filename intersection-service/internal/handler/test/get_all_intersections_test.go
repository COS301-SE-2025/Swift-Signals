package test

import (
	"context"
	"errors"
	"testing"

	handlerPkg "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/handler"
	serviceMocks "github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/mocks/service"
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/model"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

// fakeStream implements the gRPC stream interface for testing
type fakeStream struct {
	ctx  context.Context
	sent []*intersectionpb.IntersectionResponse
	err  error
}

func (f *fakeStream) Send(resp *intersectionpb.IntersectionResponse) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, resp)
	return nil
}

func (f *fakeStream) Context() context.Context {
	return f.ctx
}

// gRPC header/trailer methods with correct types
func (f *fakeStream) SendHeader(md metadata.MD) error { return nil }
func (f *fakeStream) SetHeader(md metadata.MD) error  { return nil }
func (f *fakeStream) SetTrailer(md metadata.MD)       {}
func (f *fakeStream) SendMsg(m any) error            { return nil }
func (f *fakeStream) RecvMsg(m any) error            { return nil }

func TestGetAllIntersections_Success(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.GetAllIntersectionsRequest{
		Page:     1,
		PageSize: 2,
		Filter:   "Main",
	}

	intersections := []*model.Intersection{
		{ID: "int-001", Name: "Main & First"},
		{ID: "int-002", Name: "Main & Second"},
	}

	mockService.On("GetAllIntersections",
		mock.Anything,
		int(req.GetPage()),
		int(req.GetPageSize()),
		req.GetFilter(),
	).Return(intersections, nil)

	stream := &fakeStream{ctx: ctx}

	err := h.GetAllIntersections(req, stream)
	assert.NoError(t, err)
	assert.Len(t, stream.sent, 2)
	assert.Equal(t, intersections[0].ID, stream.sent[0].GetId())
	assert.Equal(t, intersections[1].ID, stream.sent[1].GetId())

	mockService.AssertExpectations(t)
}

func TestGetAllIntersections_ServiceError(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.GetAllIntersectionsRequest{}

	mockService.On("GetAllIntersections",
		mock.Anything,
		0,
		0,
		"",
	).Return(nil, errors.New("database failure"))

	stream := &fakeStream{ctx: ctx}

	err := h.GetAllIntersections(req, stream)
	assert.Error(t, err)
	assert.Empty(t, stream.sent)

	mockService.AssertExpectations(t)
}

func TestGetAllIntersections_StreamError(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	ctx := context.Background()
	req := &intersectionpb.GetAllIntersectionsRequest{}

	intersections := []*model.Intersection{
		{ID: "int-001", Name: "Main & First"},
	}

	mockService.On("GetAllIntersections",
		mock.Anything,
		0,
		0,
		"",
	).Return(intersections, nil)

	stream := &fakeStream{
		ctx: ctx,
		err: errors.New("stream send failed"),
	}

	err := h.GetAllIntersections(req, stream)
	assert.Error(t, err)
	assert.Empty(t, stream.sent) // nothing successfully sent

	mockService.AssertExpectations(t)
}

func TestGetAllIntersections_ContextCancelled(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req := &intersectionpb.GetAllIntersectionsRequest{}

	intersections := []*model.Intersection{
		{ID: "int-001", Name: "Main & First"},
	}

	mockService.On("GetAllIntersections", mock.Anything, 0, 0, "").Return(intersections, nil)

	stream := &fakeStream{ctx: ctx}

	err := h.GetAllIntersections(req, stream)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	mockService.AssertExpectations(t)
}

func TestGetAllIntersections_ResponseNil(t *testing.T) {
	mockService := new(serviceMocks.MockIntersectionService)
	h := handlerPkg.NewIntersectionHandler(mockService)

	// Return a slice containing a nil intersection to hit the branch
	mockService.On("GetAllIntersections",
		mock.Anything,
		0,
		0,
		"",
	).Return([]*model.Intersection{nil}, nil)

	stream := &fakeStream{ctx: context.Background()}
	req := &intersectionpb.GetAllIntersectionsRequest{}

	err := h.GetAllIntersections(req, stream)
	assert.NoError(t, err)

	// Since the only element is nil, stream.sent should remain empty
	assert.Len(t, stream.sent, 0)

	mockService.AssertExpectations(t)
}
