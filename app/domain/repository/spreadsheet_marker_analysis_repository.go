package repository

import "context"

type SpreadSheetMarkerAnalysisRepository interface {
	Write(ctx context.Context) error
}
