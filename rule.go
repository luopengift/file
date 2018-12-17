package file

import (
	"strings"
	"time"
)

var (
	// TimeRule time
	TimeRule *TimeRuler
	// NullRule null
	NullRule *NullRuler
	//Map 时间通配符，用于正则表达式替换
	Map = map[string]string{
		"%Y": "2006",
		"%M": "01",
		"%D": "02",
		"%h": "15",
		"%m": "04",
		"%s": "05",
	}
)

// Handler 处理接口
type Handler interface {
	Handle(string) string
}

// TimeRuler time
type TimeRuler struct{}

// Handle handler
// eg:"test-%Y%M%D.log" ->"test-20170203.log"
// eg:"test-%Y-%M-%D.log" ->"test-2017-02-03.log"
func (t *TimeRuler) Handle(str string) string {
	for k, v := range Map {
		str = strings.Replace(str, k, time.Now().Format(v), -1)
	}
	return str
}

// NullRuler null
type NullRuler struct{}

// Handle handler
func (n *NullRuler) Handle(str string) string {
	return str
}

func init() {
	TimeRule = new(TimeRuler)
	NullRule = new(NullRuler)
}
