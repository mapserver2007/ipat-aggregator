package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/shopspring/decimal"
)

type PlaceUnHit interface {
	Convert(ctx context.Context, calculables []*analysis_entity.PlaceAllInCalculable, upperOdds float64, lowerOdds float64) error
}

type placeUnHitService struct {
	placeAllInService PlaceAllIn
}

func NewPlaceUnHit(
	placeAllInService PlaceAllIn,
) PlaceUnHit {
	return &placeUnHitService{
		placeAllInService: placeAllInService,
	}
}

func (p *placeUnHitService) Convert(
	ctx context.Context,
	calculables []*analysis_entity.PlaceAllInCalculable,
	upperOdds float64,
	lowerOdds float64,
) error {
	isHit := func(calculable *analysis_entity.PlaceAllInCalculable) bool {
		return (calculable.Entries() <= 7 && calculable.OrderNo() <= 2) || (calculable.Entries() >= 8 && calculable.OrderNo() <= 3)
	}
	decimalUpperOdds := decimal.NewFromFloat(upperOdds)
	decimalLowerOdds := decimal.NewFromFloat(lowerOdds)

	//var placeUnHitCalculables []*analysis_entity.PlaceAllInCalculable
	for _, calculable := range calculables {
		if !isHit(calculable) && calculable.WinOdds().GreaterThanOrEqual(decimalUpperOdds) && calculable.WinOdds().LessThanOrEqual(decimalLowerOdds) {
			fmt.Println(calculable)
		}
	}

	return nil
}
