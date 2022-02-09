package instagram_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"channels-instagram-dm/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Request struct {
	ID      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Response struct {
	ID      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Error   string `json:"error,omitempty"`
}

type PayloadResponse struct {
	Response
	Result interface{} `json:"result,omitempty"`
}

type ResultOK string

func (r ResultOK) Ok() bool {
	if r == "ok" {
		return true
	}

	return false
}

type DiscoveryRow struct {
	Host       string   `json:"host"`
	Active     bool     `json:"active"`
	Users      []string `json:"users"`
	ActiveUser string   `json:"active_user"`
}

type service struct {
	mux          sync.Mutex
	cancel       context.CancelFunc
	ctx          context.Context
	logger       domain.Logger
	conn         *websocket.Conn
	inbox        chan struct{}
	listenersMap map[string]listener
}

type listener struct {
	ctx    context.Context
	ch     chan PayloadResponse
	result interface{}
}

func NewRequest(method string) Request {
	return Request{
		ID:      uuid.New().String(),
		JsonRpc: "2.0",
		Method:  method,
	}
}

func NewService(ctxService context.Context, logger domain.Logger, host string) (*service, error) {
	u := url.URL{Scheme: "ws", Host: host}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Open WebSocket: Error. No connect to %s. Dial err %w ", u.String(), err)
	}

	ctx, cancel := context.WithCancel(ctxService)

	s := &service{
		mux:          sync.Mutex{},
		ctx:          ctx,
		cancel:       cancel,
		logger:       logger,
		conn:         c,
		inbox:        make(chan struct{}),
		listenersMap: make(map[string]listener),
	}

	c.SetCloseHandler(func(code int, text string) error {
		logger.Error(fmt.Sprintf("WebSocketConnection: Closed. Code [%d]. Text [%s]", code, text), nil)

		s.Close()
		return nil
	})

	go func() {

		for {
			select {
			case <-ctx.Done():
				logger.Error(ctx.Err().Error(), nil)
				return
			default:

			}

			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Error(fmt.Sprintf("Read message: Error. %s", err), nil)

				// TODO: Это остановка сервиса и нужно как то обработать
				if err.Error() == "websocket: close 1006 (abnormal closure): unexpected EOF" {
					logger.Info(fmt.Sprintf("WebSocketConnection: Need to close"), nil)

					s.Close()
					return
				}

				time.Sleep(10 * time.Second)
				continue
			}

			response := Response{}
			err = json.Unmarshal(message, &response)
			if err != nil {
				logger.Error(fmt.Sprintf("Unmarshal message: Error. %s", err), string(message))
				continue
			}

			if response.Error != "" {
				logger.Error(fmt.Sprintf("Response on [%s]: Error. %s", response.ID, string(message)), nil)
			}

			go func(resp Response, message []byte) {
				s.mux.Lock()
				defer s.mux.Unlock()

				listener, ok := s.listenersMap[resp.ID]
				if !ok {
					logger.Info(fmt.Sprintf("No listener on [%s]", resp.ID), nil)
					return
				}

				response := PayloadResponse{
					Response: resp,
				}

				if response.Error == "" && listener.result != nil {
					response.Result = listener.result
					err = json.Unmarshal(message, &response)
					if err != nil {
						logger.Error(fmt.Sprintf("Unmarshal message: Error. %s", err), string(message))
						return
					}
				}

				logger.Debug(fmt.Sprintf("Recived message on listener [%s]", resp.ID), string(message))

				select {
				case <-ctx.Done():
					return
				case <-listener.ctx.Done():
					return
				case listener.ch <- response:
					return
				}
			}(response, message)
		}
	}()

	//  Зачищаем неактуальные listeners
	go func() {
		for {
			s.mux.Lock()

			for k, listener := range s.listenersMap {
				select {
				case <-listener.ctx.Done():
					// logger.Debug(fmt.Sprintf("Listener [%s] was deleted", k), nil)

					close(listener.ch)
					delete(s.listenersMap, k)
				default:
					continue
				}
			}

			s.mux.Unlock()

			select {
			case <-ctx.Done():
				// Оповещаем слушателей, что закрыли сервис
				s.mux.Lock()

				for k, listener := range s.listenersMap {
					close(listener.ch)
					delete(s.listenersMap, k)
				}

				s.mux.Unlock()
				return
			case <-time.After(10 * time.Second):
				continue
			}
		}
	}()

	return s, nil
}

func (s *service) send(ctx context.Context, listenerID string, req Request, result interface{}) (chan PayloadResponse, error) {
	if s.IsClosed() {
		return nil, fmt.Errorf("Service has been closed")
	}

	ch := make(chan PayloadResponse)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	s.listenersMap[listenerID] = listener{
		ctx:    ctx,
		ch:     ch,
		result: result,
	}

	s.logger.Debug(fmt.Sprintf("Prepare send: Listener [%s] on method [%s] was added", listenerID, req.Method), nil)

	if err := s.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return nil, err
	}

	s.logger.Debug(fmt.Sprintf("Sent message for listener [%s]: Success", listenerID), string(payload))

	return ch, nil
}

func (s *service) Close() {
	if err := s.conn.Close(); err != nil {
		s.logger.Error(fmt.Sprintf("Close service: Error. %s", err), nil)
	}

	s.cancel()

	// TODO: Возможно стоит добавить событие, чтобы оповещать, что сервис закрылся

	s.logger.Info("Close service: Success", nil)
}

func (s *service) IsClosed() bool {
	return s.ctx.Err() != nil
}
