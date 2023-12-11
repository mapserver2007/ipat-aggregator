package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	jockey_vo "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type NetKeibaService interface {
	CreateRacingNumberUrls(ctx context.Context, tickets []*ticket_csv_entity.Ticket, rawRacingNumbers []*raw_entity.RacingNumber) ([]string, error)
	CreateRaceUrls(ctx context.Context, tickets []*ticket_csv_entity.Ticket, rawRaces []*raw_entity.Race, rawRacingNumbers []*raw_entity.RacingNumber) ([]string, error)
	CreateJockeyUrls(ctx context.Context, rawJockeys []*raw_entity.Jockey, excludeJockeyIds []int) ([]string, error)
}

const (
	raceListUrlForJRA       = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	jockeyUrl               = "https://db.netkeiba.com/jockey/%s/"
)

type netKeibaService struct {
	raceConverter RaceConverter
}

func NewNetKeibaService(
	raceConverter RaceConverter,
) NetKeibaService {
	return &netKeibaService{
		raceConverter: raceConverter,
	}
}

func (n *netKeibaService) CreateRacingNumberUrls(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	rawRacingNumbers []*raw_entity.RacingNumber,
) ([]string, error) {
	racingNumberMap := n.raceConverter.ConvertToRawRacingNumberMap(ctx, rawRacingNumbers)
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
	rawRaces []*raw_entity.Race,
	rawRacingNumbers []*raw_entity.RacingNumber,
) ([]string, error) {
	raceMap := n.raceConverter.ConvertToRawRaceMap(ctx, rawRaces)
	ticketMap := n.raceConverter.ConvertToTicketMap(ctx, tickets, rawRacingNumbers)
	raceUrlCache := map[types.RaceId]string{}

	for raceId, ticket := range ticketMap {
		var url string
		if _, ok := raceMap[raceId]; ok {
			continue
		}
		if _, ok := raceUrlCache[raceId]; ok {
			continue
		}
		if ticket.RaceCourse().JRA() {
			url = fmt.Sprintf(raceResultUrlForJRA, raceId, types.JRA)
		} else if ticket.RaceCourse().NAR() {
			url = fmt.Sprintf(raceResultUrlForNAR, raceId, types.NAR)
		} else if ticket.RaceCourse().Oversea() {
			url = fmt.Sprintf(raceResultUrlForOversea, raceId, types.OverseaOrganizer)
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
	rawJockeys []*raw_entity.Jockey,
	excludeJockeyIds []int,
) ([]string, error) {
	beginIdForJRA := 422
	endIdForJRA := 2000
	beginIdForNARandOversea := 5000
	endIdForNARandOversea := 5999

	jockeysMap := map[int]bool{}
	for _, jockey := range rawJockeys {
		jockeysMap[jockey.JockeyId] = true
	}

	excludeJockeyIdsMap := map[int]jockey_vo.JockeyId{}
	for _, rawJockeyId := range excludeJockeyIds {
		excludeJockeyIdsMap[rawJockeyId] = jockey_vo.JockeyId(rawJockeyId)
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
		jockeyId := jockey_vo.JockeyId(i)
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
		jockeyId := jockey_vo.JockeyId(i)
		urls = append(urls, fmt.Sprintf(jockeyUrl, jockeyId.Format()))
	}

	return urls, nil
}
