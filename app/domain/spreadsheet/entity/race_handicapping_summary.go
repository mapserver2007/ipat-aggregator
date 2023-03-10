package entity

func NewRaceHandicappingSummary(
	favorite, contender HorseInfo,
) RaceHandicappingSummary {
	return RaceHandicappingSummary{
		Favorite:  favorite,
		Contender: contender,
	}
}
