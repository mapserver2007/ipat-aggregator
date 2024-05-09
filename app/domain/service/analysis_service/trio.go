package analysis_service

import "context"

type Trio interface {
	Create(ctx context.Context) error
}
