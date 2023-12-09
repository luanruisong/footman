package footman

import (
	"fmt"
	"testing"
	"time"
)

func TestFootman(t *testing.T) {

	svr := NewSvr(LimitOpt(100))
	topics := []string{
		"test1",
		"test2",
	}

	testSub(svr, "t1", topics...)
	testSub(svr, "t2", topics...)

	for i := 0; ; i++ {
		for _, v := range topics {
			svr.Produce(v, fmt.Sprintf(`{"name":"%s-%d"}`, v, i))
		}
		time.Sleep(time.Second)
	}
}

func testSub(s *Svr, key string, topics ...string) {
	go func() {
		c := s.Subscribe(topics...)
		for {
			msg, err := c.ReadMessage(0)
			if err != nil && !Timeout(err) {
				fmt.Println(err.Error())
			} else {
				//msg, _ := jsoniter.MarshalToString(msg)
				for _, v := range msg {
					fmt.Println(key, v.Topic(), v.Data())
				}
			}
		}
	}()
}
