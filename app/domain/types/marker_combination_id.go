package types

import (
	"fmt"
	"strconv"
	"strings"
)

type MarkerCombinationId int

func NewMarkerCombinationId(rawMarkerCombinationId int) (MarkerCombinationId, error) {
	digits := len(strconv.Itoa(rawMarkerCombinationId))
	if digits <= 1 || rawMarkerCombinationId <= 0 {
		return 0, fmt.Errorf("invalid marker combination id: %d", rawMarkerCombinationId)
	}

	ticketTypeId, _ := strconv.Atoi(string(strconv.Itoa(rawMarkerCombinationId)[0]))
	if ticketTypeId <= 0 || ticketTypeId > 7 {
		return 0, fmt.Errorf("invalid marker combination id: %d", ticketTypeId)
	}

	return MarkerCombinationId(rawMarkerCombinationId), nil
}

func (m MarkerCombinationId) Value() int {
	return int(m)
}

func (m MarkerCombinationId) String() string {
	rawMarkerCombinationId := m.Value()
	var rawMarkerCombinationIds []int
	for rawMarkerCombinationId > 0 {
		rawMarkerCombinationIds = append([]int{rawMarkerCombinationId % 10}, rawMarkerCombinationIds...)
		rawMarkerCombinationId = rawMarkerCombinationId / 10
	}

	var (
		ticketType TicketType
		markers    []string
	)
	for idx, rawMarkerId := range rawMarkerCombinationIds {
		if idx == 0 {
			switch rawMarkerId {
			case 1:
				ticketType = Win
			case 2:
				ticketType = Place
			case 3:
				ticketType = QuinellaPlace
			case 4:
				ticketType = Quinella
			case 5:
				ticketType = Exacta
			case 6:
				ticketType = Trio
			case 7:
				ticketType = Trifecta
			}
			continue
		}

		markerId, err := NewMarker(rawMarkerId)
		if err != nil {
			return ""
		}
		markers = append(markers, markerId.String())
	}

	switch ticketType {
	case Win, Place:
		return markers[0]
	case QuinellaPlace, Quinella, Trio:
		return strings.Join(markers, QuinellaSeparator)
	case Exacta, Trifecta:
		return strings.Join(markers, ExactaSeparator)
	}

	return ""
}

func (m MarkerCombinationId) TicketType() TicketType {
	ticketTypeId, _ := strconv.Atoi(string(strconv.Itoa(m.Value())[0]))
	switch ticketTypeId {
	case 1:
		return Win
	case 2:
		return Place
	case 3:
		return QuinellaPlace
	case 4:
		return Quinella
	case 5:
		return Exacta
	case 6:
		return Trio
	case 7:
		return Trifecta
	}

	return UnknownTicketType
}
