package main

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func Test_C(t *testing.T) {
	errr := tt()
	fmt.Println(errr.Error())
}

func tt() error{
	var err error
	defer func() {
		var er = recover()
		if er != nil{
			err = er.(error)
		}
	}()
	panic(errors.New("测试"))
	return err
}

