package model

type Totalretentions struct {
	TotalRetentionInitiated int `json:"totalretentioninitiated"`
	TotalRetained           int `json:"totalretained"`
	//Pages      float64 `json:"pages"`
	Data []Retention `json:"data"`
}

type Retention struct {
	ID                 int    `json:"id,omitempty"`
	ProjectName        string `json:"project_name"`
	ProjectManagerName string `json:"project_manager_name,omitempty"`
	RetentionInitiated int    `json:"retention_initiated"`
	Retained           int    `json:"retained"`
	CreatedAt          string `json:"createdat,omitempty"`
	UpdatedAt          string `json:"updatedat,omitempty"`
	IsActive           int    `json:"isactive,omitempty"`
}
