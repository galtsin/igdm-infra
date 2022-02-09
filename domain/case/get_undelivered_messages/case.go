package get_undelivered_messages

import (
	"fmt"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/model"
)

type Request struct {
	Account model.Account
}

type Response struct {
	Messages []model.Message
}

func Run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	runtimeContext.Logger().Info("[get_undelivered_messages] Case run", nil)

	resp, err := run(runtimeContext, req)
	if err != nil {
		runtimeContext.Logger().Error(fmt.Sprintf("[get_undelivered_messages] Case err [%s]", err), nil)
		return resp, err
	}

	return resp, nil
}

func run(runtimeContext domain.RuntimeContext, req Request) (Response, error) {
	resp := Response{}

	limit := 50

	messages := make([]model.Message, 0, limit)

	// Обходим в порядке приоритете
	filter := runtimeContext.Repository().MessageRepository().Filter()
	filter.WithAccountID(req.Account.ID)

	messagesRecent, err := runtimeContext.Repository().MessageRepository().WhereChannelsDeliveredFailedRecentAt(filter, 2*time.Hour, limit)
	if err != nil {
		return resp, err
	}

	messages = append(messages, messagesRecent...)

	if len(messages) == limit {
		resp.Messages = messages
		return resp, nil
	}

	messagesRecent, err = runtimeContext.Repository().MessageRepository().WhereChannelsDeliveredNone(filter, limit-len(messages))
	if err != nil {
		return resp, err
	}

	messages = append(messages, messagesRecent...)

	if len(messages) == limit {
		resp.Messages = messages
		return resp, nil
	}

	filter = runtimeContext.Repository().MessageRepository().Filter()
	filter.WithAccountID(req.Account.ID)

	messagesRecent, err = runtimeContext.Repository().MessageRepository().WhereInstagramDeliveredFailedRecentAt(filter, 2*time.Hour, limit-len(messages))
	if err != nil {
		return resp, err
	}

	messages = append(messages, messagesRecent...)

	if len(messages) == limit {
		resp.Messages = messages
		return resp, nil
	}

	messagesRecent, err = runtimeContext.Repository().MessageRepository().WhereInstagramDeliveredNone(filter, limit-len(messages))
	if err != nil {
		return resp, err
	}

	messages = append(messages, messagesRecent...)

	resp.Messages = messages

	return resp, nil
}
