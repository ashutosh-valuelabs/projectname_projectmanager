package model

type OpenPositionAging struct {
	Data []Position `json:"data"`
}

type Pagination struct {
	TotalData  int         `json:"total_data"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"current_page"`
	Data       interface{} `json:"data"`
}
type UserRole struct {
	Role []string `json:"role"`
}

type Project struct {
	Id             int    `json:"Id,omitempty"`
	ProjectName    string `json:"project_name,omitempty"`
	SubProjectName string `json:"subproject_name,omitempty"`
	ManagerName    string `json:"project_manager_name,omitempty"`
	ManagerEmailID string `json:"project_manager_email_id,omitempty"`
}

type Position struct {
	ProjectName    string `json:"projectname"`
	Between0to15   int    `json:"<_15_Days"`
	Between15to30  int    `json:"15_to_30_Days"`
	Between30to60  int    `json:"30_to_60_Days"`
	Between60to90  int    `json:"60_to_90_Days"`
	Between90to120 int    `json:"90_to_120_Days"`
	Greaterthen120 int    `json:">_120_Days"`
}

type Profile struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

//Daily : struct
type Daily struct {
	Id                string `json:"id"`
	Project_name      string `json:"project_name"`
	Designation       string `json:"designation"`
	Ageing            int    `json:"ageing"`
	Type_position     string `json:"new_replacement_resignation"`
	Position          string `json:"positions"`
	Priority          string `json:"priority"`
	Additonal_comment string `json:"additional_comments"`
	L1_due            string `json:"due_for_l1"`
	L2_due            string `json:"due_for_l2"`
	Client_due        string `json:"due_for_client"`
}

//Weekly : struct
type Weekly struct {
	Id                string `json:"id"`
	Project_name      string `json:"project_name"`
	Designation       string `json:"designation"`
	Ageing            int    `json:"ageing"`
	Type_position     string `json:"new_replacement_resignation"`
	Position          string `json:"positions"`
	Priority          string `json:"priority"`
	Additonal_comment string `json:"additional_comments"`
	L1_Happened       string `json:"l1_happened_till_now"`
	L2_Happened       string `json:"l2_happened_till_now"`
	Client_Happened   string `json:"client_happened_till_now"`
}
