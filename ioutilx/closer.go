package ioutilx

import (
	"io"
	"sync"

	"github.com/MineTakaki/go-utils"
	"github.com/pkg/errors"
)

type (
	//CloseHolder 複数のCloserをまとめてCloseします
	CloseHolder interface {
		io.Closer

		Append(...io.Closer)
	}

	closeHolder struct {
		mutex sync.Mutex
		arr   []io.Closer
	}

	closeWrap struct {
		c io.Closer
	}
)

//Close 指定して全ての Closer を呼び出します
// エラーが起きても最後まで実行し、最初のエラーを返します
func Close(args ...io.Closer) (err error) {
	for _, c := range args {
		if utils.IsNil(c) {
			continue
		}
		if e := c.Close(); err == nil {
			err = e
		}
	}
	return
}

//NewCloserHolder CloserHolderを生成します
func NewCloserHolder(args ...io.Closer) CloseHolder {
	x := &closeHolder{}
	for _, c := range args {
		if !utils.IsNil(c) {
			x.arr = append(x.arr, c)
		}
	}
	return x
}

func (ch *closeHolder) Close() (err error) {
	ch.mutex.Lock()
	for _, c := range ch.arr {
		if e := c.Close(); err == nil {
			err = e
		}
	}
	ch.arr = nil
	ch.mutex.Unlock()
	return
}

func (ch *closeHolder) Append(args ...io.Closer) {
	if len(args) == 0 {
		return
	}
	ch.mutex.Lock()
	for _, c := range args {
		if !utils.IsNil(c) {
			ch.arr = append(ch.arr, c)
		}
	}
	ch.mutex.Unlock()
}

//CloseWithStack io.Closerのerrorにスタックトレースを追加します
func CloseWithStack(c io.Closer) io.Closer {
	if utils.IsNil(c) {
		return nil
	}
	return &closeWrap{c: c}
}

func (cw *closeWrap) Close() error {
	return errors.WithStack(cw.c.Close())
}
