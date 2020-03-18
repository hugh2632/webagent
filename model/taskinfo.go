package model

type TaskInfo struct{
	Taskinfo_id	uint64
	Taskinfo_webid	uint64
	Taskinfo_createtime	string
	Taskinfo_onlyfirst	string
	Taskinfo_rebuild	string
	Taskinfo_starttime	*string
	Taskinfo_endtime	*string
	Taskinfo_status	int
}

