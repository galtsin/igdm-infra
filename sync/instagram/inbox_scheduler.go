package instagram

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"channels-instagram-dm/domain"
)

const (
	RoundDurationDefault = 12 * time.Minute // 15
	RoundDurationMax     = 3 * time.Hour

	AttemptRoundDurationDefault = 1 * time.Minute
	AttemptsNumMax              = 3
)

type scheduler struct {
	isAttemptMode               bool
	attemptsNum                 int
	attemptsNumMax              int
	attemptRoundDurationDefault time.Duration
	roundDurationDefault        time.Duration
	roundDurationMax            time.Duration
	roundDuration               time.Duration
	ctx                         context.Context
	ticker                      *time.Ticker
	mux                         *sync.Mutex
	logger                      domain.Logger
}

func (s *scheduler) Ticker() *time.Ticker {
	return s.ticker
}

func (s *scheduler) setRound(d time.Duration) {
	s.logger.Info(fmt.Sprintf("Duration set in [%f] sec, at [%s]", d.Seconds(), time.Now().Add(d).Format(time.RFC3339)), nil)
	s.ticker.Reset(d)
}

// Fail Включается режим попытки на 1 раунд
func (s *scheduler) Fail() *scheduler {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.isAttemptMode = true

	return s
}

func (s *scheduler) Reset() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.isAttemptMode = false
	s.attemptsNum = 0
	s.roundDuration = s.roundDurationDefault
}

func (s *scheduler) Next() {
	s.mux.Lock()
	defer s.mux.Unlock()

	// Если не было ошибок на предыдущем раунде, то сбрасываем счетчик попыток
	if !s.isAttemptMode {
		s.attemptsNum = 0
	}

	// Режим попытки действует только на 1 раунд
	// Если произошла повторная ошибка - вызвать Fail(), иначе будем считать, что все хорошо
	if s.isAttemptMode {
		s.isAttemptMode = false

		if s.attemptsNum < s.attemptsNumMax {
			roundDuration := calcAttemptDuration(s.attemptsNum, s.attemptRoundDurationDefault)
			s.setRound(jitter(roundDuration))

			s.attemptsNum++

			s.logger.Info(fmt.Sprintf("Round attempt [%d]", s.attemptsNum), nil)
			return
		}

		if s.attemptsNum > s.attemptsNumMax {
			s.attemptsNum = 0
			s.logger.Info(fmt.Sprintf("Round attempt [%d]: Limit [%d] reached", s.attemptsNum, s.attemptsNumMax), nil)
		}
	}

	roundDuration := calcRoundDuration(s.roundDuration, s.roundDurationMax)
	s.roundDuration = roundDuration

	s.setRound(jitter(roundDuration))
}

func NewScheduler(ctx context.Context, logger domain.Logger) *scheduler {
	sch := &scheduler{
		attemptsNumMax:              AttemptsNumMax,
		attemptRoundDurationDefault: AttemptRoundDurationDefault,
		roundDurationDefault:        RoundDurationDefault,
		roundDurationMax:            RoundDurationMax,
		roundDuration:               RoundDurationDefault,
		mux:                         &sync.Mutex{},
		logger:                      logger,
		ticker:                      time.NewTicker(RoundDurationDefault),
		ctx:                         ctx,
	}

	return sch
}

func jitter(t time.Duration) time.Duration {
	df := 0.15 // Джиттер в пределах 15%
	rand.Seed(time.Now().UnixNano())

	a := float64(t.Nanoseconds())
	j := a*df + a - float64(rand.Intn(int(a*2*df)))

	return time.Duration(j).Round(time.Second)
}

func calcRoundDuration(timePrev time.Duration, roundDurationMax time.Duration) time.Duration {
	df := 0.3 // Приращение
	newDuration := float64(timePrev.Nanoseconds()) * df
	roundDuration := timePrev + time.Duration(newDuration).Round(time.Second)

	if roundDuration > roundDurationMax {
		return roundDurationMax
	}

	return roundDuration
}

func calcAttemptDuration(attemptsNum int, attemptRoundDuration time.Duration) time.Duration {
	newDuration := 1 << attemptsNum * attemptRoundDuration
	roundDuration := time.Duration(newDuration).Round(time.Second)

	return roundDuration
}
