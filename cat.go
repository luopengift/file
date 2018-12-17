package file

import (
	"errors"
	"fmt"

	"github.com/luopengift/log"
)

type Cat struct {
	Paths []string
	Files []string
	Regex string
	tails map[string]*Tail
	buf   chan []byte
}

func NewCatPaths(paths ...string) *Cat {
	cat := &Cat{
		Paths: paths,
		tails: make(map[string]*Tail),
	}
	return cat
}

func NewCatFiles(files ...string) *Cat {
	cat := &Cat{
		Files: files,
		tails: make(map[string]*Tail),
	}
	return cat
}
func (c *Cat) ReadLine() error {
	// regexp, err := regexp.Compile(c.Regex)
	// if err != nil {
	// 	return err
	// }
	c.buf = make(chan []byte, 1000)
	// fun := func(path string, info os.FileInfo, err error) error {
	// 	if info.IsDir() {
	// 		return nil
	// 	}
	// 	if regexp.MatchString(info.Name()) {
	// 		c.Tails[path] = NewTail(path, TimeRule)
	// 	}
	// 	return nil
	// }

	for _, f := range c.Files {
		c.tails[f] = NewTail(f, TimeRule)
	}
	for _, tail := range c.tails {
		go func(tail *Tail) {
			for msg := range tail.NextLine() {
				c.buf <- msg
			}
		}(tail)
		go func(tail *Tail) {
			if err := tail.ReadLine(); err != nil {
				log.Error("%v", err)
			}
		}(tail)
	}
	return nil
}

func (c *Cat) NextLine() <-chan []byte {
	return c.buf
}

func (c *Cat) Read(b []byte) (int, error) {
	msg, ok := <-c.buf
	if !ok {
		return 0, fmt.Errorf("file is closed")
	}
	if len(msg) > len(b) {
		return 0, errors.New("message is large than buf")
	}
	n := copy(b, msg)
	return n, nil
}