package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	result_summary_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	spreadsheet_vo "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	secretFileName                   = "secret.json"
	spreadSheetSummaryFileName       = "spreadsheet_summary.json"
	spreadSheetTicketSummaryFileName = "spreadsheet_ticket_summary.json"
	spreadSheetListFileName          = "spreadsheet_list.json"
	spreadSheetAnalyzeFileName       = "spreadsheet_analyze.json"
)

type SpreadSheetClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetClient(
	ctx context.Context,
) repository.SpreadSheetClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfigOld(ctx, spreadSheetSummaryFileName)
	return &SpreadSheetClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

type SpreadSheetMonthlyBettingTicketClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetMonthlyBettingTicketClient(
	ctx context.Context,
) repository.SpreadSheetMonthlyBettingTicketClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfigOld(ctx, spreadSheetTicketSummaryFileName)
	return &SpreadSheetMonthlyBettingTicketClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

type SpreadSheetListClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	sheetId           int64
}

func NewSpreadSheetListClient(
	ctx context.Context,
) repository.SpreadSheetListClient {
	service, spreadSheetConfig, sheetId := getSpreadSheetConfigOld(ctx, spreadSheetListFileName)
	return &SpreadSheetListClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetId:           sheetId,
	}
}

type SpreadSheetAnalyzeClient struct {
	client            *sheets.Service
	spreadSheetConfig spreadsheet_entity.SpreadSheetAnalyzeConfig
	sheetMap          map[spreadsheet_vo.AnalyzeType]*sheets.SheetProperties
}

func NewSpreadSheetAnalyzeClient(
	ctx context.Context,
) repository.SpreadSheetAnalyzeClient {
	service, spreadSheetConfig, sheetMap := getSpreadSheetAnalyzeConfig(ctx, spreadSheetAnalyzeFileName)
	return &SpreadSheetAnalyzeClient{
		client:            service,
		spreadSheetConfig: spreadSheetConfig,
		sheetMap:          sheetMap,
	}
}

func getSpreadSheetConfigOld(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, spreadsheet_entity.SpreadSheetConfig, int64) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		panic(err)
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		panic(err)
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		panic(err)
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		panic(err)
	}

	var spreadSheetConfig spreadsheet_entity.SpreadSheetConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &spreadSheetConfig); err != nil {
		panic(err)
	}

	response, err := service.Spreadsheets.Get(spreadSheetConfig.Id).Do()
	if err != nil {
		panic(err)
	}

	var sheetId int64
	for _, sheet := range response.Sheets {
		if sheet.Properties.Title == spreadSheetConfig.SheetName {
			sheetId = sheet.Properties.SheetId
		}
	}

	return service, spreadSheetConfig, sheetId
}

func getSpreadSheetAnalyzeConfig(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, spreadsheet_entity.SpreadSheetAnalyzeConfig, map[spreadsheet_vo.AnalyzeType]*sheets.SheetProperties) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		panic(err)
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		panic(err)
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		panic(err)
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		panic(err)
	}

	var spreadSheetConfig spreadsheet_entity.SpreadSheetAnalyzeConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &spreadSheetConfig); err != nil {
		panic(err)
	}

	response, err := service.Spreadsheets.Get(spreadSheetConfig.Id).Do()
	if err != nil {
		panic(err)
	}

	var sheetReverseMap = map[string]spreadsheet_vo.AnalyzeType{}
	for _, sheetName := range spreadSheetConfig.SheetNames {
		sheetReverseMap[sheetName.Name] = spreadsheet_vo.AnalyzeType(sheetName.Type)
	}

	sheetMap := map[spreadsheet_vo.AnalyzeType]*sheets.SheetProperties{}
	for _, sheet := range response.Sheets {
		if analyzeType, ok := sheetReverseMap[sheet.Properties.Title]; ok {
			sheetMap[analyzeType] = sheet.Properties
		}
	}

	return service, spreadSheetConfig, sheetMap
}

