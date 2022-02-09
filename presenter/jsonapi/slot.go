package jsonapi

import (
	"encoding/json"

	"channels-instagram-dm/domain"
)

type SlotPresenter interface {
	MarshalList([]domain.SlotContainer) ([]byte, error)
}

type slotPresenter struct {
	included []*Resource
}

type Slot struct {
	Type
	Attributes SlotAttributes `json:"attributes"`
}

type SlotAttributes struct {
	Host       string   `json:"host"`
	Status     string   `json:"status"`
	Username   string   `json:"username"`
	Users      []string `json:"users"`
	ActiveUser string   `json:"active_user"`
}

func NewSlotPresenter() SlotPresenter {
	return &slotPresenter{}
}

func (p *slotPresenter) MarshalList(list []domain.SlotContainer) ([]byte, error) {
	slots := make([]Slot, 0, len(list))

	for _, sc := range list {
		slot := Slot{}
		slot.fromModel(sc)
		slots = append(slots, slot)
	}

	result := struct {
		Data []Slot `json:"data"`
	}{
		Data: slots,
	}

	return json.Marshal(result)
}

func (s *Slot) fromModel(sc domain.SlotContainer) {
	s.Type.ID = sc.Slot.Host
	s.Type.Type = "slot"

	s.Attributes.Host = sc.Slot.Host
	s.Attributes.Status = string(sc.Slot.Status)
	s.Attributes.Username = sc.Username
	s.Attributes.Users = sc.Metadata.Users
	s.Attributes.ActiveUser = sc.Metadata.ActiveUser
}
