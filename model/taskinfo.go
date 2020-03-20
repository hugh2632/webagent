package model

type TaskInfo struct {
	Taskinfo_id         string
	Taskinfo_webid      string
	Taskinfo_createtime string
	Taskinfo_onlyfirst  string
	Taskinfo_rebuild    string
	Taskinfo_starttime  *string
	Taskinfo_endtime    *string
	Taskinfo_status     int
}

type AjaxTaskinfo struct{
	TaskInfo
	Webinfo_name           string
	Webinfo_url            string
}
