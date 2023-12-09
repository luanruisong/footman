package footman

import "github.com/pkg/errors"

var (
	RootErr        = errors.New("Footman Error")
	ErrNoData      = errors.Wrap(RootErr, "no data")
	ErrOutOfRange  = errors.Wrap(RootErr, "out of range")
	ErrReadTimeout = errors.Wrap(RootErr, "timeout")
	ErrTopicType   = errors.Wrap(RootErr, "topic type err")
)
