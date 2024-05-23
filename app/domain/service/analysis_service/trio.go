package analysis_service

import "context"

type Trio interface {
	Execute(ctx context.Context) error
}
