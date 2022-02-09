package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"channels-instagram-dm/domain"
	"channels-instagram-dm/presenter/jsonapi"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	RequestIDHeaderName = "X-Request-Id" // Заголовок для отслеживания Trace ID
	CacheTime           = "X-Cache-Time" // Признак того, что в ответе может использоваться кэш, не старше чем CacheTime
)

type Session struct {
	ctx       context.Context
	requestID string
}

type runtimeContext struct {
	ctx        context.Context
	syncer     domain.Syncer
	logger     domain.Logger
	repository domain.Repository
	service    domain.Service
	eventBus   domain.EventBus
	mq         domain.MQ
}

func RuntimeContext(withContext context.Context, rep domain.Repository, service domain.Service, logger domain.Logger, mq domain.MQ, eventBus domain.EventBus, syncer domain.Syncer) domain.RuntimeContext {
	return &runtimeContext{
		ctx:        withContext,
		syncer:     syncer,
		logger:     logger,
		repository: rep,
		service:    service,
		eventBus:   eventBus,
		mq:         mq,
	}
}

func NewSession(req *http.Request) *Session {
	return &Session{
		ctx:       req.Context(),
		requestID: req.Header.Get(RequestIDHeaderName),
	}
}

func (c *runtimeContext) WithContext(ctx context.Context) domain.RuntimeContext {
	return RuntimeContext(ctx, c.Repository(), c.Service(), c.Logger(), c.MQ(), c.EventBus(), c.Syncer())
}

func (c *runtimeContext) WithLogger(logger domain.Logger) domain.RuntimeContext {
	return RuntimeContext(c.ctx, c.Repository(), c.Service(), logger, c.MQ(), c.EventBus(), c.Syncer())
}

func (c *runtimeContext) Context() context.Context {
	return c.ctx
}

func (c *runtimeContext) Logger() domain.Logger {
	return c.logger
}

func (c *runtimeContext) Repository() domain.Repository {
	return c.repository
}

func (c *runtimeContext) Service() domain.Service {
	return c.service
}

func (c *runtimeContext) MQ() domain.MQ {
	return c.mq
}

func (c *runtimeContext) EventBus() domain.EventBus {
	return c.eventBus
}

func (c *runtimeContext) Syncer() domain.Syncer {
	return c.syncer
}

type routeHandler func(domain.RuntimeContext, *http.Request) ([]byte, error)

func RouteHandler(rc domain.RuntimeContext, r *mux.Router, path string, handler routeHandler) *mux.Route {
	return r.HandleFunc(path, func(resp http.ResponseWriter, req *http.Request) {
		runtimeContext := rc.WithLogger(
			rc.Logger().Copy(uuid.New().String()), // TraceID
		)

		runtimeContext.Logger().Debug(fmt.Sprintf("Request %s:%s", req.Method, req.URL), nil)

		result, err := handler(runtimeContext, req)

		if err != nil {
			runtimeContext.Logger().Error(fmt.Sprintf("Respond %s:%s with error %s", req.Method, req.URL, err), nil)
			RespondWithError(resp, err)
			return
		}

		runtimeContext.Logger().Debug(fmt.Sprintf("Respond %s:%s", req.Method, req.URL), string(result))

		if req.Method == http.MethodPost {
			RespondWithJSON(resp, http.StatusCreated, result)
			return
		}

		RespondWithJSON(resp, http.StatusOK, result)
	})
}

func RespondWithError(w http.ResponseWriter, err error) {
	httpCode := extractHttpErrorCode(err)
	presenter := jsonapi.NewErrorPresenter()

	var bs []byte
	if errBase, ok := err.(domain.BaseError); ok {
		bs = presenter.Marshal(errBase)
	} else {
		bs = presenter.Marshal(domain.NewError("", err).(domain.BaseError))
	}

	RespondWithJSON(w, httpCode, bs)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

func RespondWithEmptyError(w http.ResponseWriter, err error) {
	httpCode := extractHttpErrorCode(err)
	w.WriteHeader(httpCode)
	w.Write(nil)
}

func extractHttpErrorCode(err error) int {
	if errors.Is(err, domain.ErrorInvalidArgument) || errors.Is(err, domain.ErrorNotFound) {
		return http.StatusBadRequest
	}

	if errors.Is(err, domain.ErrorPermissionDenied) {
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}