func (s *SpreadSheetClient) WriteForTotalSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A1")
	values := [][]interface{}{
		{
			"累計",
			"",
		},
		{
			"投資",
			summary.GetPayment(),
		},
		{
			"回収",
			summary.GetPayout(),
		},
		{
			"回収率",
			summary.GetRecoveryRate(),
		},
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForCurrentMonthSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "C1")
	values := [][]interface{}{
		{
			"月間累計",
			"",
		},
		{
			"投資",
			summary.GetPayment(),
		},
		{
			"回収",
			summary.GetPayout(),
		},
		{
			"回収率",
			summary.GetRecoveryRate(),
		},
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForCurrentYearSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "E1")
	values := [][]interface{}{
		{
			"年間累計",
			"",
		},
		{
			"投資",
			summary.GetPayment(),
		},
		{
			"回収",
			summary.GetPayout(),
		},
		{
			"回収率",
			summary.GetRecoveryRate(),
		},
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForTotalBettingTicketRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetBettingTicketSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A6")
	values := [][]interface{}{
		{
			"券種別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	winSummary := summary.GetWinSummary()
	placeSummary := summary.GetPlaceSummary()
	quinellaSummary := summary.GetQuinellaSummary()
	exactaSummary := summary.GetExactaSummary()
	quinellaPlaceSummary := summary.GetQuinellaPlaceSummary()
	trioSummary := summary.GetTrioSummary()
	trifectaSummary := summary.GetTrifectaSummary()
	totalSummary := summary.GetTotalSummary()

	values = append(values, []interface{}{
		betting_ticket_vo.Win.Name(),
		winSummary.GetRaceCount(),
		winSummary.GetBetCount(),
		winSummary.GetHitCount(),
		winSummary.GetHitRate(),
		winSummary.GetPayment(),
		winSummary.GetPayout(),
		winSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.Place.Name(),
		placeSummary.GetRaceCount(),
		placeSummary.GetBetCount(),
		placeSummary.GetHitCount(),
		placeSummary.GetHitRate(),
		placeSummary.GetPayment(),
		placeSummary.GetPayout(),
		placeSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.Quinella.Name(),
		quinellaSummary.GetRaceCount(),
		quinellaSummary.GetBetCount(),
		quinellaSummary.GetHitCount(),
		quinellaSummary.GetHitRate(),
		quinellaSummary.GetPayment(),
		quinellaSummary.GetPayout(),
		quinellaSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.Exacta.Name(),
		exactaSummary.GetRaceCount(),
		exactaSummary.GetBetCount(),
		exactaSummary.GetHitCount(),
		exactaSummary.GetHitRate(),
		exactaSummary.GetPayment(),
		exactaSummary.GetPayout(),
		exactaSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.QuinellaPlace.Name(),
		quinellaPlaceSummary.GetRaceCount(),
		quinellaPlaceSummary.GetBetCount(),
		quinellaPlaceSummary.GetHitCount(),
		quinellaPlaceSummary.GetHitRate(),
		quinellaPlaceSummary.GetPayment(),
		quinellaPlaceSummary.GetPayout(),
		quinellaPlaceSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.Trio.Name(),
		trioSummary.GetRaceCount(),
		trioSummary.GetBetCount(),
		trioSummary.GetHitCount(),
		trioSummary.GetHitRate(),
		trioSummary.GetPayment(),
		trioSummary.GetPayout(),
		trioSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		betting_ticket_vo.Trifecta.Name(),
		trifectaSummary.GetRaceCount(),
		trifectaSummary.GetBetCount(),
		trifectaSummary.GetHitCount(),
		trifectaSummary.GetHitRate(),
		trifectaSummary.GetPayment(),
		trifectaSummary.GetPayout(),
		trifectaSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		"累計",
		totalSummary.GetRaceCount(),
		totalSummary.GetBetCount(),
		totalSummary.GetHitCount(),
		totalSummary.GetHitRate(),
		totalSummary.GetPayment(),
		totalSummary.GetPayout(),
		totalSummary.GetRecoveryRate(),
	})

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForRaceClassRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetClassSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A15")

	values := [][]interface{}{
		{
			"クラス別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	grade1Summary := summary.GetGrade1Summary()
	grade2Summary := summary.GetGrade2Summary()
	grade3Summary := summary.GetGrade3Summary()
	jpn1Summary := summary.GetJpn1Summary()
	jpn2Summary := summary.GetJpn2Summary()
	jpn3Summary := summary.GetJpn3Summary()
	openClassSummary := summary.GetOpenClassSummary()
	threeWinClassSummary := summary.GetThreeWinClassSummary()
	twoWinClassSummary := summary.GetTwoWinClassSummary()
	oneWinClassSummary := summary.GetOneWinClassSummary()
	maidenClassSummary := summary.GetMaidenClassSummary()
	makeDebutClassSummary := summary.GetMakeDebutClassSummary()

	values = append(values, []interface{}{
		race_vo.Grade1.String(),
		grade1Summary.GetRaceCount(),
		grade1Summary.GetBetCount(),
		grade1Summary.GetHitCount(),
		grade1Summary.GetHitRate(),
		grade1Summary.GetPayment(),
		grade1Summary.GetPayout(),
		grade1Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Grade2.String(),
		grade2Summary.GetRaceCount(),
		grade2Summary.GetBetCount(),
		grade2Summary.GetHitCount(),
		grade2Summary.GetHitRate(),
		grade2Summary.GetPayment(),
		grade2Summary.GetPayout(),
		grade2Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Grade3.String(),
		grade3Summary.GetRaceCount(),
		grade3Summary.GetBetCount(),
		grade3Summary.GetHitCount(),
		grade3Summary.GetHitRate(),
		grade3Summary.GetPayment(),
		grade3Summary.GetPayout(),
		grade3Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Jpn1.String(),
		jpn1Summary.GetRaceCount(),
		jpn1Summary.GetBetCount(),
		jpn1Summary.GetHitCount(),
		jpn1Summary.GetHitRate(),
		jpn1Summary.GetPayment(),
		jpn1Summary.GetPayout(),
		jpn1Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Jpn2.String(),
		jpn2Summary.GetRaceCount(),
		jpn2Summary.GetBetCount(),
		jpn2Summary.GetHitCount(),
		jpn2Summary.GetHitRate(),
		jpn2Summary.GetPayment(),
		jpn2Summary.GetPayout(),
		jpn2Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Jpn3.String(),
		jpn3Summary.GetRaceCount(),
		jpn3Summary.GetBetCount(),
		jpn3Summary.GetHitCount(),
		jpn3Summary.GetHitRate(),
		jpn3Summary.GetPayment(),
		jpn3Summary.GetPayout(),
		jpn3Summary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.OpenClass.String(),
		openClassSummary.GetRaceCount(),
		openClassSummary.GetBetCount(),
		openClassSummary.GetHitCount(),
		openClassSummary.GetHitRate(),
		openClassSummary.GetPayment(),
		openClassSummary.GetPayout(),
		openClassSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.ThreeWinClass.String(),
		threeWinClassSummary.GetRaceCount(),
		threeWinClassSummary.GetBetCount(),
		threeWinClassSummary.GetHitCount(),
		threeWinClassSummary.GetHitRate(),
		threeWinClassSummary.GetPayment(),
		threeWinClassSummary.GetPayout(),
		threeWinClassSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.TwoWinClass.String(),
		twoWinClassSummary.GetRaceCount(),
		twoWinClassSummary.GetBetCount(),
		twoWinClassSummary.GetHitCount(),
		twoWinClassSummary.GetHitRate(),
		twoWinClassSummary.GetPayment(),
		twoWinClassSummary.GetPayout(),
		twoWinClassSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.OneWinClass.String(),
		oneWinClassSummary.GetRaceCount(),
		oneWinClassSummary.GetBetCount(),
		oneWinClassSummary.GetHitCount(),
		oneWinClassSummary.GetHitRate(),
		oneWinClassSummary.GetPayment(),
		oneWinClassSummary.GetPayout(),
		oneWinClassSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Maiden.String(),
		maidenClassSummary.GetRaceCount(),
		maidenClassSummary.GetBetCount(),
		maidenClassSummary.GetHitCount(),
		maidenClassSummary.GetHitRate(),
		maidenClassSummary.GetPayment(),
		maidenClassSummary.GetPayout(),
		maidenClassSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.MakeDebut.String(),
		makeDebutClassSummary.GetRaceCount(),
		makeDebutClassSummary.GetBetCount(),
		makeDebutClassSummary.GetHitCount(),
		makeDebutClassSummary.GetHitRate(),
		makeDebutClassSummary.GetPayment(),
		makeDebutClassSummary.GetPayout(),
		makeDebutClassSummary.GetRecoveryRate(),
	})

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForCourseCategoryRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetCourseCategorySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "I6")

	values := [][]interface{}{
		{
			"路面別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	turfSummary := summary.GetCourseCategorySummary(race_vo.Turf)
	dirtSummary := summary.GetCourseCategorySummary(race_vo.Dirt)
	jumpSummary := summary.GetCourseCategorySummary(race_vo.Jump)

	values = append(values, []interface{}{
		race_vo.Turf.String(),
		turfSummary.GetRaceCount(),
		turfSummary.GetBetCount(),
		turfSummary.GetHitCount(),
		turfSummary.GetHitRate(),
		turfSummary.GetPayment(),
		turfSummary.GetPayout(),
		turfSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Dirt.String(),
		dirtSummary.GetRaceCount(),
		dirtSummary.GetBetCount(),
		dirtSummary.GetHitCount(),
		dirtSummary.GetHitRate(),
		dirtSummary.GetPayment(),
		dirtSummary.GetPayout(),
		dirtSummary.GetRecoveryRate(),
	})
	values = append(values, []interface{}{
		race_vo.Jump.String(),
		jumpSummary.GetRaceCount(),
		jumpSummary.GetBetCount(),
		jumpSummary.GetHitCount(),
		jumpSummary.GetHitRate(),
		jumpSummary.GetPayment(),
		jumpSummary.GetPayout(),
		jumpSummary.GetRecoveryRate(),
	})

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForDistanceCategoryRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetDistanceCategorySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "I10")

	values := [][]interface{}{
		{
			"距離別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	distanceCategories := []race_vo.DistanceCategory{
		race_vo.TurfSprint,
		race_vo.TurfMile,
		race_vo.TurfIntermediate,
		race_vo.TurfLong,
		race_vo.TurfExtended,
		race_vo.DirtSprint,
		race_vo.DirtMile,
		race_vo.DirtIntermediate,
		race_vo.DirtLong,
		race_vo.JumpAllDistance,
	}

	for _, distanceCategory := range distanceCategories {
		distanceCategorySummary := summary.GetDistanceCategorySummary(distanceCategory)
		values = append(values, []interface{}{
			distanceCategory.String(),
			distanceCategorySummary.GetRaceCount(),
			distanceCategorySummary.GetBetCount(),
			distanceCategorySummary.GetHitCount(),
			distanceCategorySummary.GetHitRate(),
			distanceCategorySummary.GetPayment(),
			distanceCategorySummary.GetPayout(),
			distanceCategorySummary.GetRecoveryRate(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForRaceCourseRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetRaceCourseSummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "I21")

	values := [][]interface{}{
		{
			"開催別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	raceCourses := []race_vo.RaceCourse{
		race_vo.Sapporo, race_vo.Hakodate, race_vo.Fukushima, race_vo.Niigata, race_vo.Tokyo, race_vo.Nakayama, race_vo.Chukyo, race_vo.Kyoto, race_vo.Hanshin, race_vo.Kokura,
		race_vo.Monbetsu, race_vo.Morioka, race_vo.Urawa, race_vo.Hunabashi, race_vo.Ooi, race_vo.Kawasaki, race_vo.Kanazawa, race_vo.Nagoya, race_vo.Sonoda, race_vo.Kouchi, race_vo.Saga,
		race_vo.Overseas,
	}

	for _, raceCourse := range raceCourses {
		raceCourseSummary := summary.GetRaceCourseSummary(raceCourse)
		values = append(values, []interface{}{
			raceCourse.Name(),
			raceCourseSummary.GetRaceCount(),
			raceCourseSummary.GetBetCount(),
			raceCourseSummary.GetHitCount(),
			raceCourseSummary.GetHitRate(),
			raceCourseSummary.GetPayment(),
			raceCourseSummary.GetPayout(),
			raceCourseSummary.GetRecoveryRate(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteForMonthlyRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetMonthlySummary) error {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A28")

	values := [][]interface{}{
		{
			"月別",
			"投票レース数",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
		},
	}

	monthlySummaryMap := summary.GetMonthlySummaryMap()

	var (
		dateList  []int
		summaries []result_summary_entity.DetailSummary
	)
	for key := range monthlySummaryMap {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] > dateList[j]
	})

	for _, date := range dateList {
		summaries = append(summaries, monthlySummaryMap[date])
	}
	for idx, monthlySummary := range summaries {
		values = append(values, []interface{}{
			strconv.Itoa(dateList[idx]),
			monthlySummary.GetRaceCount(),
			monthlySummary.GetBetCount(),
			monthlySummary.GetHitCount(),
			monthlySummary.GetHitRate(),
			monthlySummary.GetPayment(),
			monthlySummary.GetPayout(),
			monthlySummary.GetRecoveryRate(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForTotalSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   2,
						EndRowIndex:      1,
					},
				},
			},
			// 1列目のセル幅調整
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    s.sheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 95,
					},
					Fields: "pixelSize",
				},
			},
			// 2列目のセル幅調整
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   2,
						SheetId:    s.sheetId,
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    1,
						EndColumnIndex:   1,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   2,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForCurrentMonthlySummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 2,
						StartRowIndex:    0,
						EndColumnIndex:   4,
						EndRowIndex:      1,
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    s.sheetId,
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 2,
						StartRowIndex:    1,
						EndColumnIndex:   3,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 2,
						StartRowIndex:    0,
						EndColumnIndex:   4,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForCurrentYearSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			// 1行目のセルをマージ
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 4,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      1,
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    s.sheetId,
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 4,
						StartRowIndex:    1,
						EndColumnIndex:   5,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 4,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      4,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForTotalBettingTicketRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    5,
						EndColumnIndex:   8,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    6,
						EndColumnIndex:   1,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    5,
						EndColumnIndex:   8,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    6,
						EndColumnIndex:   1,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForRaceClassRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    14,
						EndColumnIndex:   8,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      27,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    14,
						EndColumnIndex:   8,
						EndRowIndex:      15,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    15,
						EndColumnIndex:   1,
						EndRowIndex:      27,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForCourseCategoryRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   16,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    6,
						EndColumnIndex:   9,
						EndRowIndex:      9,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   16,
						EndRowIndex:      6,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    5,
						EndColumnIndex:   9,
						EndRowIndex:      9,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForDistanceCategoryRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    9,
						EndColumnIndex:   16,
						EndRowIndex:      10,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    10,
						EndColumnIndex:   9,
						EndRowIndex:      20,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    9,
						EndColumnIndex:   16,
						EndRowIndex:      10,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    10,
						EndColumnIndex:   9,
						EndRowIndex:      20,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForRaceCourseRateSummary(ctx context.Context) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    20,
						EndColumnIndex:   16,
						EndRowIndex:      21,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    21,
						EndColumnIndex:   9,
						EndRowIndex:      43,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    20,
						EndColumnIndex:   16,
						EndRowIndex:      21,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 8,
						StartRowIndex:    21,
						EndColumnIndex:   9,
						EndRowIndex:      43,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) WriteStyleForMonthlyRateSummary(ctx context.Context, rowCount int) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    27,
						EndColumnIndex:   8,
						EndRowIndex:      28,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.6,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    28,
						EndColumnIndex:   1,
						EndRowIndex:      28 + int64(rowCount),
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    27,
						EndColumnIndex:   8,
						EndRowIndex:      28,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    28,
						EndColumnIndex:   1,
						EndRowIndex:      28 + int64(rowCount),
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetMonthlyBettingTicketClient) Write(ctx context.Context, summary *spreadsheet_entity.SpreadSheetMonthlyBettingTicketSummary) error {
	summaryMap := summary.GetMonthlyBettingTicketSummaryMap()
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A1")
	var dateList []int
	for key := range summaryMap {
		dateList = append(dateList, key)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] > dateList[j]
	})

	// 単勝
	values := [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}
	summaries := make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetWinSummary())
	}
	for idx, winSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "単勝"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			winSummary.GetBetCount(),
			winSummary.GetHitCount(),
			winSummary.GetHitRate(),
			winSummary.GetPayment(),
			winSummary.GetPayout(),
			winSummary.GetRecoveryRate(),
			winSummary.GetAveragePayout(),
			winSummary.GetMaxPayout(),
			winSummary.GetMinPayout(),
		})
	}

	// 複勝
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetPlaceSummary())
	}
	for idx, placeSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "複勝"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			placeSummary.GetBetCount(),
			placeSummary.GetHitCount(),
			placeSummary.GetHitRate(),
			placeSummary.GetPayment(),
			placeSummary.GetPayout(),
			placeSummary.GetRecoveryRate(),
			placeSummary.GetAveragePayout(),
			placeSummary.GetMaxPayout(),
			placeSummary.GetMinPayout(),
		})
	}

	// 馬連
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetQuinellaSummary())
	}
	for idx, quinellaSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "馬連"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			quinellaSummary.GetBetCount(),
			quinellaSummary.GetHitCount(),
			quinellaSummary.GetHitRate(),
			quinellaSummary.GetPayment(),
			quinellaSummary.GetPayout(),
			quinellaSummary.GetRecoveryRate(),
			quinellaSummary.GetAveragePayout(),
			quinellaSummary.GetMaxPayout(),
			quinellaSummary.GetMinPayout(),
		})
	}

	// 馬単
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetExactaSummary())
	}
	for idx, exactaSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "馬単"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			exactaSummary.GetBetCount(),
			exactaSummary.GetHitCount(),
			exactaSummary.GetHitRate(),
			exactaSummary.GetPayment(),
			exactaSummary.GetPayout(),
			exactaSummary.GetRecoveryRate(),
			exactaSummary.GetAveragePayout(),
			exactaSummary.GetMaxPayout(),
			exactaSummary.GetMinPayout(),
		})
	}

	// ワイド
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetQuinellaPlaceSummary())
	}
	for idx, quinellaPlaceSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "ワイド"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			quinellaPlaceSummary.GetBetCount(),
			quinellaPlaceSummary.GetHitCount(),
			quinellaPlaceSummary.GetHitRate(),
			quinellaPlaceSummary.GetPayment(),
			quinellaPlaceSummary.GetPayout(),
			quinellaPlaceSummary.GetRecoveryRate(),
			quinellaPlaceSummary.GetAveragePayout(),
			quinellaPlaceSummary.GetMaxPayout(),
			quinellaPlaceSummary.GetMinPayout(),
		})
	}

	// 3連複
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetTrioSummary())
	}
	for idx, trioPlaceSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "3連複"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			trioPlaceSummary.GetBetCount(),
			trioPlaceSummary.GetHitCount(),
			trioPlaceSummary.GetHitRate(),
			trioPlaceSummary.GetPayment(),
			trioPlaceSummary.GetPayout(),
			trioPlaceSummary.GetRecoveryRate(),
			trioPlaceSummary.GetAveragePayout(),
			trioPlaceSummary.GetMaxPayout(),
			trioPlaceSummary.GetMinPayout(),
		})
	}

	// 3連単
	values = append(values, [][]interface{}{
		{
			"",
			"月別",
			"投票回数",
			"的中回数",
			"的中率",
			"投資額",
			"回収額",
			"回収率",
			"平均払戻金額",
			"最大払戻金額",
			"最小払戻金額",
		},
	}...)
	summaries = make([]result_summary_entity.DetailSummary, 0)
	for _, date := range dateList {
		summaries = append(summaries, summaryMap[date].GetTrifectaSummary())
	}
	for idx, trifectaPlaceSummary := range summaries {
		headerColumn := ""
		if idx == 0 {
			headerColumn = "3連単"
		}
		values = append(values, []interface{}{
			headerColumn,
			strconv.Itoa(dateList[idx]),
			trifectaPlaceSummary.GetBetCount(),
			trifectaPlaceSummary.GetHitCount(),
			trifectaPlaceSummary.GetHitRate(),
			trifectaPlaceSummary.GetPayment(),
			trifectaPlaceSummary.GetPayout(),
			trifectaPlaceSummary.GetRecoveryRate(),
			trifectaPlaceSummary.GetAveragePayout(),
			trifectaPlaceSummary.GetMaxPayout(),
			trifectaPlaceSummary.GetMinPayout(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetMonthlyBettingTicketClient) WriteStyle(ctx context.Context, rowCount int) error {
	endRowFunc := func(idx int) int64 {
		r := (rowCount + 1) * idx
		return int64(r)
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(0) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(1),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(1) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(2),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(2) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(3),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(3) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(4),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(4) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(5),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(5) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(6),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(6) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(0) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(1),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(1) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(2),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(2) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(3),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(3) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(4),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(4) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(5),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(5) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    endRowFunc(6),
						EndColumnIndex:   11,
						EndRowIndex:      endRowFunc(6) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    s.sheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 45,
					},
					Fields: "pixelSize",
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetListClient) WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle, error) {
	writeRange := fmt.Sprintf("%s!%s", s.spreadSheetConfig.SheetName, "A1")
	values := [][]interface{}{
		{
			"レース条件",
			"",
			"",
			"",
			"",
			"",
			"投資額",
			"回収額",
			"回収率",
			"本命",
			"騎手",
			"人気",
			"オッズ",
			"対抗",
			"騎手",
			"人気",
			"オッズ",
			"1着",
			"騎手",
			"人気",
			"オッズ",
			"2着",
			"騎手",
			"人気",
			"オッズ",
		},
	}

	var rivalHorseName, rivalPopularNumber, rivalOdds, rivalJockeyName string

	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race().StartTime() > records[j].Race().StartTime()
	})
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Race().RaceDate() > records[j].Race().RaceDate()
	})

	styleMap := map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle{}
	for idx, record := range records {
		var (
			favoriteColor, rivalColor         spreadsheet_vo.PlaceColor
			firstPlaceColor, secondPlaceColor spreadsheet_vo.PopularColor
			gradeClassColor                   spreadsheet_vo.GradeClassColor
			repaymentComments                 spreadsheet_vo.RepaymentComments
		)
		raceResults := record.Race().RaceResults()
		raceResultOfFirst := raceResults[0]
		raceResultOfSecond := raceResults[1]
		raceResultOfThird := raceResults[2]

		if record.FavoriteHorse().HorseName() == raceResultOfFirst.HorseName() {
			favoriteColor = spreadsheet_vo.FirstPlace
		} else if record.FavoriteHorse().HorseName() == raceResultOfSecond.HorseName() {
			favoriteColor = spreadsheet_vo.SecondPlace
		} else if record.FavoriteHorse().HorseName() == raceResultOfThird.HorseName() {
			favoriteColor = spreadsheet_vo.ThirdPlace
		}

		if record.RivalHorse() != nil {
			if record.RivalHorse().HorseName() == raceResultOfFirst.HorseName() {
				rivalColor = spreadsheet_vo.FirstPlace
			} else if record.RivalHorse().HorseName() == raceResultOfSecond.HorseName() {
				rivalColor = spreadsheet_vo.SecondPlace
			} else if record.RivalHorse().HorseName() == raceResultOfThird.HorseName() {
				rivalColor = spreadsheet_vo.ThirdPlace
			}
		}

		if raceResultOfFirst.PopularNumber() == 1 {
			firstPlaceColor = spreadsheet_vo.FirstPopular
		} else if raceResultOfFirst.PopularNumber() == 2 {
			firstPlaceColor = spreadsheet_vo.SecondPopular
		} else if raceResultOfFirst.PopularNumber() == 3 {
			firstPlaceColor = spreadsheet_vo.ThirdPopular
		}

		if raceResultOfSecond.PopularNumber() == 1 {
			secondPlaceColor = spreadsheet_vo.FirstPopular
		} else if raceResultOfSecond.PopularNumber() == 2 {
			secondPlaceColor = spreadsheet_vo.SecondPopular
		} else if raceResultOfSecond.PopularNumber() == 3 {
			secondPlaceColor = spreadsheet_vo.ThirdPopular
		}

		switch record.Race().Class() {
		case race_vo.Grade1, race_vo.Jpn1, race_vo.JumpGrade1:
			gradeClassColor = spreadsheet_vo.Grade1
		case race_vo.Grade2, race_vo.Jpn2, race_vo.JumpGrade2:
			gradeClassColor = spreadsheet_vo.Grade2
		case race_vo.Grade3, race_vo.Jpn3, race_vo.JumpGrade3:
			gradeClassColor = spreadsheet_vo.Grade3
		}

		if record.WinningTickets() != nil {
			for _, winningTicket := range record.WinningTickets() {
				repaymentComments = append(repaymentComments,
					fmt.Sprintf("%s %s %s倍 %d円 %d人気", winningTicket.BettingTicket().Name(), winningTicket.BetNumber().String(), winningTicket.Odds(), winningTicket.Repayment(), winningTicket.Popular()))
			}
		}

		styleMap[record.Race().RaceId()] = spreadsheet_entity.NewSpreadSheetStyle(
			idx+1, favoriteColor, rivalColor, firstPlaceColor, secondPlaceColor, gradeClassColor, repaymentComments,
		)

		if record.RivalHorse() != nil {
			rivalHorseName = record.RivalHorse().HorseName()
			rivalPopularNumber = strconv.Itoa(record.RivalHorse().PopularNumber())
			rivalOdds = record.RivalHorse().Odds()
			rivalJockeyName = record.RivalJockey().JockeyName()
		} else {
			rivalHorseName = "-"
			rivalPopularNumber = "-"
			rivalOdds = "-"
			rivalJockeyName = "-"
		}

		values = append(values, []interface{}{
			record.Race().RaceDate().DateFormat(),
			record.Race().Class().String(),
			record.Race().CourseCategory().String(),
			fmt.Sprintf("%d%s", record.Race().Distance(), "m"),
			record.Race().TrackCondition(),
			fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", record.Race().Url(), record.Race().RaceName()),
			record.Payment(),
			record.Repayment(),
			fmt.Sprintf("%.0f%s", float64(record.Repayment())*float64(100)/float64(record.Payment()), "%"),
			record.FavoriteHorse().HorseName(),
			record.FavoriteJockey().JockeyName(),
			record.FavoriteHorse().PopularNumber(),
			record.FavoriteHorse().Odds(),
			rivalHorseName,
			rivalJockeyName,
			rivalPopularNumber,
			rivalOdds,
			raceResultOfFirst.HorseName(),
			raceResultOfFirst.JockeyName(),
			raceResultOfFirst.PopularNumber(),
			raceResultOfFirst.Odds(),
			raceResultOfSecond.HorseName(),
			raceResultOfSecond.JockeyName(),
			raceResultOfSecond.PopularNumber(),
			raceResultOfSecond.Odds(),
		})
	}

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return nil, err
	}

	return styleMap, nil
}

