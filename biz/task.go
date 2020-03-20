package biz

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"webagent/model"
	"webagent/pkg/crawler"
	"webagent/pkg/setting"
	"webagent/util"
	"webagent/util/bloomfilter"
	"webagent/util/htmlUtil"
	mysqlutil "webagent/util/mysql"
)

func TaskListTaskInfo(sort string, order string, pageindex int, rowcount int ) (res []model.AjaxTaskinfo, total int, err error) {
	if rowcount <10 {
		rowcount = 10
	}
	if sort == "" {
		sort = "Taskinfo_createtime"
	}
	if order == "" {
		order = "desc"
	}
	var minLimit = (rowcount - 1)
	var sqlStr = "SELECT Taskinfo_id, Taskinfo_webid, Taskinfo_createtime,Taskinfo_onlyfirst, Taskinfo_rebuild, Taskinfo_starttime, Taskinfo_endtime, Taskinfo_status,Webinfo_name, Webinfo_url FROM taskinfo A inner join webinfo B on A.Taskinfo_webid = B.Webinfo_id "
	var sqlCount = "SELECT count(1) FROM taskinfo A inner join webinfo B on A.Taskinfo_webid = B.Webinfo_id "
	if pageindex >= 0 {
		sqlStr += fmt.Sprintf(`order by Taskinfo_createtime desc limit %v, %v`, minLimit, rowcount)
	}
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		rows, err := db.Query(sqlStr)
		if err != nil {
			return
		}
		for rows.Next() {
			var tmp model.AjaxTaskinfo
			err = rows.Scan(&tmp.Taskinfo_id, &tmp.Taskinfo_webid, &tmp.Taskinfo_createtime,&tmp.Taskinfo_onlyfirst, &tmp.Taskinfo_rebuild , &tmp.Taskinfo_starttime, &tmp.Taskinfo_endtime,&tmp.Taskinfo_status, &tmp.Webinfo_name, &tmp.Webinfo_url )
			if err != nil {
				log.Println("跳过一条反序列化，" + err.Error())
				break
			}
			res = append(res, tmp)
		}
		rows, err = db.Query(sqlCount)
		if err != nil {
			return
		}
		if rows.Next() {
			err = rows.Scan(&total)
			if err != nil {
				return
			}
		}
	})
	return res, total, err
}

func TaskListSite() (res []model.WebInfo, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		rows, err := db.Query(`SELECT
			webinfo.Webinfo_id,
			webinfo.Webinfo_name,
			webinfo.Webinfo_url,
			webinfo.Webinfo_pagenationrule,
			webinfo.Webinfo_spiderrule
			FROM
			webinfo`)
		if err != nil {
			panic(err.Error())
		}
		for rows.Next() {
			var tmp model.WebInfo
			err = rows.Scan(&tmp.Webinfo_id, &tmp.Webinfo_name, &tmp.Webinfo_url, &tmp.Webinfo_pagenationrule, &tmp.Webinfo_spiderrule)
			if err != nil {
				log.Println("跳过一条反序列化，" + err.Error())
				break
			}
			res = append(res, tmp)
		}
	})
	return res, err
}

func TaskNewWebinfo(name string, url string, pagenationrule string, spiderrule string) (id string, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		idRow, err := db.Query("select UUID_Short()")
		if err != nil {
			panic("获取UUID失败," + err.Error())
		}
		idRow.Next()
		_ = idRow.Scan(&id)
		_, err = db.Exec(`INSERT INTO webinfo(Webinfo_id, Webinfo_name,Webinfo_url,Webinfo_pagenationrule,Webinfo_spiderrule) VALUES ('%v', '%v', '%v', '%v','%v')`, id, name, url, pagenationrule, spiderrule)
		if err != nil {
			panic("插入web信息失败," + err.Error())
		}
	})
	return id, err
}

