package handler

import (
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/internal/service"
	intersectionpb "github.com/COS301-SE-2025/Swift-Signals/protos/gen/intersection"
	// "google.golang.org/protobuf/types/known/emptypb"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	intersectionpb.UnimplementedIntersectionServiceServer
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}
