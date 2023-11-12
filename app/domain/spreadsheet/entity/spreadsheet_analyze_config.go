package entity

type SpreadSheetAnalyzeConfig struct {
	Id         string      `json:"spreadsheet_id"`
	SheetNames []SheetName `json:"spreadsheet_sheet_names"`
}

type SheetName struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}