func TaskGetWebinfo(id string) (res model.WebInfo, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		row, err := db.Query("select webinfo.Webinfo_id,webinfo.Webinfo_name,webinfo.Webinfo_url,webinfo.Webinfo_pagenationrule,webinfo.Webinfo_spiderrule,webinfo.Webinfo_snapshot from webinfo where Webinfo_id = '" + id + "'")
		if err != nil {
			panic("获取网站信息失败," + err.Error())
		}
		if row.Next(){
			err = row.Scan(&res.Webinfo_id, &res.Webinfo_name, &res.Webinfo_url, &res.Webinfo_pagenationrule, &res.Webinfo_spiderrule, &res.Webinfo_snapshot)
		}

	})
	return res, err
}

func TaskNewTask(webid string, onlyfirst string, rebuild string) (id string, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	if webid == "" {
		return id, errors.New("网站id为空")
	}
	if onlyfirst != "no" {
		onlyfirst = "yes"
	}
	if rebuild != "yes" {
		rebuild = "no"
	}
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		idRow, err := db.Query("select UUID_Short()")
		if err != nil {
			panic("获取UUID失败," + err.Error())
		}
		idRow.Next()
		_ = idRow.Scan(&id)
		_, err = db.Exec(fmt.Sprintf(`insert into taskinfo (Taskinfo_id, Taskinfo_webid, Taskinfo_createtime, Taskinfo_onlyfirst, Taskinfo_rebuild, Taskinfo_status) Values('%v', '%v', now(), '%v','%v','%v')`, id, webid, onlyfirst, rebuild, -1))
		if err != nil {
			panic("插入web任务失败," + err.Error())
		}
	})
	return id, err
}

func TaskRun(taskid string) (err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	var webinfo model.WebInfo
	var onlyfirst *string
	var rebuild *string
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		rows, err := db.Query(`select B.Webinfo_id, B.Webinfo_name, B. webinfo_url, B.webinfo_spiderrule, B.webinfo_pagenationrule,A.Taskinfo_onlyfirst,A.Taskinfo_rebuild from taskinfo A
inner join webinfo B on A.taskinfo_webid = B.webinfo_id
where taskinfo_id = ` + taskid)
		if err != nil {
			panic("获取任务信息失败" + err.Error())
		}
		if rows.Next() {
			err = rows.Scan(&webinfo.Webinfo_id, &webinfo.Webinfo_name, &webinfo.Webinfo_url, &webinfo.Webinfo_spiderrule, &webinfo.Webinfo_pagenationrule, &onlyfirst, &rebuild)
			if err != nil {
				panic("webinfo数据有误" + err.Error())
			}
		}
	})
	var spider = ""
	if webinfo.Webinfo_spiderrule != nil {
		spider = *webinfo.Webinfo_spiderrule
	}
	var pagerule = ""
	if webinfo.Webinfo_pagenationrule != nil {
		pagerule = *webinfo.Webinfo_pagenationrule
	}
	var ofbool = true
	if onlyfirst != nil && *onlyfirst == "no" {
		ofbool = false
	}
	var rebuildbool = false
	if rebuild != nil && *rebuild == "yes" {
		rebuildbool = true
	}
	err = TaskRunTask(taskid, webinfo.Webinfo_id, webinfo.Webinfo_url, pagerule, spider, ofbool, rebuildbool)
	return err
}

