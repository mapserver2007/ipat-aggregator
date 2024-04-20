package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"time"
)

type NetKeibaService interface {
	CreateRacingNumberUrls(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber) ([]string, error)
	CreateRaceUrls(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber) ([]string, error)
	CreateJockeyUrls(ctx context.Context, jockeys []*data_cache_entity.Jockey, excludeJockeyIds []int) ([]string, error)
	CreateRaceIdUrls(ctx context.Context, dateFrom, dateTo string, raceIdMap map[types.RaceDate][]types.RaceId, excludeDates []types.RaceDate) ([]string, error)
	CreateAnalysisRaceUrls(ctx context.Context, races []*data_cache_entity.Race, raceIdMap map[types.RaceId]types.RaceDate) ([]string, error)
	CreatePredictionRaceUrls(ctx context.Context, raceIds []types.RaceId) ([]string, []string, []string)
	CreateOddsUrls(ctx context.Context, oddsList []*data_cache_entity.Odds, raceIdMap map[types.RaceId]types.RaceDate) ([]string, error)
}

const (
	raceListUrlForJRA       = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d&race_date=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=%d&race_date=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d&race_date=%d"
	jockeyUrl               = "https://db.netkeiba.com/jockey/%s/"
	analysisRaceResultUrl   = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=1&race_date=%d"
	predictionRaceUrl       = "https://race.netkeiba.com/race/shutuba.html?race_id=%s"
	oddsUrl                 = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=%d&action=update"
	predictionRaceResultUrl = "https://race.netkeiba.com/race/result.html?race_id=%s"
)

type netKeibaService struct {
	raceConverter   RaceConverter
	ticketConverter TicketConverter
}

func NewNetKeibaService(
	raceConverter RaceConverter,
	ticketConverter TicketConverter,
) NetKeibaService {
	return &netKeibaService{
		raceConverter:   raceConverter,
		ticketConverter: ticketConverter,
	}
}

func (n *netKeibaService) CreateRacingNumberUrls(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) ([]string, error) {
	racingNumberMap := n.raceConverter.ConvertToRawRacingNumberMap(ctx, racingNumbers)
	racingNumberUrlCache := map[types.RaceDate]string{}
	for _, ticket := range tickets {
		// JRA以外は日付からレース番号の特定が可能のため処理しない
		if !ticket.RaceCourse().JRA() {
			continue
		}
		racingNumberId := types.NewRacingNumberId(
			ticket.RaceDate(),
			ticket.RaceCourse(),
		)
		if _, ok := racingNumberMap[racingNumberId]; ok {
			continue
		}
		if _, ok := racingNumberUrlCache[ticket.RaceDate()]; ok {
			continue
		}
		racingNumberUrlCache[ticket.RaceDate()] = fmt.Sprintf(raceListUrlForJRA, ticket.RaceDate().Value())
	}

	racingNumberUrls := make([]string, 0, len(racingNumberUrlCache))
	for _, racingNumberUrl := range racingNumberUrlCache {
		racingNumberUrls = append(racingNumberUrls, racingNumberUrl)
	}

	return racingNumberUrls, nil
}

func (n *netKeibaService) CreateRaceUrls(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	races []*data_cache_entity.Race,
	racingNumbers []*data_cache_entity.RacingNumber,
) ([]string, error) {
	raceMap := n.raceConverter.ConvertToRaceMap(ctx, races)
	ticketsMap := n.ticketConverter.ConvertToRaceIdMap(ctx, tickets, racingNumbers)
	raceUrlCache := map[types.RaceId]string{}

	for raceId, ticketsByRaceId := range ticketsMap {
		// 馬券からレース情報が抜ければ良いので要素1つだけ抜く
		if len(ticketsByRaceId) == 0 {
			continue
		}
		ticket := ticketsByRaceId[0]

		var url string
		if _, ok := raceMap[raceId]; ok {
			continue
		}
		if _, ok := raceUrlCache[raceId]; ok {
			continue
		}
		if ticket.RaceCourse().JRA() {
			url = fmt.Sprintf(raceResultUrlForJRA, raceId, types.JRA, ticket.RaceDate())
		} else if ticket.RaceCourse().NAR() {
			url = fmt.Sprintf(raceResultUrlForNAR, raceId, types.NAR, ticket.RaceDate())
		} else if ticket.RaceCourse().Oversea() {
			url = fmt.Sprintf(raceResultUrlForOversea, raceId, types.OverseaOrganizer, ticket.RaceDate())
		} else {
			return nil, fmt.Errorf("undefined organizer: race_date %d, race_no %d", ticket.RaceDate(), ticket.RaceNo())
		}

		raceUrlCache[raceId] = url
	}

	raceUrls := make([]string, 0, len(raceUrlCache))
	for _, raceUrl := range raceUrlCache {
		raceUrls = append(raceUrls, raceUrl)
	}

	return raceUrls, nil
}

