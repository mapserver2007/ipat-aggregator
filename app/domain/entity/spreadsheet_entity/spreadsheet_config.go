package spreadsheet_entity

type SpreadSheetConfig struct {
	spreadSheetId string
	sheetId       int64
	sheetName     string
}

func NewSpreadSheetConfig(spreadSheetId string, sheetId int64, sheetName string) *SpreadSheetConfig {
	return &SpreadSheetConfig{
		spreadSheetId: spreadSheetId,
		sheetId:       sheetId,
		sheetName:     sheetName,
	}
}

func (s *SpreadSheetConfig) SpreadSheetId() string {
	return s.spreadSheetId
}

func (s *SpreadSheetConfig) SheetId() int64 {
	return s.sheetId
}

func (s *SpreadSheetConfig) SheetName() string {
	return s.sheetName
}
