package raw_entity

type SpreadSheetConfig struct {
	Id        string `json:"spreadsheet_id"`
	SheetName string `json:"spreadsheet_sheet_name"`
}

type SpreadSheetConfigs struct {
	Id         string   `json:"spreadsheet_id"`
	SheetNames []string `json:"spreadsheet_sheet_name"`
}
