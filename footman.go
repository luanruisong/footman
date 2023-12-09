package footman

import (
	"github.com/pkg/errors"
)

func LimitOpt(limit int) Option {
	return func(topic *Topic) {
		topic.Limit(limit)
	}
}

func Timeout(err error) bool {
	return errors.Is(err, ErrReadTimeout)
}

type (
	Option func(topic *Topic)
)
