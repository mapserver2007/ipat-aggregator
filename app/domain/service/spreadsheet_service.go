package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"google.golang.org/api/sheets/v4"
)

type SpreadSheetService interface {
	GetCellColor(ctx context.Context, colorType types.CellColorType) *sheets.Color
}

type spreadSheetService struct{}

func NewSpreadSheetService() SpreadSheetService {
	return &spreadSheetService{}
}

func (s *spreadSheetService) GetCellColor(
	ctx context.Context,
	colorType types.CellColorType,
) *sheets.Color {
	switch colorType {
	case types.FirstColor:
		return &sheets.Color{
			Red:   1.0,
			Green: 0.937,
			Blue:  0.498,
		}
	case types.SecondColor:
		return &sheets.Color{
			Red:   0.796,
			Green: 0.871,
			Blue:  1.0,
		}
	case types.ThirdColor:
		return &sheets.Color{
			Red:   0.937,
			Green: 0.78,
			Blue:  0.624,
		}
	}
	return &sheets.Color{
		Red:   1.0,
		Blue:  1.0,
		Green: 1.0,
	}
}
