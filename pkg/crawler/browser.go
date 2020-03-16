package crawler

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
	"log"
	"sync"
	"time"
)

var IsDubug bool = true
var Timeout time.Duration = time.Duration(30) * time.Second
var Capacity int = 5
var _instace *chromeBrowser
var locker sync.RWMutex
var FirstPage = "about:blank"
var UrlTimeout error = errors.New("网站已超时")

type chromeBrowser struct {
	ctx *context.Context
	cancel *context.CancelFunc
	sync.RWMutex
	count int
}

func Instance() *chromeBrowser{
	if _instace == nil {
		locker.Lock()
		defer locker.Unlock()
		if _instace == nil {
			ctx, cancel := newctx()
			_instace = &chromeBrowser{
				ctx: &ctx,
				cancel:&cancel,
				count:0,
			}

		}
	}
	return _instace
}

func newctx() (context.Context, context.CancelFunc){
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", !IsDubug),
		chromedp.Flag("ignore-certificate-errors", true),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	br, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	err := chromedp.Run(br, chromedp.Navigate(FirstPage))
	if err != nil {
		panic("无法启动Chrome,请确认有没有安装" + err.Error())
	}
	return br, cancel
}

func (self *chromeBrowser) NewTab() *ChromeTab{
	for {
		self.Lock()
		if self.count < Capacity{
			self.count++
			var brctx = *self.ctx
			var ers = brctx.Err()
			if ers != nil && ers == context.Canceled {
				log.Println("浏览器被关闭，强制重开一个")
				ctx, cl := newctx()
				self.ctx = &ctx
				self.cancel = &cl
				self.count = 0
			}
			self.Unlock()
			taskCtx, cancel := chromedp.NewContext(brctx)
			var tab = ChromeTab{
				Context:    taskCtx,
				CancelFunc: cancel,
				browser: self,
				ch: make(chan struct{}),
				msgchan: make(chan bool),
			}
			tab.listen()
			return &tab
		}
		self.Unlock()
		time.Sleep(time.Second)
	}
}

func (self *chromeBrowser) Destroy(tab *ChromeTab) {
	self.Lock()
	tab = nil
	self.count--
	self.Unlock()
}