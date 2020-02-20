package model

type Totalretained struct {
	TotalActiveResignation int                `json:"totalactiveresignation"`
	TotalPip               int                `json:"totalperformanceimproplan"`
	TotalTbr               int                `json:"totaltoberetained"`
	Data                   []Toberetaineddata `json:"data"`
}
type Toberetaineddata struct {
	ID                   int    `json:"id,omitempty"`
	ManagerName          string `json:"managername,omitempty"`
	ProjectName          string `json:"project_name"`
	ActiveResignation    int    `json:"active_resignation"`
	PerformanceImproPlan int    `json:"PIP"`
	ToBeRetained         int    `json:"to_be_retained"`
	IsActive             int    `json:"isactive,omitempty"`
}

// type Toberetainedgetall struct {
// 	Manager              string `json:"manager"`
// 	ActiveResignation    int    `json:"activeresignation"`
// 	PerformanceImproPlan int    `json:"performanceimproplan"`
// 	ToBeRetained         int    `json:"toberetained"`
// }
