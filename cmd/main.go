package main

import (
	"github.com/luopengift/file"
	"github.com/luopengift/log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("close")
		}
	}()
	f := file.NewCatFiles("test.log")
	f.ReadLine()
	for v := range f.NextLine() {
		log.Info(string(v))
	}
	select {}
}
