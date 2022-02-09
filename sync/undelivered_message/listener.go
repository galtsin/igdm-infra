package undelivered_message

import (
	"fmt"
	"sync"
	"time"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/domain/case/get_undelivered_messages"
	"channels-instagram-dm/domain/model"
)

func Listen(runtimeContext domain.RuntimeContext, wg *sync.WaitGroup, ch chan model.MessagesBatch, account model.Account) error {
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()

		tickDuration := 1 * time.Minute
		ticker := time.NewTicker(3 * time.Minute) // First run duration

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-runtimeContext.Context().Done():
				runtimeContext.Logger().Debug("Context was closed", nil)
				return
			case <-ticker.C:
			}

			response, err := get_undelivered_messages.Run(runtimeContext, get_undelivered_messages.Request{
				Account: account,
			})

			if err != nil {
				runtimeContext.Logger().Error(fmt.Sprintf("%s", err), nil)
				continue
			}

			runtimeContext.Logger().Info(fmt.Sprintf("Recieved [%d]", len(response.Messages)), nil)

			// Группируем по беседе
			batch := make(model.MessagesBatch)
			for _, message := range response.Messages {
				ms, ok := batch[message.ConversationID]
				if !ok {
					batch[message.ConversationID] = make([]model.Message, 0)
				}

				batch[message.ConversationID] = append(ms, message)
			}

			ch <- batch
			ticker.Reset(tickDuration)
		}
	}()

	return nil
}