func TaskRunTask(taskid string, webid string, url string, pagerule string, spiderrule string, onlyfirst bool, rebuild bool) (err error) {
	var starttime = time.Now().Format("2006-01-02 15:04:05")
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		_, err = db.Exec(fmt.Sprintf(`update taskinfo set Taskinfo_starttime = '%v', Taskinfo_status = '%v' where Taskinfo_id = '%v'`, starttime, 0, taskid))
		if err != nil {
			log.Println("更新任务状态失败," + err.Error())
		}
	})
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
				var endtime = time.Now().Format("2006-01-02 15:04:05")
				_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
					_, err1 := db.Exec(fmt.Sprintf(`update taskinfo set Taskinfo_endtime = '%v', Taskinfo_status = '%v' where Taskinfo_id = '%v'`, endtime, -2, taskid))
					if err1 != nil {
						log.Println("更新任务状态失败," + err.Error())
					}
				})
			}
		}
	}()
	var dataList []model.CrawlerData
	var tab = crawler.Instance().NewTab()
	defer tab.Close()
	var data []model.CrawlerData
	err = tab.NavigateEvaluate(url, spiderrule, &dataList)
	if err != nil && err.Error() != "encountered an undefined value" {
		except, ok := err.(*runtime.ExceptionDetails)
		if ok {
			panic(spiderrule + "执行脚本失败," + except.Exception.Description)
		}
		panic(spiderrule + "执行脚本失败," + err.Error())
	}
	dataList = append(dataList, data...)
	if strings.TrimSpace(pagerule) != "" && !onlyfirst {
		if err != nil {
			panic("主页超时，任务失败")
		}
		var page, _ = NewPagenationRule(tab, webid, pagerule, spiderrule)
		var vm = otto.New()
		err := vm.Set("tab", page)
		if err != nil {
			panic(err)
		}
		_, err = vm.Run(page.Pagenationrule)
		if err != nil {
			val, ok:= err.(parser.ErrorList)
			if ok {
				var errMsg = ""
				for _, vv := range val{
					errMsg += "第" + strconv.Itoa(vv.Position.Line)  + "第" + strconv.Itoa(vv.Position.Column) + "列分页规则有错误,信息:" + vv.Message + "\n"
				}
				panic(errMsg)
			}
			panic("分页规则有错误!")
		}
		dataList = append(dataList, *page.Datalist...)
	}
	if len(dataList) == 0 {
		err = errors.New("无捕获的任务")
	}
	go func() {
		var bloom = bloomfilter.NewSqlFilter(webid, 4096, setting.MysqlDataSource, bloomfilter.DefaultHash...)
		defer func() {
			var endtime = time.Now().Format("2006-01-02 15:04:05")
			_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
				_, err = db.Exec(fmt.Sprintf(`update taskinfo set Taskinfo_endtime = '%v', Taskinfo_status = '%v' where Taskinfo_id = '%v'`, endtime, 1, taskid))
				if err != nil {
					log.Println("更新任务状态失败," + err.Error())
				}
			})
			bloom.Write()
		}()
		var ncount = len(dataList)
		var wg = sync.WaitGroup{}
		wg.Add(ncount)
		for i := 0; i < ncount; i++ {
			go func(index int) {
				defer func() {
					if p := recover(); p != nil {
						str, ok := p.(string)
						if ok {
							log.Println(str)
							fmt.Println(str)
						}
					}
					wg.Done()
				}()

				var v = dataList[index]
				newDate, _ := util.ParseAnyTime(v.Date)
				var htmlStr = "" //源html文档
				var status = -1  //爬取结果,-1表示获取网站内容失败,-2代表保存文件失败， -3代表html源码获取成功，但是脱皮失败,-4代表超时，0代表跳过,1代表成功
				var str = ""     //保存在数据库的脱皮数据
				var ha = fmt.Sprintf("%x", md5.Sum([]byte(v.Title+newDate.Format("20060102"))))
				var path = "backup/" + webid + "/"
				var fileName = path + ha + ".html"
				defer func() {
					//保存数据
					_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
						var sqlstr = fmt.Sprintf(`insert into taskres(
							taskres.Taskres_taskid,
							taskres.Taskres_webid,
							taskres.Taskres_pageurl,
							taskres.Taskres_pagetitle,
							taskres.Taskres_pagedate,
							taskres.Taskres_pagetext,
							taskres.Taskres_pagepath,
							taskres.Taskres_status,
							taskres.Taskres_savetime
							)
							VALUES('%v', '%v', '%v', '%v','%v','%v','%v','%v',now())
							ON DUPLICATE KEY UPDATE Taskres_pagetext='%v', taskres.Taskres_pagepath='%v',taskres.Taskres_status='%v',taskres.Taskres_savetime=now()
`, taskid, webid, v.Url, v.Title, newDate.Format("2006-01-02"), str, fileName, status, str, fileName, status)
						_, mer := db.Exec(sqlstr)
						if mer != nil {
							log.Println("数据保存失败!" + mer.Error() + "\n" + sqlstr)
						} else if status == 1 {
							//保存bloom
							bloom.Push([]byte(v.Url))
						}
					})
				}()

				if !rebuild && bloom.Exists([]byte(v.Url)) {
					status = 0
					return
				}

				var aTab = crawler.Instance().NewTab()
				defer aTab.Close()

				htmlStr, err := aTab.Gethtml(v.Url)
				if err != nil {
					if err == crawler.UrlTimeout {
						status = -4
					} else {
						status = -1
					}
				} else {
					err = util.EnsurePath(path)
					if err != nil {
						status = -2
					} else {
						err = ioutil.WriteFile(fileName, []byte(htmlStr), 0644)
						if err != nil {
							status = -2
						} else {
							doc, err1 := html.Parse(strings.NewReader(htmlStr))
							if err1 != nil {
								status = -3
							} else {
								var pellNodes func(*html.Node) error
								pellNodes = func(node *html.Node) error {
									if node.Type == html.ElementNode && (node.Data == "script" || node.Data == "noscript" || node.Data == "a" || node.Data == "style") {
										return nil
									}
									var err error
									txt, er := htmlUtil.GetSelfNodeStr(node)
									if er != nil {
										return errors.New("获取文本节点内容失败" + er.Error())
									}
									str += txt
									for n := node.FirstChild; n != nil; n = n.NextSibling {
										err = pellNodes(n)
										if err != nil {
											return errors.New("获取文本节点内容失败" + err.Error())
										}
									}
									return err
								}
								dberr := pellNodes(doc)
								if dberr != nil {
									status = -3
								} else {
									status = 1
								}
							}
						}
					}
				}
			}(i)
		}
		wg.Wait()
	}()
	return err
}

