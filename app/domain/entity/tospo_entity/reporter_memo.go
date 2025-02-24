package tospo_entity

import (
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type ReporterMemo struct {
	horseNumber types.HorseNumber
	memos       []*Memo
}

type Memo struct {
	comment string
	date    time.Time
}

func NewReporterMemo(
	horseNumber types.HorseNumber,
	memos []*Memo,
) *ReporterMemo {
	return &ReporterMemo{
		horseNumber: horseNumber,
		memos:       memos,
	}
}

func (r *ReporterMemo) HorseNumber() types.HorseNumber {
	return r.horseNumber
}

func (r *ReporterMemo) Memos() []*Memo {
	return r.memos
}

func NewMemo(
	comment string,
	rawDate string,
) (*Memo, error) {
	layout := "2006-01-02" // このレイアウトは固定
	dateTime, err := time.Parse(layout, rawDate)
	if err != nil {
		return nil, err
	}

	return &Memo{
		comment: comment,
		date:    dateTime,
	}, nil
}

func (m *Memo) Comment() string {
	return m.comment
}

func (m *Memo) Date() time.Time {
	return m.date
}
