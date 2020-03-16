package util

import (
	"github.com/pkg/errors"
	"time"
)

//todo 必要时要添加格式
var examples = []string{
	"3:04:05.000 PM Mon Jan",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	"2006/01/02 15:04:05",

}

func ParseAnyTime(input string) (time.Time, error){
	for _, d := range examples {
		t, er := time.Parse(d, input)
		if er == nil {
			return t, nil
		}
	}
	return time.Now(), errors.New("都转不了")
}