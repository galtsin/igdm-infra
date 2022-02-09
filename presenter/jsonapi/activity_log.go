package jsonapi

import (
	"encoding/json"
	"time"

	"channels-instagram-dm/domain/model"
)

type ActivityLogPresenter interface {
	MarshalList([]model.ActivityLog) ([]byte, error)
}

type activityLogPresenter struct{}

type ActivityLog struct {
	Type
	Attributes ActivityLogAttributes `json:"attributes"`
}

type ActivityLogAttributes struct {
	Log       string    `json:"log"`
	CreatedAt time.Time `json:"created_at"`
}

func NewActivityLogPresenter() ActivityLogPresenter {
	return &activityLogPresenter{}
}

func (p *activityLogPresenter) MarshalList(list []model.ActivityLog) ([]byte, error) {
	presenterModelList := make([]ActivityLog, 0, len(list))

	for _, item := range list {
		presenterModel := ActivityLog{}
		presenterModel.fromModel(item)
		presenterModelList = append(presenterModelList, presenterModel)
	}

	result := struct {
		Data []ActivityLog `json:"data"`
	}{
		Data: presenterModelList,
	}

	return json.Marshal(result)
}

func (a *ActivityLog) fromModel(m model.ActivityLog) {
	a.Type.Type = "activity_log"

	a.Attributes.Log = m.Log
	a.Attributes.CreatedAt = m.CreatedAt
}
