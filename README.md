# Footman

取自 《魔兽争霸:冰封王座》 人族步兵

## usage

```go
//svr
svr := NewSvr(LimitOpt(100))
//consumer
c := svr.Subscribe(topics...)
for {
    msgs, err := c.ReadMessage(-1)// 0,-1,>0 timeout
    if err != nil && !footman.Timeout(err) {
        fmt.Println(err.Error())
    }   
    for _,v := range msgs {
        fmt.Println(v.Data())
    }
}
//producer
svr.Produce(topic, "123")
```

