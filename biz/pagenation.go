package biz

import (
	"log"
	"time"
	"webagent/model"
	"webagent/pkg/crawler"
)

type pagenation struct {
	mtab           *crawler.ChromeTab
	Pagenationrule string
	SpiderRule     string
	Datalist       *[]model.CrawlerData
}

func NewPagenationRule(tab *crawler.ChromeTab, webinfoid string, pagerule string, spdierrule string) (*pagenation, error) {
	var p pagenation
	var datas = make([]model.CrawlerData, 0)
	p.Datalist = &datas
	p.Pagenationrule = pagerule
	p.mtab = tab
	p.SpiderRule = spdierrule
	return &p, nil
}

func (p pagenation) has(url string) bool {
	for i, _ := range *p.Datalist {
		var v = *p.Datalist
		if v[i].Url == url {
			return true
		}
	}
	return false
}

func (p pagenation) RunDynic(pagerule string, millisecond int) bool {
	err := p.mtab.NoWaitEvaluate(pagerule, nil)
	if err != nil && err.Error() != "encountered an undefined value" {
		return false
	}
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	var datas []model.CrawlerData
	err = p.mtab.NoWaitEvaluate(p.SpiderRule, &datas)
	if err != nil {
		log.Println("执行脚本失败或超过分页范围," + err.Error())
		return false
	}
	if datas == nil || len(datas) == 0 {
		return false
	}
	for i, _ := range datas {
		if p.has(datas[i].Url) {
			log.Println("爬取到重复项，执行失败。请确认是执行间隔过短，或者对方设置问题")
			return false
		}
	}
	*p.Datalist = append(*p.Datalist, datas...)
	return true
}

func (p pagenation) RunStatic(url string) bool {
	var tab = crawler.Instance().NewTab()
	defer tab.Close()
	var datas []model.CrawlerData
	err := tab.NavigateEvaluate(url, p.SpiderRule, &datas)
	if err != nil && err.Error() != "encountered an undefined value" {
		log.Println(p.SpiderRule + "执行脚本失败," + err.Error())
		return false
	}
	if datas == nil || len(datas) == 0 {
		//未取到数据，已经到底或者错误
		return false
	}
	for i, _ := range datas {
		if p.has(datas[i].Url) {
			log.Println("爬取到重复项，执行失败。请确认规则设计有无错误")
			return false
		}
	}
	*p.Datalist = append(*p.Datalist, datas...)
	return true
}
