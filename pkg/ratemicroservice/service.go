package rate

import (
	"context"
)

type Service interface {
	GetRate(ctx context.Context) (float64, error)
	ServiceStatus(ctx context.Context) (int, error)
}
