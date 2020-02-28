package model

// Pagination : to apply pagenation
// type Pagination struct {
// 	TotalData   int                `json:"total_data,omitempty"`
// 	Limit       int                `json:"limit,omitempty"`
// 	TotalPages  int                `json:"total_pages,omitempty"`
// 	CurrentPage int                `json:"current_page,omitempty"`
// 	AI          []ActionItemClosed `json:"data"`
// }

// ActionItemClosed : action item structure
type ActionItemClosed struct {
	ID          int    `json:"id"`
	ProjectName string `json:"project_name,omitempty"`
	ActionItem  string `json:"action_item,omitempty"`
	Owner       string `json:"owner,omitempty"`
	MeetingDate string `json:"meeting_date,omitempty"`
	TargetDate  string `json:"target_date,omitempty"`
	Status      string `json:"status,omitempty"`
	ClosedDate  string `json:"closed_date,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Flag        string `json:"flag,omitempty"`
}

//Error : error
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