func (n *netKeibaService) CreateJockeyUrls(
	ctx context.Context,
	jockeys []*data_cache_entity.Jockey,
	excludeJockeyIds []int,
) ([]string, error) {
	beginIdForJRA := 422
	endIdForJRA := 2000
	beginIdForNARandOversea := 5000
	endIdForNARandOversea := 5999

	jockeysMap := map[int]bool{}
	for _, jockey := range jockeys {
		jockeysMap[jockey.JockeyId().Value()] = true
	}

	excludeJockeyIdsMap := map[int]types.JockeyId{}
	for _, rawJockeyId := range excludeJockeyIds {
		excludeJockeyIdsMap[rawJockeyId] = types.JockeyId(rawJockeyId)
	}

	var urls []string
	for i := beginIdForJRA; i <= endIdForJRA; i++ {
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[i]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[i]; ok {
			continue
		}
		jockeyId := types.JockeyId(i)
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId.Format()))
	}
	for i := beginIdForNARandOversea; i <= endIdForNARandOversea; i++ {
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[i]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[i]; ok {
			continue
		}
		jockeyId := types.JockeyId(i)
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId.Format()))
	}

	return urls, nil
}

func (n *netKeibaService) CreateRaceIdUrls(
	ctx context.Context,
	dateFrom, dateTo string,
	raceIdMap map[types.RaceDate][]types.RaceId,
	excludeDates []types.RaceDate,
) ([]string, error) {
	urls := make([]string, 0)
	excludeDateMap := map[types.RaceDate]bool{}
	for _, excludeDate := range excludeDates {
		excludeDateMap[excludeDate] = true
	}

	startTime, _ := time.Parse("20060102", dateFrom)
	endTime, _ := time.Parse("20060102", dateTo)
	for d := startTime; d.Before(endTime) || d.Equal(endTime); d = d.AddDate(0, 0, 1) {
		date, err := types.NewRaceDate(d.Format("20060102"))
		if err != nil {
			return nil, err
		}
		if excludeDateMap != nil {
			if _, ok := excludeDateMap[date]; ok {
				continue
			}
		}
		if raceIdMap == nil {
			urls = append(urls, fmt.Sprintf(raceListUrlForJRA, date))
		} else {
			if _, ok := raceIdMap[date]; !ok {
				urls = append(urls, fmt.Sprintf(raceListUrlForJRA, date))
			}
		}
	}

	return urls, nil
}

func (n *netKeibaService) CreateAnalysisRaceUrls(
	ctx context.Context,
	races []*data_cache_entity.Race,
	raceIdMap map[types.RaceId]types.RaceDate,
) ([]string, error) {
	var raceUrls []string
	raceMap := map[types.RaceId]*data_cache_entity.Race{}

	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	for _, raceId := range SortedRaceIdKeys(raceIdMap) {
		if _, ok := raceMap[raceId]; !ok {
			raceUrls = append(raceUrls, fmt.Sprintf(analysisRaceResultUrl, raceId, raceIdMap[raceId]))
		}
	}

	return raceUrls, nil
}

func (n *netKeibaService) CreatePredictionRaceUrls(ctx context.Context, raceIds []types.RaceId) ([]string, []string, []string) {
	var raceUrls, oddsUrls, raceResultUrls []string
	for _, raceId := range raceIds {
		raceUrls = append(raceUrls, fmt.Sprintf(predictionRaceUrl, raceId))
		oddsUrls = append(oddsUrls, fmt.Sprintf(oddsUrl, raceId, 1))
		raceResultUrls = append(raceResultUrls, fmt.Sprintf(predictionRaceResultUrl, raceId))
	}
	return raceUrls, oddsUrls, raceResultUrls
}

func (n *netKeibaService) CreateOddsUrls(
	ctx context.Context,
	oddsList []*data_cache_entity.Odds,
	raceIdMap map[types.RaceId]types.RaceDate,
) ([]string, error) {
	var trioOddsUrls []string
	oddsMap := map[types.RaceId]bool{}

	for _, odds := range oddsList {
		oddsMap[odds.RaceId()] = true
	}

	for _, raceId := range SortedRaceIdKeys(raceIdMap) {
		if _, ok := oddsMap[raceId]; !ok {
			trioOddsUrls = append(trioOddsUrls, fmt.Sprintf(oddsUrl, raceId, 7))
		}
	}

	return trioOddsUrls, nil
}
