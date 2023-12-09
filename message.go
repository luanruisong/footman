package footman

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/luanruisong/footmain/ptr"
	"time"
)

type (
	message struct {
		topicName    *string
		offset       *uint64
		ts, originTs *int64
		data         any
	}
)

func (m *message) MarshalJSON() ([]byte, error) {
	data, err := jsoniter.MarshalToString(m.data)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{"topic_name":"%s","offset":%d,"data":%s}`, *m.topicName, *m.offset, data)), nil
}

func (m *message) Topic() string {
	return *m.topicName
}

func (m *message) Offset() uint64 {
	return *m.offset
}

func (m *message) Data() any {
	return m.data
}

func (m *message) Ts() int64 {
	return *m.ts
}

func (m *message) OriginTs() int64 {
	return *m.originTs
}

func NewMessage(topic *string, data any) *message {
	return &message{
		topicName: topic,
		data:      data,
		ts:        ptr.Int64(time.Now().UnixMilli()),
	}
}
