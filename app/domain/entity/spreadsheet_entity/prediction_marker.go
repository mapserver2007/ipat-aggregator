package spreadsheet_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type PredictionMarker struct {
	raceId                   string
	favoriteHorseNumber      int
	rivalHorseNumber         int
	brackTriangleHorseNumber int
	whiteTriangleHorseNumber int
	starHorseNumber          int
	checkHorseNumber         int
}

func NewPredictionMarker(
	raceId types.RaceId,
	favoriteHorseNumber types.HorseNumber,
	rivalHorseNumber types.HorseNumber,
	brackTriangleHorseNumber types.HorseNumber,
	whiteTriangleHorseNumber types.HorseNumber,
	starHorseNumber types.HorseNumber,
	checkHorseNumber types.HorseNumber,
) *PredictionMarker {
	return &PredictionMarker{
		raceId:                   raceId.String(),
		favoriteHorseNumber:      favoriteHorseNumber.Value(),
		rivalHorseNumber:         rivalHorseNumber.Value(),
		brackTriangleHorseNumber: brackTriangleHorseNumber.Value(),
		whiteTriangleHorseNumber: whiteTriangleHorseNumber.Value(),
		starHorseNumber:          starHorseNumber.Value(),
		checkHorseNumber:         checkHorseNumber.Value(),
	}
}

func (m *PredictionMarker) RaceId() string {
	return m.raceId
}

func (m *PredictionMarker) FavoriteHorseNumber() int {
	return m.favoriteHorseNumber
}

func (m *PredictionMarker) RivalHorseNumber() int {
	return m.rivalHorseNumber
}

func (m *PredictionMarker) BrackTriangleHorseNumber() int {
	return m.brackTriangleHorseNumber
}

func (m *PredictionMarker) WhiteTriangleHorseNumber() int {
	return m.whiteTriangleHorseNumber
}

func (m *PredictionMarker) StarHorseNumber() int {
	return m.starHorseNumber
}

func (m *PredictionMarker) CheckHorseNumber() int {
	return m.checkHorseNumber
}
