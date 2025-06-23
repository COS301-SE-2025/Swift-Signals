package intersection

import (
	"github.com/COS301-SE-2025/Swift-Signals/intersection-service/db"
)

type Service struct {
	repo db.IntersectionRepository
}

func NewService(r db.IntersectionRepository) *Service {
	return &Service{repo: r}
}
