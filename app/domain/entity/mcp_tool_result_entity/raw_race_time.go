package mcp_tool_result_entity

type RaceTimeResult struct {
	FastestTime RaceTime `json:"fastest_time"`
	SlowestTime RaceTime `json:"slowest_time"`
}

type RaceTime struct {
	RaceTime string `json:"race_time"`
	RaceName string `json:"race_name"`
	RaceDate string `json:"race_date"`
	Distance int    `json:"distance"`
	Class    string `json:"class"`
}
