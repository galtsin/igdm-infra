package get_activity_log

import (
	"fmt"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

type Request struct {
	ExternalID string
}

type Response struct {
	ActivityLogList []model.ActivityLog
}

func validate(req Request) error {
	if req.ExternalID == "" {
		return fmt.Errorf("ExternalID should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[get_activity_log] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[get_activity_log] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}

	if err := validate(req); err != nil {
		return resp, domain.NewErrorInvalidArgument(err.Error())
	}

	accountRepository := runtimeContext.Repository().AccountRepository()

	account, err := accountRepository.WhereExternalID(req.ExternalID)
	if err != nil {
		return resp, err
	}

	activityLogRepository := runtimeContext.Repository().ActivityLogRepository()

	activityLogList, err := activityLogRepository.WhereAccountID(account.ID, 10, 0)
	if err != nil {
		return resp, err
	}

	resp.ActivityLogList = activityLogList

	return resp, nil

}
