package validate

import (
	"errors"
	"time"
)

var (
	errTimeExpired = errors.New("expired time")
)

type TimeParser interface {
	GetTime(d []byte) time.Time
}

type Expire struct {
	timeParser TimeParser
	lifeTime   time.Duration
	futureTime time.Duration
}

func NewExpire(timeParser TimeParser, lifeTime time.Duration, futureTime time.Duration) *Expire {
	return &Expire{timeParser: timeParser, lifeTime: lifeTime, futureTime: futureTime}
}

func (e *Expire) Validate(_ string, data []byte) error {
	now := time.Now()
	t := e.timeParser.GetTime(data)
	if now.Add(-e.lifeTime).After(t) || now.Add(e.futureTime).Before(t) {
		return errTimeExpired
	}
	return nil
}
