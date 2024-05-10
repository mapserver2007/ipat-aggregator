package aggregation_usecase

import "context"

type List interface {
	Execute(ctx context.Context) error
}

type ListInput struct {
}

type list struct {
}

func NewList() List {
	return &list{}
}

func (l *list) Execute(ctx context.Context) error {
	return nil
}
