package model

type Overview struct {
	GrandTotal        int64 `json:"Open"`
	ClosedInTime      int64 `json:"Closed_In-Time"`
	ClosedOutTime     int64 `json:"Closed_Out-Time"`
	InProgressInTime  int64 `json:"Inprogress_In-Time"`
	InProgressOutTime int64 `json:"Inprogress_Out-Time"`
	OnHold            int64 `json:"Onhold"`
}

type Manager struct {
	Name  string  `json:"name"`
	SayDo float64 `json:"saydo"`
}

type Saydo struct {
	WeekStart string  `json:"weekstart"`
	WeekEnd   string  `json:"weekend"`
	Saydo     float32 `json:"saydo"`
}

type Prospects struct {
	ID         string `json:"id"`
	Project    string `json:"project"`
	Manager    string `json:"manager,omitempty"`
	Prospect   string `json:"prospect"`
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	Challenges string `json:"challenges"`
}

type Updates struct {
	ID             string `json:"id"`
	Manager        string `json:"manager,omitempty"`
	Project        string `json:"project"`
	Ups            string `json:"ups"`
	Downs          string `json:"downs"`
	ProjectUpdates string `json:"project_updates"`
	GeneralUpdates string `json:"general_updates"`
	Challenges     string `json:"challenges"`
	NeedHelp       string `json:"need_help"`
	ClientVisits   string `json:"client_visits"`
	TeamSize       int    `json:"team_size"`
	OpenPositions  int    `json:"open_positions"`
	HighPerformer  string `json:"high_performer"`
	LowPerformer   string `json:"low_performer"`
	// Is_active      string `json:"is_active"`
}
type SayDoData struct {
	Data []Saydo `json:"data"`
}

type Data struct {
	Data []Manager `json:"data"`
}
type UpdatesData struct {
	UData []Updates `json:"data"`
}
