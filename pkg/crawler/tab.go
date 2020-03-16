package crawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"sync"
)

type ChromeTab struct {
	context.Context
	context.CancelFunc
	loaded  bool
	stopped bool
	ch      chan struct{}
	sync.Once
	browser    *ChromeBrowser
	SpiderRule string
}

func (self ChromeTab) Close() {
	if self.browser != nil {
		self.browser.Destroy(&self)
	} else {
		fmt.Println("已退出")
	}
	self.CancelFunc()
}

func (self ChromeTab) Reset() {
	self.stopped = false
	self.loaded = false
}

func (self ChromeTab) Stop() {
	//self.CancelFunc()
}

type _baseGetAction func(context.Context) error

//func (self ChromeTab) baseGet(url string, f _baseGetAction) error{
//	ch := make(chan struct{})
//	defer self.Close()
//	var err error
//
//	err = chromedp.Run(self,
//		chromedp.Navigate(url))
//	if err != nil {
//		fmt.Println("url是否有误？url:" + url + ",error:" + err.Error())
//		return err
//	}
//	select {
//	case <-time.After(time.Millisecond * time.Duration(Timeout)):
//		fmt.Println(url + ",加载超时")
//		return errors.New("加载超时")
//	case <-ch:
//		return f(self)
//	}
//}

//获取pdf字节流
func (self ChromeTab) GetPdfBytes(url string) ([]byte, error) {
	var er error
	var pdfBuffer []byte
	er = chromedp.Run(self,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	return pdfBuffer, er
}

//获取html文本
func (self ChromeTab) Gethtml(url string) (string, error) {
	var err error
	var res string
	err = self.Navigate(url)
	if err != nil {
		return res, err
	}
	var ids []cdp.NodeID
	err = chromedp.Run(self,
		chromedp.NodeIDs(`document`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if len(ids) > 0 {
				res, err = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
				if err != nil {
					return err
				}
				return nil
			} else {
				return errors.New(`获取HTML文档失败,长度为0`)
			}
		}),
	)
	return res, err
}

func (self ChromeTab) NavigateEvaluate(url string, rule string, v interface{}) error {
	var err = chromedp.Run(self, chromedp.Navigate(url))
	if err != nil {
		return err
	}
	err = chromedp.Run(self, chromedp.Evaluate(rule, &v))
	return err
}

func (self ChromeTab) Navigate(url string) error {
	return chromedp.Run(self, chromedp.Navigate(url))
}

func (self ChromeTab) NoWaitEvaluate(rule string, v interface{}) error {
	return chromedp.Run(self, chromedp.Evaluate(rule, &v))
}

func (self ChromeTab) GetData() (interface{}, error) {
	var v interface{}
	err := chromedp.Run(self, chromedp.Evaluate(self.SpiderRule, &v))
	return v, err
}
