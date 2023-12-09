package footman

import (
	"context"
	"errors"
	"sync"
	"time"
)

type (
	processor struct {
		offset uint64
		topic  *Topic
	}
	processors struct {
		m   *sync.Map
		svr *Svr
	}
	Consumer struct {
		timerAfter time.Duration
		svr        *Svr
		ps         *processors
	}
)

func (c *Consumer) Timer(d time.Duration) *Consumer {
	c.timerAfter = d
	return c
}

func (c *Consumer) Subscribe(topic ...string) *Consumer {
	c.ps.Subscribe(topic...)
	return c
}

func (c *Consumer) init() *Consumer {
	c.ps = &processors{
		m:   &sync.Map{},
		svr: c.svr,
	}
	c.Timer(time.Second / 10)
	return c
}

func (c *Consumer) ReadMessage(d time.Duration) ([]*message, error) {
	timer := time.NewTimer(0)
	ctx, cancel := context.WithCancel(context.Background())
	if d > 0 {
		ctx, cancel = context.WithTimeout(ctx, d)
	} else if d == 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Second/2)
	}
	defer cancel()
	for {
		select {
		case <-timer.C:
			ret, err := c.ps.Find()
			if err != nil {
				return nil, err
			}
			if len(ret) > 0 {
				return ret, nil
			}
			timer.Reset(c.timerAfter)
		case <-ctx.Done():
			return nil, ErrReadTimeout
		}
	}
}

func (p *processor) find() ([]*message, error) {
	msg, err := p.topic.Find(p.offset)
	if err == nil && len(msg) > 0 {
		p.offset = msg[len(msg)-1].Offset()
	}
	return msg, err
}

func (p *processors) Subscribe(topic ...string) {
	for _, v := range topic {
		if _, ok := p.m.Load(v); !ok {
			top := p.svr.LoadTopic(v)
			p.m.Store(v, &processor{
				offset: top.Offset(),
				topic:  top,
			})
		}
	}
}

func (p *processors) Find() ([]*message, error) {
	var (
		ret   = make([]*message, 0)
		limit int
		ch    = make(chan []*message)
		ech   = make(chan error)
		err   error
	)
	p.m.Range(func(key, value any) bool {
		limit++
		go func(pp *processor) {
			info, ce := pp.find()
			if ce != nil && !errors.Is(ce, ErrNoData) {
				ech <- ce
				return
			}
			ch <- info
		}(value.(*processor))
		return true
	})
	for i := 0; i < limit; i++ {
		select {
		case info := <-ch:
			ret = append(ret, info...)
		case e := <-ech:
			if !errors.Is(e, ErrNoData) {
				err = e
			}
		}
	}
	close(ch)
	close(ech)
	return ret, err
}

func NewConsumer(s *Svr) *Consumer {
	return &Consumer{
		svr: s,
	}
}
