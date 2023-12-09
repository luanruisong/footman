package footman

import (
	"github.com/luanruisong/footmain/ptr"
	"sync/atomic"
	"unsafe"
)

type (
	Topic struct {
		name  *string
		limit *int
		seq   *uint64
		data  []*message
	}
)

func (t *Topic) init() {
	t.seq = ptr.Uint64(0)
	if t.limit == nil || *t.limit == 0 {
		t.limit = ptr.Int(10000)
	}
	t.data = make([]*message, *t.limit)
}

func (t *Topic) Limit(i int) {
	t.limit = ptr.Int(i)
}

func (t *Topic) Offset() uint64 {
	return atomic.LoadUint64(t.seq)
}

func (t *Topic) idx(offset uint64) int {
	return int(offset) % *t.limit
}

func (t *Topic) Find(offset uint64) ([]*message, error) {
	seq := atomic.LoadUint64(t.seq)
	switch {
	case offset > seq:
		return nil, ErrOutOfRange
	case offset == seq:
		return nil, ErrNoData
	default:
		ret := make([]*message, 0)
		// 套圈问题
		if ul := uint64(*t.limit); seq > ul {
			if off := seq - ul; off > offset {
				offset = off
			}
		}
		for offset = offset + 1; offset <= seq; offset++ {
			idx := t.idx(offset)
			node := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&t.data[idx])))
			if node != nil {
				ret = append(ret, (*message)(node))
			}
		}
		return ret, nil
	}
}

func (t *Topic) Append(tg any) {
	msg := NewMessage(t.name, tg)
	t.AppendMessage(msg)
}

func (t *Topic) AppendMessage(msg *message) {
	seq := atomic.AddUint64(t.seq, 1)
	idx := t.idx(seq)
	msg.offset = ptr.Uint64(seq)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&t.data[idx])), unsafe.Pointer(msg))
}
