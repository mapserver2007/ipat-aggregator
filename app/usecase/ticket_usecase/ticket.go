package ticket_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
	"strings"
)

type TicketUseCase struct {
	csvRepository repository.TicketCsvRepository
}

func NewTicketUseCase(
	csvRepository repository.TicketCsvRepository,
) *TicketUseCase {
	return &TicketUseCase{
		csvRepository: csvRepository,
	}
}

func (t *TicketUseCase) Read(ctx context.Context) ([]*ticket_csv_entity.Ticket, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv")
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var allTickets []*ticket_csv_entity.Ticket
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if filepath.Ext(filePath) != ".csv" {
			continue
		}
		if !strings.Contains(filePath, "_tohyo") {
			continue
		}

		tickets, err := t.csvRepository.Read(ctx, filePath)
		if err != nil {
			return nil, err
		}
		allTickets = append(allTickets, tickets...)
	}

	return allTickets, nil
}
