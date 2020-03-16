package main

import (
	"testing"
	"webagent/pkg/crawler"
)

func Test_C(t *testing.T) {
	tab := crawler.Instance().NewTab()
	_ = tab.Navigate("https://www.ics-cert.org.cn/portal/page/111/index_3.html#")
	select {}
}
