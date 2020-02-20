package model

type Active struct { //Structure of total active resignation
	ActiveHeadCount   string `json:"Active Head Count"`
	Billable          string `json:"Billable"`
	BillingOnHold     string `json:"Billing On Hold"`
	ValueTrade        string `json:"Value Trade"`
	ProjectInvestment string `json:"Project Investment"`
	Others            string `json:"Others"`
}
