package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/luopengift/log"
)

// Tail tail
type Tail struct {
	*File
	cname    string //config name
	buf      chan []byte
	reader   *bufio.Reader
	interval int // 轮转间隔
	endwait  int // 停止等待时间
	Handler      // file name handle interface
}

// NewTail new tail
func NewTail(cname string, handler Handler) *Tail {
	name := handler.Handle(cname)
	file := NewFile(name, os.O_RDONLY)

	tail := new(Tail)
	tail.File = NewFile(name, os.O_RDONLY)
	tail.cname = cname
	tail.buf = make(chan []byte, 1)
	tail.reader = bufio.NewReader(file.fd)
	tail.interval = 10 * 1000                //10s
	tail.endwait = 365 * 1000 * 60 * 60 * 24 //365d
	tail.Handler = handler
	return tail
}

// SetEndWaitTime set endwait time
func (t *Tail) SetEndWaitTime(times int) {
	t.endwait = times
}

// SetInterval set interval time
func (t *Tail) SetInterval(times int) {
	t.interval = times
}

// Close close
func (t *Tail) Close() error {
	close(t.buf)
	return t.File.Close()
}

// ReOpen re open
func (t *Tail) ReOpen() error {
	if err := t.File.Close(); err != nil {
		return log.Errorf("name: %v, %v", t.name, err)
	}
	t.name = t.Handler.Handle(t.cname)
	err := t.Open()
	if err != nil {
		return err
	}
	t.reader = bufio.NewReader(t.fd)
	return nil
}

// Stop stop
func (t *Tail) Stop() {
	t.File.Close()
	close(t.buf)
}

func (t *Tail) rotate() error {
	if t.name == t.cname {
		ok, err := t.IsSameFile(t.name)
		if err != nil {
			return err
		}
		if !ok {
			return t.ReOpen()
		}
	}
	if t.name != t.Handler.Handle(t.cname) { //检测是否需要按时间轮转新文件
		return t.ReOpen()
	}
	return nil
}

// ReadLine read line
func (t *Tail) ReadLine() error {
	offset, err := t.TrancateOffsetByLF(t.seek)
	if err != nil {
		log.Error("<Trancate offset:%d,Error:%+v>", t.seek, err)
		return err
	}
	_, err = t.Seek(offset, 0)
	if err != nil {
		log.Error("<seek offset[%d] error:%+v>", t.seek, err)
		return err
	}
	endintervel := time.Duration(t.endwait) * time.Millisecond
	endTimer := time.NewTimer(endintervel)

	renameintervel := time.Duration(t.interval) * time.Millisecond
	renameTimer := time.NewTimer(renameintervel)
	for {
		line, err := t.reader.ReadBytes('\n')
		switch {
		case err == io.EOF:
			select {
			case <-endTimer.C:
				t.Stop()
				return log.Errorf("over %v second, no more new msg", t.endwait)
			case <-renameTimer.C:
				t.rotate()
			}
		case err != nil && err != io.EOF:
			select {
			case <-renameTimer.C:
				t.rotate()
			}
		default:
			msg := bytes.TrimRight(line, "\n")
			t.buf <- msg
			t.seek += int64(len(line))
			endTimer.Reset(endintervel)
		}
		renameTimer.Reset(renameintervel)
	}

}

// NextLine nextline
func (t *Tail) NextLine() chan []byte {
	return t.buf
}

// Read read
func (t *Tail) Read(p []byte) (int, error) {
	msg, ok := <-t.buf
	if !ok {
		return 0, fmt.Errorf("file is closed")
	}
	if len(msg) > len(p) {
		return 0, errors.New("message is large than buf")
	}
	n := copy(p, msg)
	return n, nil
}
