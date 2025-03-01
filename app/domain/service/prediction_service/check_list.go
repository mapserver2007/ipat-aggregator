package prediction_service

import "context"

type CheckList interface {
	GetPositivePoint(ctx context.Context) error
	GetNegativePoint(ctx context.Context) error
}

type checkList struct {
}
