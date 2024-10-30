package prediction_service

const (
	raceCardUrl            = "https://race.netkeiba.com/race/shutuba.html?race_id=%s"
	raceListUrlForJRA      = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	oddsUrl                = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=1&action=update"
	raceResultUrl          = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=1&race_date=%s"
	raceMarkerUrl          = "https://race.netkeiba.com/api/api_post_social_cart.html?race_id=%s"
	horseUrl               = "https://db.netkeiba.com/horse/%s"
	trainerUrl             = "https://db.netkeiba.com/trainer/%s"
	raceForecastUrl        = "https://tospo-keiba.jp/race/detail/%s/forecast"
	raceTrainingCommentUrl = "https://tospo-keiba.jp/race/detail/%s/comment"
	jockeyFileName         = "jockey.json"
)
