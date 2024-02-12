package list_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type ListRow struct {
	race           *Race
	favoriteHorse  *Horse
	rivalHorse     *Horse
	favoriteJockey *Jockey
	rivalJockey    *Jockey
	payment        types.Payment
	payout         types.Payout
	hitTickets     []*Ticket
	status         types.PredictStatus
}

func NewListRow(
	race *Race,
	favoriteHorse *Horse,
	rivalHorse *Horse,
	favoriteJockey *Jockey,
	rivalJockey *Jockey,
	payment types.Payment,
	payout types.Payout,
	hitTickets []*Ticket,
	status types.PredictStatus,
) *ListRow {
	return &ListRow{
		race:           race,
		favoriteHorse:  favoriteHorse,
		rivalHorse:     rivalHorse,
		favoriteJockey: favoriteJockey,
		rivalJockey:    rivalJockey,
		payment:        payment,
		payout:         payout,
		hitTickets:     hitTickets,
		status:         status,
	}
}

func (l *ListRow) Race() *Race {
	return l.race
}

func (l *ListRow) FavoriteHorse() *Horse {
	return l.favoriteHorse
}

func (l *ListRow) RivalHorse() *Horse {
	return l.rivalHorse
}

func (l *ListRow) FavoriteJockey() *Jockey {
	return l.favoriteJockey
}

func (l *ListRow) RivalJockey() *Jockey {
	return l.rivalJockey
}

func (l *ListRow) Payment() types.Payment {
	return l.payment
}

func (l *ListRow) Payout() types.Payout {
	return l.payout
}

func (l *ListRow) HitTickets() []*Ticket {
	return l.hitTickets
}

func (l *ListRow) Status() types.PredictStatus {
	return l.status
}
