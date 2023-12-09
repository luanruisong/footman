package footman

import (
	"github.com/luanruisong/footmain/ptr"
	"sync"
)

type (
	Svr struct {
		topics *sync.Map
		opts   []Option
	}
)

func NewSvr(o ...Option) *Svr {
	return &Svr{
		topics: &sync.Map{},
		opts:   o,
	}
}

func (s *Svr) t(topicName string) *Topic {
	t := &Topic{
		name: ptr.String(topicName),
	}
	for i := range s.opts {
		s.opts[i](t)
	}
	return t
}

func (s *Svr) LoadTopic(topicName string) *Topic {
	t, loaded := s.topics.LoadOrStore(topicName, s.t(topicName))
	topic := t.(*Topic)
	if !loaded {
		topic.init()
	}
	return topic
}

func (s *Svr) RemoveTopic(topicName string) {
	_, loaded := s.topics.Load(topicName)
	if loaded {
		s.topics.Delete(topicName)
	}
}

func (s *Svr) Subscribe(topic ...string) *Consumer {
	return NewConsumer(s).init().Subscribe(topic...)
}

func (s *Svr) Produce(topic string, data any) {
	s.LoadTopic(topic).Append(data)
}

func (s *Svr) ProduceMessage(msg *message) {
	s.LoadTopic(msg.Topic()).AppendMessage(msg)
}
