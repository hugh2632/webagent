package crawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

type ChromeTab struct {
	context.Context
	context.CancelFunc
	loaded  bool
	stopped bool
	ch      chan struct{}
	msgchan chan bool
	browser *chromeBrowser
	sync.Once
	SpiderRule string
}

func (self ChromeTab) Close() {
	if self.browser != nil {
		self.browser.Destroy(&self)
	}
	if self.Err() != context.DeadlineExceeded{
		self.CancelFunc()
	}
}

func (self ChromeTab) Reset() {
	self.stopped = false
	self.loaded = false
}

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

func (self ChromeTab) reset(){
	self.loaded = false
	self.stopped = false
}

func (self ChromeTab) listen(){
	go func() {
		self.Do(func() {
			chromedp.ListenTarget(self, func(ev interface{}) {
				if IsDubug{
					te :=  reflect.Indirect(reflect.ValueOf(ev)).Type()
					name := te.String()
					if strings.HasPrefix(name, "page."){
						fmt.Println(name + "\t" + reflect.ValueOf(ev).Elem().String())
					}
				}
				switch ev.(type) {
				case *page.EventLoadEventFired://两个事件确保加载完了页面
					go func() {
						self.loaded= true
						if self.stopped{
							self.ch <- struct{}{}
						}
					}()
				case *page.EventFrameStoppedLoading://会多次触发。。 不知道原因
					go func() {
						self.stopped = true
						if self.loaded {
							self.ch <- struct{}{}
						}
					}()
				}
			})
		})
	}()
	go func() {
		select{
		case <-time.After(time.Duration(Timeout)):
			log.Println("加载超时")
			self.msgchan <- false
		case <-self.ch:
			self.reset()
			self.msgchan <- true
		}
	}()
}


//获取html文本
func (self ChromeTab) Gethtml(url string) (string, error) {
	var err error
	var res string
	err = self.Navigate(url)
	if err != nil {
		return res, err
	}
	b := <- self.msgchan
	if !b {
		return res, UrlTimeout
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
	b := <- self.msgchan
	if !b {
		return UrlTimeout
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
