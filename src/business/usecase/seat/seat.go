package seat

import (
	categoryDom "go-clean/src/business/domain/category"
	eventDom "go-clean/src/business/domain/event"
	seatDom "go-clean/src/business/domain/seat"
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
)

type Interface interface {
	GetList(param entity.SeatParam) (entity.SeatResponse, error)
}

type seat struct {
	seat     seatDom.Interface
	category categoryDom.Interface
	event    eventDom.Interface
}

func Init(sd seatDom.Interface, cd categoryDom.Interface, ed eventDom.Interface) Interface {
	s := &seat{
		seat:     sd,
		category: cd,
		event:    ed,
	}

	return s
}

func (s *seat) GetList(param entity.SeatParam) (entity.SeatResponse, error) {
	res := entity.SeatResponse{}

	event, err := s.event.Get(entity.EventParam{
		ID: param.EventId,
	})
	if err != nil {
		return res, errors.NewError("failed to get event data", err.Error())
	}

	category, err := s.category.Get(entity.CategoryParam{
		ID: param.CategoryId,
	})
	if err != nil {
		return res, errors.NewError("failed to get category data", err.Error())
	}

	seats, err := s.seat.GetList(param)
	if err != nil {
		return res, errors.NewError("failed to get seat data", err.Error())
	}

	if len(seats) == 0 {
		return res, errors.NewError("seat not found", "failed to get seat, param : %#v", param)
	}

	seatMap := make(map[string][]entity.SeatListRows)
	for _, s := range seats {
		slr := entity.SeatListRows{
			SeatID:     s.ID,
			Number:     s.Number,
			IsReserved: s.IsReserved,
		}
		if _, ok := seatMap[s.Row]; ok {
			seatMap[s.Row] = append(seatMap[s.Row], slr)
			continue
		}
		seatMap[s.Row] = append(seatMap[s.Row], slr)
	}

	res.EventName = event.Name
	res.Category = category.Name
	for k, v := range seatMap {
		res.Rows = append(res.Rows, entity.SeatRows{
			Row:   k,
			Seats: v,
		})
	}

	return res, nil
}
