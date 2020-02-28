package model

type HeadCount struct {
	ID             int    `json:"Id,omitempty"`
	ProjectName    string `json:"Project_Name,omitempty"`
	ManagerName    string `json:"Manager_Name,omitempty"`
	BillablesCount int    `json:"Billable"`
	BillingOnHold  int    `json:"Billing_On_Hold"`
	VtCount        int    `json:"Value_Trials"`
	PiICount       int    `json:"Project_Investment"`
	Others         int    `json:"Others"`
	Net            int    `json:"Net"`
	CreatedAt      string `json:"createdat,omitempty"`
	UpdatedAt      string `json:"updatedat,omitempty"`
	IsActive       int    `json:"is_active,omitempty"`
}
