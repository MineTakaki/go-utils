package async

import (
	"context"
	"io"
	"sync"
)

type (
	// Writer 非同期Writerインターフェイス
	Writer interface {
		io.Writer

		// Cancel 処理をキャンセルします
		Cancel()

		// Error 発生した最初のエラーを取得します
		Error() error

		// TotalBytes 書き込み済みの総バイト数を取得します
		TotalBytes() int64
	}

	asyncWriter struct {
		ctx       context.Context
		fnCan     context.CancelFunc
		err       error
		w         io.Writer
		writeChan chan []byte
		wg        sync.WaitGroup
		cnt       int64
	}
)

// NewWriter 非同期書き込みを行うWriterを生成します
func NewWriter(ctx context.Context, w io.Writer) Writer {
	aw := &asyncWriter{
		w:         w,
		writeChan: make(chan []byte),
	}
	aw.ctx, aw.fnCan = context.WithCancel(ctx)
	aw.wg.Add(1)
	go aw.writeRoutine()
	return aw
}

func (aw *asyncWriter) Cancel() {
	if aw.err == nil {
		aw.fnCan()
	}
}

func (aw *asyncWriter) Error() error {
	return aw.err
}

func (aw *asyncWriter) TotalBytes() int64 {
	return aw.cnt
}

func (aw *asyncWriter) Write(data []byte) (int, error) {
	if aw.err != nil {
		return 0, aw.err
	}
	aw.writeChan <- data
	return len(data), nil
}

func (aw *asyncWriter) Close() error {
	close(aw.writeChan)
	aw.wg.Wait()
	return aw.err
}

func (aw *asyncWriter) writeRoutine() {
	defer aw.wg.Done()
	for {
		select {
		case <-aw.ctx.Done():
			if aw.err == nil {
				aw.err = aw.ctx.Err()
			}
			return
		case data, ok := <-aw.writeChan:
			if !ok {
				return
			}
			n, err := aw.w.Write(data)
			if err != nil && aw.err == nil {
				aw.err = err
			} else {
				aw.cnt += int64(n)
			}
		}
	}
}
