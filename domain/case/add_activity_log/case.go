package add_activity_log

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

type Request struct {
	AccountID string
	Log       string
}

type Response struct {
	ActivityLog model.ActivityLog
}

func validate(req Request) error {
	if req.AccountID == "" {
		return fmt.Errorf("AccountID should not be empty")
	}

	if req.Log == "" {
		return fmt.Errorf("Log should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[add_activity_log] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[add_activity_log] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}

	if err := validate(req); err != nil {
		return resp, domain.NewErrorInvalidArgument(err.Error())
	}

	// Activity log
	activityLog := model.NewActivityLog(req.AccountID)
	activityLog.SetLog(req.Log)

	activityLog, err := runtimeContext.Repository().ActivityLogRepository().Store(activityLog)
	if err != nil {
		return resp, err
	}

	resp.ActivityLog = activityLog

	return resp, nil
}
