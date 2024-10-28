package raw_entity

type MarkerInfo struct {
	Data map[string]MarkerData `json:"data,omitempty"`
}

type MarkerData struct {
	Code string `json:"_cd"`
}
