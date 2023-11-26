package ticket_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
)

type ticket struct {
	csvRepository repository.TicketCsvRepository
}

func NewTicket(
	csvRepository repository.TicketCsvRepository,
) *ticket {
	return &ticket{
		csvRepository: csvRepository,
	}
}

func (t *ticket) Read(ctx context.Context) error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv")
	if err != nil {
		return err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	//var results []*ticket_csv_entity.Ticket
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if filepath.Ext(filePath) != ".csv" {
			continue
		}
		t.csvRepository.Read(ctx, filePath)

	}

	return nil
}

func (t *ticket) Write() error {
	//TODO implement me
	panic("implement me")
}
