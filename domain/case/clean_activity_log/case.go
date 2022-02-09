package clean_activity_log

import (
	"fmt"

	"channels-instagram-dm/domain"
)

const ListLimit = 100

type Request struct {
	AccountID string
}

func validate(req Request) error {
	if req.AccountID == "" {
		return fmt.Errorf("AccountID should not be empty")
	}

	return nil
}

func Run(runtimeContext domain.RuntimeContext, req Request) error {
	runtimeContext.Logger().Info("[clean_activity_log] Case run", nil)

	if err := run(runtimeContext, req); err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[clean_activity_log] Case err [%s]", err), nil)
		return err
	}

	return nil
}

func run(runtimeContext domain.RuntimeContext, req Request) error {
	if err := validate(req); err != nil {
		return domain.NewErrorInvalidArgument(err.Error())
	}

	activityLogRepository := runtimeContext.Repository().ActivityLogRepository()

	activityLogList, err := activityLogRepository.WhereAccountID(req.AccountID, 1, ListLimit)
	if err != nil {
		return err
	}

	if len(activityLogList) == 0 {
		return nil
	}

	activityLog := activityLogList[0]

	return activityLogRepository.DeleteWhereAccountIDCreatedAtBefore(activityLog.AccountID, activityLog.CreatedAt)
}
