package infrastructure

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/umaca_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	ticketUmacaDataSuffix = "_tohyo_umaca"
)

type umacaTicketRepository struct {
	pathOptimizer file_gateway.PathOptimizer
}

func NewUmacaTicketRepository(
	pathOptimizer file_gateway.PathOptimizer,
) repository.UmacaTicketRepository {
	return &umacaTicketRepository{
		pathOptimizer: pathOptimizer,
	}
}

func (u *umacaTicketRepository) GetMaster(
	ctx context.Context,
	path string,
) ([]*umaca_csv_entity.UmacaMaster, error) {
	rootPath, err := u.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var masters []*umaca_csv_entity.UmacaMaster
	reader := csv.NewReader(f) // UTF-8決め打ちする
	rowNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if rowNum == 0 {
			rowNum++
			continue
		}

		master, err := umaca_csv_entity.NewUmacaMaster(
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
		)
		if err != nil {
			return nil, err
		}

		masters = append(masters, master)
		rowNum++
	}

	return masters, nil
}

func (u *umacaTicketRepository) List(ctx context.Context, path string) ([]string, error) {
	rootPath, err := u.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(absPath, "*.csv")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !strings.Contains(file, ticketUmacaDataSuffix) {
			continue
		}
		fileNames = append(fileNames, filepath.Base(file))
	}

	return fileNames, nil
}

func (u *umacaTicketRepository) Write(
	ctx context.Context,
	path string,
	data [][]string,
) error {
	rootPath, err := u.pathOptimizer.GetProjectRoot()
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(transform.NewWriter(f, japanese.ShiftJIS.NewEncoder()))
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}
