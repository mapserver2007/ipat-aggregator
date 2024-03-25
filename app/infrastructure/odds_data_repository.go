package infrastructure

import (
	"context"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
)

type oddsDataRepository struct {
	client *colly.Collector
}

func NewOddsDataRepository() repository.OddsDataRepository {
	return &oddsDataRepository{
		client: colly.NewCollector(),
	}
}

func (o *oddsDataRepository) Fetch(ctx context.Context, url string) ([]*netkeiba_entity.Odds, error) {

	err := o.client.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to visit url: %s, %v", url, err)
	}

	return nil, nil
}
