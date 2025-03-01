package gateway

import (
	"context"
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetAnalysisPlaceUnhitFileName = "spreadsheet_analysis_place_unhit.json"
)

type SpreadSheetAnalysisPlaceUnhitGateway interface {
	Write(ctx context.Context, analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit) error
}

type spreadSheetAnalysisPlaceUnhitGateway struct {
	spreadSheetConfigGateway SpreadSheetConfigGateway
	logger                   *logrus.Logger
}

func NewSpreadSheetAnalysisPlaceUnhitGateway(
	spreadSheetConfigGateway SpreadSheetConfigGateway,
	logger *logrus.Logger,
) SpreadSheetAnalysisPlaceUnhitGateway {
	return &spreadSheetAnalysisPlaceUnhitGateway{
		spreadSheetConfigGateway: spreadSheetConfigGateway,
		logger:                   logger,
	}
}

func (g *spreadSheetAnalysisPlaceUnhitGateway) Write(
	ctx context.Context,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	client, config, err := g.spreadSheetConfigGateway.GetConfig(ctx, spreadSheetAnalysisPlaceUnhitFileName)
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit start")
	writeRange := fmt.Sprintf("%s!%s", config.SheetName(), "A1")
	values := [][]any{
		{
			"レース条件",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"頭数",
			"馬番",
			"馬名",
			"騎手",
			"人",
			"着",
		},
	}

	for _, analysisPlaceUnhit := range analysisPlaceUnhits {
		values = append(values, []any{
			analysisPlaceUnhit.RaceDate(),
			analysisPlaceUnhit.Class().String(),
			analysisPlaceUnhit.RaceName(),
			analysisPlaceUnhit.CourseCategory().String(),
			analysisPlaceUnhit.Distance(),
			analysisPlaceUnhit.RaceWeightCondition().String(),
			analysisPlaceUnhit.TrackCondition().String(),
			analysisPlaceUnhit.Entries(),
			analysisPlaceUnhit.HorseNumber(),
			analysisPlaceUnhit.HorseName(),
			analysisPlaceUnhit.JockeyName(),
			analysisPlaceUnhit.PopularNumber(),
			analysisPlaceUnhit.Odds().String(),
			analysisPlaceUnhit.OrderNo(),
		})
	}

	_, err = client.Spreadsheets.Values.Update(config.SpreadSheetId(), writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}

	g.logger.Infof("write analysis place unhit end")

	return nil
}
