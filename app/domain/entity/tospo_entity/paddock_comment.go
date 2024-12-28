package tospo_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type PaddockComment struct {
	horseNumber types.HorseNumber
	comment     string
	evaluation  int
}

func NewPaddockComment(
	horseNumber types.HorseNumber,
	comment string,
	evaluation int,
) *PaddockComment {
	return &PaddockComment{
		horseNumber: horseNumber,
		comment:     comment,
		evaluation:  evaluation,
	}
}

func (p *PaddockComment) HorseNumber() types.HorseNumber {
	return p.horseNumber
}

func (p *PaddockComment) Comment() string {
	return p.comment
}

func (p *PaddockComment) Evaluation() int {
	return p.evaluation
}
