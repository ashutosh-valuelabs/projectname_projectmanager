package model

type Resignations struct { //Resignations structure
	ID                int    `json:"id,omitempty"`
	Empname           string `json:"name,omitempty"`
	Project           string `json:"project,omitempty"`
	Manager           string `json:"manager,omitempty"`
	Backfillrequired  string `json:"backfill_required,omitempty"`
	Regrenonregre     string `json:"regrettable,omitempty"`
	Status            string `json:"status,omitempty"`
	Dateofresignation string `json:"date_of_resignation,omitempty"`
	Dateofleaving     string `json:"date_of_leaving,omitempty"`
	CreatedAt         string `json:"createdat,omitempty"`
	UpdatedAt         string `json:"updatedat,omitempty"`
	IsActive          int    `json:"isactive,omitempty"`
}
