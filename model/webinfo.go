package model

type WebInfo struct{
	Webinfo_id	uint64
	Webinfo_name	string
	Webinfo_url	string
	Webinfo_spiderrule	*string
	Webinfo_pagenationrule	*string
	Webinfo_snapshot	*string
}