func TaskGetRes(taskid string) (info model.TaskInfo, res []model.TaskRes, err error) {
	defer func() {
		if p := recover(); p != nil {
			str, ok := p.(string)
			if ok {
				err = errors.New(str)
				log.Println(str)
				fmt.Println(str)
			} else {
				err = errors.New("panic")
			}
		}
	}()
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		inforow := db.QueryRow(fmt.Sprintf(`SELECT taskinfo.Taskinfo_id,
			taskinfo.Taskinfo_webid,
			taskinfo.Taskinfo_createtime,
			taskinfo.Taskinfo_rebuild,
			taskinfo.Taskinfo_onlyfirst,
			taskinfo.Taskinfo_starttime,
			taskinfo.Taskinfo_endtime,
			taskinfo.Taskinfo_status FROM taskinfo where Taskinfo_id = '%v'`, taskid))
		err = inforow.Scan(&info.Taskinfo_id, &info.Taskinfo_webid, &info.Taskinfo_createtime, &info.Taskinfo_rebuild, &info.Taskinfo_onlyfirst, &info.Taskinfo_starttime, &info.Taskinfo_endtime, &info.Taskinfo_status)
		if err != nil {
			panic("查询taskinfo失败" + err.Error())
		}
		resrow, err := db.Query(fmt.Sprintf(`select taskres.Taskres_taskid,
			taskres.Taskres_webid,
			taskres.Taskres_pageurl,
			taskres.Taskres_pagetitle,
			taskres.Taskres_pagedate,
			taskres.Taskres_pagetext,
			taskres.Taskres_pagepath,
			taskres.Taskres_savetime,
			taskres.Taskres_status from taskres where Taskres_taskid = '%v'`, taskid))
		if err != nil {
			panic("查询taskres失败," + err.Error())
		}
		for resrow.Next() {
			var tmp model.TaskRes
			err = resrow.Scan(&tmp.Taskres_taskid, &tmp.Taskres_webid, &tmp.Taskres_pageurl, &tmp.Taskres_pagetitle, &tmp.Taskres_pagedate, &tmp.Taskres_pagetext, &tmp.Taskres_pagepath, &tmp.Taskres_savetime, &tmp.Taskres_status)
			//转一下
			tmp.Taskres_pagetext = html.UnescapeString(tmp.Taskres_pagetext)
			if err != nil {
				log.Println(err.Error())
			}
			res = append(res, tmp)
		}

	})
	return info, res, err
}
