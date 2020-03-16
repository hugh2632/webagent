package crawler

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

var IsDubug bool = true
var Timeout time.Duration = time.Duration(30) * time.Second
var Capacity int = 5
var _instace *ChromeBrowser
var locker sync.RWMutex
var FirstPage = "https://www.baidu.com"

type ChromeBrowser struct {
	context.Context
	context.CancelFunc
	sync.RWMutex
	count int
}

func Instance() *ChromeBrowser{
	if _instace == nil {
		locker.Lock()
		defer locker.Unlock()
		if _instace == nil {
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
			_instace = &ChromeBrowser{
				Context:    br,
				CancelFunc: cancel,
				count: 0,
			}
		}
	}
	return _instace
}

func (self *ChromeBrowser) NewTab() *ChromeTab{
	for {
		self.Lock()
		if self.count < Capacity{
			self.count++
			self.Unlock()
			taskCtx, cancel := chromedp.NewContext(self)
			var tab = ChromeTab{
				Context:    taskCtx,
				CancelFunc: cancel,
				browser:self,
			}
			go func() {
				tab.Do(func() {
					chromedp.ListenTarget(tab, func(ev interface{}) {
						if IsDubug{
							te :=  reflect.Indirect(reflect.ValueOf(ev)).Type()
							name := te.String()
							//if strings.HasPrefix(name, "page.") || strings.HasPrefix(name, "network.") { //页面事件
							if strings.HasPrefix(name, "page."){
								fmt.Println(name + "\t" + reflect.ValueOf(ev).Elem().String())
							}
						}
						switch ev.(type) {
						case *page.EventLoadEventFired://两个事件确保加载完了页面
							go func() {
								tab.loaded= true
								if tab.stopped{
									tab.ch <- struct{}{}
								}
							}()
						case *page.EventFrameStoppedLoading://会多次触发。。 不知道原因
							go func() {
								tab.stopped = true
								if tab.loaded {
									tab.ch <- struct{}{}
								}
							}()
						}
					})
				})
			}()
			go func() {
				select {
				case <- tab.ch:
					self.Destroy(&tab)
				}
			}()
			return &tab
		}
		self.Unlock()
		time.Sleep(time.Second)
	}
}

func (self *ChromeBrowser) Destroy(tab *ChromeTab) {
	self.Lock()
	tab = nil
	self.count--
	self.Unlock()
}