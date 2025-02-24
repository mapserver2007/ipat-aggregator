package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const secretFileName = "secret.json"

func getSpreadSheetConfig(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, *spreadsheet_entity.SpreadSheetConfig, error) {
	rootPath, err := getProjectRoot()
	if err != nil {
		return nil, nil, err
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		return nil, nil, err
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		return nil, nil, err
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		return nil, nil, err
	}

	var rawSpreadSheetConfig raw_entity.SpreadSheetConfig
	if err = json.Unmarshal(spreadSheetConfigBytes, &rawSpreadSheetConfig); err != nil {
		return nil, nil, err
	}

	response, err := service.Spreadsheets.Get(rawSpreadSheetConfig.Id).Do()
	if err != nil {
		return nil, nil, err
	}

	var spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
	for _, sheet := range response.Sheets {
		if sheet.Properties.Title == rawSpreadSheetConfig.SheetName {
			spreadSheetConfig = spreadsheet_entity.NewSpreadSheetConfig(rawSpreadSheetConfig.Id, sheet.Properties.SheetId, sheet.Properties.Title)
		}
	}

	return service, spreadSheetConfig, nil
}

func getSpreadSheetConfigs(
	ctx context.Context,
	spreadSheetConfigFileName string,
) (*sheets.Service, []*spreadsheet_entity.SpreadSheetConfig, error) {
	rootPath, err := getProjectRoot()
	if err != nil {
		return nil, nil, err
	}

	secretFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, secretFileName))
	if err != nil {
		return nil, nil, err
	}
	spreadSheetConfigFilePath, err := filepath.Abs(fmt.Sprintf("%s/secret/%s", rootPath, spreadSheetConfigFileName))
	if err != nil {
		return nil, nil, err
	}

	credential := option.WithCredentialsFile(secretFilePath)
	service, err := sheets.NewService(ctx, credential)
	if err != nil {
		return nil, nil, err
	}

	spreadSheetConfigBytes, err := os.ReadFile(spreadSheetConfigFilePath)
	if err != nil {
		return nil, nil, err
	}

	var rawSpreadSheetConfigs raw_entity.SpreadSheetConfigs
	if err = json.Unmarshal(spreadSheetConfigBytes, &rawSpreadSheetConfigs); err != nil {
		return nil, nil, err
	}

	response, err := service.Spreadsheets.Get(rawSpreadSheetConfigs.Id).Do()
	if err != nil {
		return nil, nil, err
	}

	var spreadSheetConfigs []*spreadsheet_entity.SpreadSheetConfig
	for _, sheet := range response.Sheets {
		for _, sheetName := range rawSpreadSheetConfigs.SheetNames {
			if sheet.Properties.Title == sheetName {
				spreadSheetConfigs = append(spreadSheetConfigs, spreadsheet_entity.NewSpreadSheetConfig(rawSpreadSheetConfigs.Id, sheet.Properties.SheetId, sheet.Properties.Title))
			}
		}
	}

	return service, spreadSheetConfigs, nil
}

const (
	targetFileName = "go.mod"
)

// FIXME di化
func getProjectRoot() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)

	for {
		modPath := filepath.Join(execDir, targetFileName)
		if _, err := os.Stat(modPath); err == nil {
			return execDir, nil
		}

		parentDir := filepath.Dir(execDir)
		if parentDir == execDir {
			break
		}

		execDir = parentDir
	}

	rootPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return rootPath, nil
}