func (s *SpreadSheetListClient) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle) error {
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				MergeCells: &sheets.MergeCellsRequest{
					MergeType: "MERGE_ROWS",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   6,
						EndRowIndex:      1,
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   25,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   25,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    s.sheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 80,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   2,
						SheetId:    s.sheetId,
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 90,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   3,
						SheetId:    s.sheetId,
						StartIndex: 2,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 45,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   4,
						SheetId:    s.sheetId,
						StartIndex: 3,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   5,
						SheetId:    s.sheetId,
						StartIndex: 4,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 25,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   6,
						SheetId:    s.sheetId,
						StartIndex: 5,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 130,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   9,
						SheetId:    s.sheetId,
						StartIndex: 6,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 60,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   10,
						SheetId:    s.sheetId,
						StartIndex: 9,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   11,
						SheetId:    s.sheetId,
						StartIndex: 10,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   13,
						SheetId:    s.sheetId,
						StartIndex: 11,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   14,
						SheetId:    s.sheetId,
						StartIndex: 13,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   15,
						SheetId:    s.sheetId,
						StartIndex: 14,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   17,
						SheetId:    s.sheetId,
						StartIndex: 15,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   18,
						SheetId:    s.sheetId,
						StartIndex: 17,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   19,
						SheetId:    s.sheetId,
						StartIndex: 18,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   21,
						SheetId:    s.sheetId,
						StartIndex: 19,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   22,
						SheetId:    s.sheetId,
						StartIndex: 21,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   23,
						SheetId:    s.sheetId,
						StartIndex: 22,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 75,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   26,
						SheetId:    s.sheetId,
						StartIndex: 23,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	var requests []*sheets.Request
	for _, record := range records {
		if style, ok := styleMap[record.Race().RaceId()]; ok {
			if style.GetGradeClassColor() != spreadsheet_vo.NonGrade {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GetGradeClassColor() == spreadsheet_vo.Grade1 {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GetGradeClassColor() == spreadsheet_vo.Grade2 {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GetGradeClassColor() == spreadsheet_vo.Grade3 {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 1,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   2,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.GetFavoriteColor() != spreadsheet_vo.OtherPlace {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GetFavoriteColor() == spreadsheet_vo.FirstPlace {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GetFavoriteColor() == spreadsheet_vo.SecondPlace {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GetFavoriteColor() == spreadsheet_vo.ThirdPlace {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 9,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   10,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.GetRivalColor() != spreadsheet_vo.OtherPlace {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GetRivalColor() == spreadsheet_vo.FirstPlace {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GetRivalColor() == spreadsheet_vo.SecondPlace {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GetRivalColor() == spreadsheet_vo.ThirdPlace {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 13,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   14,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.GetFirstPlaceColor() != spreadsheet_vo.OtherPopular {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GetFirstPlaceColor() == spreadsheet_vo.FirstPopular {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GetFirstPlaceColor() == spreadsheet_vo.SecondPopular {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GetFirstPlaceColor() == spreadsheet_vo.ThirdPopular {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 17,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   18,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if style.GetSecondPlaceColor() != spreadsheet_vo.OtherPopular {
				color := &sheets.Color{
					Red:   1.0,
					Blue:  1.0,
					Green: 1.0,
				}
				if style.GetSecondPlaceColor() == spreadsheet_vo.FirstPopular {
					color = &sheets.Color{
						Red:   1.0,
						Green: 0.937,
						Blue:  0.498,
					}
				} else if style.GetSecondPlaceColor() == spreadsheet_vo.SecondPopular {
					color = &sheets.Color{
						Red:   0.796,
						Green: 0.871,
						Blue:  1.0,
					}
				} else if style.GetSecondPlaceColor() == spreadsheet_vo.ThirdPopular {
					color = &sheets.Color{
						Red:   0.937,
						Green: 0.78,
						Blue:  0.624,
					}
				}

				cellRequest := &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 21,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   22,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: color,
						},
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
			if len(style.GetRepaymentComment()) > 0 {
				cellRequest := &sheets.RepeatCellRequest{
					Fields: "note",
					Range: &sheets.GridRange{
						SheetId:          s.sheetId,
						StartColumnIndex: 7,
						StartRowIndex:    int64(style.GetRowIndex()),
						EndColumnIndex:   8,
						EndRowIndex:      int64(style.GetRowIndex()) + 1,
					},
					Cell: &sheets.CellData{
						Note: style.GetRepaymentComment(),
					},
				}
				requests = append(requests, &sheets.Request{
					RepeatCell: cellRequest,
				})
			}
		}
	}

	_, err = s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetAnalyzeClient) WriteWinPopular(ctx context.Context, summary *analyze_entity.WinAnalyzeSummary) error {
	sheetProperties, ok := s.sheetMap[spreadsheet_vo.Win]
	if !ok {
		return fmt.Errorf("sheet not found")
	}

	writeRange := fmt.Sprintf("%s!%s", sheetProperties.Title, "A1")
	var values [][]interface{}
	addValues := func(values [][]interface{}, summaries []*analyze_entity.WinPopularAnalyzeSummary, title string) [][]interface{} {
		values = append(values, [][]interface{}{
			{
				title,
				"1人気",
				"2人気",
				"3人気",
				"4人気",
				"5人気",
				"6人気",
				"7人気",
				"8人気",
				"9人気",
				"10人気",
				"11人気",
				"12人気",
				"13人気",
				"14人気",
				"15人気",
				"16人気",
				"17人気",
				"18人気",
			},
		}...)
		capacity := 19
		betCounts := make([]interface{}, 0, capacity)
		hitCounts := make([]interface{}, 0, capacity)
		hitRates := make([]interface{}, 0, capacity)
		payoutRates := make([]interface{}, 0, capacity)
		averagePayoutRateAtHits := make([]interface{}, 0, capacity)
		medianPayoutRateAtHits := make([]interface{}, 0, capacity)
		payoutUpside := make([]interface{}, 0, capacity)
		averageOddsAtVotes := make([]interface{}, 0, capacity)
		averageOddsAtHits := make([]interface{}, 0, capacity)
		averageOddsAtUnHits := make([]interface{}, 0, capacity)
		maxOddsAtHits := make([]interface{}, 0, capacity)
		minOddsAtHits := make([]interface{}, 0, capacity)
		totalPayments := make([]interface{}, 0, capacity)
		totalPayouts := make([]interface{}, 0, capacity)
		averagePayments := make([]interface{}, 0, capacity)
		averagePayouts := make([]interface{}, 0, capacity)
		medianPayments := make([]interface{}, 0, capacity)
		medianPayouts := make([]interface{}, 0, capacity)
		maxPayouts := make([]interface{}, 0, capacity)
		minPayouts := make([]interface{}, 0, capacity)

		betCounts = append(betCounts, "投票回数")
		hitCounts = append(hitCounts, "的中回数")
		hitRates = append(hitRates, "的中率")
		payoutRates = append(payoutRates, "回収率")
		averagePayoutRateAtHits = append(averagePayoutRateAtHits, "的中回収率(平均値)")
		medianPayoutRateAtHits = append(medianPayoutRateAtHits, "的中回収率(中央値)")
		payoutUpside = append(payoutUpside, "回収上振れ率")
		averageOddsAtVotes = append(averageOddsAtVotes, "投票時平均オッズ")
		averageOddsAtHits = append(averageOddsAtHits, "的中時平均オッズ")
		averageOddsAtUnHits = append(averageOddsAtUnHits, "不的中時平均オッズ")
		maxOddsAtHits = append(maxOddsAtHits, "的中時最大オッズ")
		minOddsAtHits = append(minOddsAtHits, "的中時最小オッズ")
		totalPayments = append(totalPayments, "投票金額合計")
		totalPayouts = append(totalPayouts, "払戻金額合計")
		averagePayments = append(averagePayments, "平均投票金額")
		medianPayments = append(medianPayments, "中央値投票金額")
		averagePayouts = append(averagePayouts, "平均払戻金額")
		medianPayouts = append(medianPayouts, "中央値払戻金額")
		maxPayouts = append(maxPayouts, "最大払戻金額")
		minPayouts = append(minPayouts, "最小払戻金額")

		for _, record := range summaries {
			betCounts = append(betCounts, record.BetCount())
			hitCounts = append(hitCounts, record.HitCount())
			hitRates = append(hitRates, record.FormattedHitRate())
			payoutRates = append(payoutRates, record.FormattedPayoutRate())
			averagePayoutRateAtHits = append(averagePayoutRateAtHits, record.FormattedAveragePayoutRateAtHit())
			medianPayoutRateAtHits = append(medianPayoutRateAtHits, record.FormattedMedianPayoutRateAtHit())
			payoutUpside = append(payoutUpside, record.FormattedPayoutUpsideRate())
			averageOddsAtVotes = append(averageOddsAtVotes, record.AverageOddsAtVote())
			averageOddsAtHits = append(averageOddsAtHits, record.AverageOddsAtHit())
			averageOddsAtUnHits = append(averageOddsAtUnHits, record.AverageOddsAtUnHit())
			maxOddsAtHits = append(maxOddsAtHits, record.MaxOddsAtHit())
			minOddsAtHits = append(minOddsAtHits, record.MinOddsAtHit())
			totalPayments = append(totalPayments, record.TotalPayment())
			totalPayouts = append(totalPayouts, record.TotalPayout())
			averagePayments = append(averagePayments, record.AveragePayment())
			medianPayments = append(medianPayments, record.MedianPayment())
			averagePayouts = append(averagePayouts, record.AveragePayout())
			medianPayouts = append(medianPayouts, record.MedianPayout())
			maxPayouts = append(maxPayouts, record.MaxPayout())
			minPayouts = append(minPayouts, record.MinPayout())
		}

		values = append(values, betCounts, hitCounts, hitRates, payoutRates, averagePayoutRateAtHits, medianPayoutRateAtHits, payoutUpside, averageOddsAtVotes,
			averageOddsAtHits, averageOddsAtUnHits, maxOddsAtHits, minOddsAtHits, totalPayments, totalPayouts, averagePayments, medianPayments, averagePayouts, medianPayouts, maxPayouts, minPayouts)

		return values
	}

	values = addValues(values, summary.AllSummaries(), "全レース集計")
	values = addValues(values, summary.Grade1Summaries(), "JRA G1集計")
	values = addValues(values, summary.Grade2Summaries(), "JRA G2集計")
	values = addValues(values, summary.Grade3Summaries(), "JRA G3集計")
	values = addValues(values, summary.AllowanceClassSummaries(), "JRA 平場集計")

	_, err := s.client.Spreadsheets.Values.Update(s.spreadSheetConfig.Id, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetAnalyzeClient) WriteStyleWinPopular(ctx context.Context, summary *analyze_entity.WinAnalyzeSummary) error {
	sheetProperties, _ := s.sheetMap[spreadsheet_vo.Win]

	// 全レース
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   1,
						EndRowIndex:      110,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   1,
						SheetId:    sheetProperties.SheetId,
						StartIndex: 0,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 135,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						Dimension:  "COLUMNS",
						EndIndex:   19,
						SheetId:    sheetProperties.SheetId,
						StartIndex: 1,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 50,
					},
					Fields: "pixelSize",
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   19,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    0,
						EndColumnIndex:   19,
						EndRowIndex:      1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    1,
						EndColumnIndex:   1,
						EndRowIndex:      8,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    8,
						EndColumnIndex:   1,
						EndRowIndex:      13,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    13,
						EndColumnIndex:   1,
						EndRowIndex:      21,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    3,
						EndColumnIndex:   19,
						EndRowIndex:      8,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    21,
						EndColumnIndex:   19,
						EndRowIndex:      22,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    21,
						EndColumnIndex:   19,
						EndRowIndex:      22,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    22,
						EndColumnIndex:   1,
						EndRowIndex:      29,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    29,
						EndColumnIndex:   1,
						EndRowIndex:      34,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    34,
						EndColumnIndex:   1,
						EndRowIndex:      42,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    26,
						EndColumnIndex:   19,
						EndRowIndex:      29,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    42,
						EndColumnIndex:   19,
						EndRowIndex:      43,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    42,
						EndColumnIndex:   19,
						EndRowIndex:      43,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    43,
						EndColumnIndex:   1,
						EndRowIndex:      50,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    50,
						EndColumnIndex:   1,
						EndRowIndex:      55,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    55,
						EndColumnIndex:   1,
						EndRowIndex:      63,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    45,
						EndColumnIndex:   19,
						EndRowIndex:      50,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    63,
						EndColumnIndex:   19,
						EndRowIndex:      64,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    63,
						EndColumnIndex:   19,
						EndRowIndex:      64,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    64,
						EndColumnIndex:   1,
						EndRowIndex:      71,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    71,
						EndColumnIndex:   1,
						EndRowIndex:      76,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    76,
						EndColumnIndex:   1,
						EndRowIndex:      84,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    66,
						EndColumnIndex:   19,
						EndRowIndex:      71,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.bold",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    84,
						EndColumnIndex:   19,
						EndRowIndex:      85,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								Bold: true,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    84,
						EndColumnIndex:   19,
						EndRowIndex:      85,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    85,
						EndColumnIndex:   1,
						EndRowIndex:      92,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   1.0,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    92,
						EndColumnIndex:   1,
						EndRowIndex:      97,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  0,
								Green: 0.75,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    97,
						EndColumnIndex:   1,
						EndRowIndex:      105,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.35,
								Blue:  1.0,
								Green: 1.0,
							},
						},
					},
				},
			},
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.numberFormat",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 1,
						StartRowIndex:    87,
						EndColumnIndex:   19,
						EndRowIndex:      92,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type:    "PERCENT",
								Pattern: "0%",
							},
						},
					},
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	payoutUpsideIndexMap := map[int64][]int64{}
	noPaymentIndexMap := map[int64][]int64{}
	payoutUpsideRowIndexForAll := int64(7)
	payoutUpsideRowIndexForG1 := int64(28)
	payoutUpsideRowIndexForG2 := int64(49)
	payoutUpsideRowIndexForG3 := int64(70)
	payoutUpsideRowIndexForAllowance := int64(91)
	noPaymentRowIndexForAll := int64(1)
	noPaymentRowIndexForG1 := int64(22)
	noPaymentRowIndexForG2 := int64(43)
	noPaymentRowIndexForG3 := int64(64)
	noPaymentRowIndexForAllowance := int64(85)
	payoutUpsideIndexMap[payoutUpsideRowIndexForAll] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG1] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG2] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForG3] = []int64{}
	payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForAll] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG1] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG2] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForG3] = []int64{}
	noPaymentIndexMap[noPaymentRowIndexForAllowance] = []int64{}

	for idx, record := range summary.AllSummaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForAll] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForAll], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForAll] = append(noPaymentIndexMap[noPaymentRowIndexForAll], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade1Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG1] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG1], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG1] = append(noPaymentIndexMap[noPaymentRowIndexForG1], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade2Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG2] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG2], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG2] = append(noPaymentIndexMap[noPaymentRowIndexForG2], int64(idx+1))
		}
	}
	for idx, record := range summary.Grade3Summaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForG3] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForG3], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForG3] = append(noPaymentIndexMap[noPaymentRowIndexForG3], int64(idx+1))
		}
	}
	for idx, record := range summary.AllowanceClassSummaries() {
		if record.PayoutUpsideRate() < 0 {
			payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance] = append(payoutUpsideIndexMap[payoutUpsideRowIndexForAllowance], int64(idx+1))
		}
		if record.BetCount() == 0 {
			noPaymentIndexMap[noPaymentRowIndexForAllowance] = append(noPaymentIndexMap[noPaymentRowIndexForAllowance], int64(idx+1))
		}
	}

	var requests []*sheets.Request
	// 回収上振れ率の専用style
	for rowIndex, columnIndexes := range payoutUpsideIndexMap {
		for _, columnIndex := range columnIndexes {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.textFormat.foregroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: columnIndex,
						StartRowIndex:    rowIndex,
						EndColumnIndex:   columnIndex + 1,
						EndRowIndex:      rowIndex + 1,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{
								ForegroundColor: &sheets.Color{
									Red:   1.0,
									Blue:  0,
									Green: 0,
								},
							},
						},
					},
				},
			})
		}

	}

	// 購入実績なしのセルはグレーアウト
	for rowIndex, columnIndexes := range noPaymentIndexMap {
		for _, columnIndex := range columnIndexes {
			requests = append(requests, &sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "userEnteredFormat.backgroundColor",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: columnIndex,
						StartRowIndex:    rowIndex,
						EndColumnIndex:   columnIndex + 1,
						EndRowIndex:      rowIndex + 20,
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							BackgroundColor: &sheets.Color{
								Red:   0.8,
								Blue:  0.8,
								Green: 0.8,
							},
						},
					},
				},
			})
		}
	}

	_, err = s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetClient) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.sheetId,
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   16,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetMonthlyBettingTicketClient) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.sheetId,
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   11,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetListClient) Clear(ctx context.Context) error {
	requests := []*sheets.Request{
		{
			RepeatCell: &sheets.RepeatCellRequest{
				Fields: "*",
				Range: &sheets.GridRange{
					SheetId:          s.sheetId,
					StartColumnIndex: 0,
					StartRowIndex:    0,
					EndColumnIndex:   25,
					EndRowIndex:      9999,
				},
				Cell: &sheets.CellData{},
			},
		},
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheetAnalyzeClient) Clear(ctx context.Context) error {
	for _, sheetProperties := range s.sheetMap {
		requests := []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Fields: "*",
					Range: &sheets.GridRange{
						SheetId:          sheetProperties.SheetId,
						StartColumnIndex: 0,
						StartRowIndex:    0,
						EndColumnIndex:   25,
						EndRowIndex:      9999,
					},
					Cell: &sheets.CellData{},
				},
			},
		}
		_, err := s.client.Spreadsheets.BatchUpdate(s.spreadSheetConfig.Id, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Do()

		if err != nil {
			return err
		}
	}

	return nil
}
