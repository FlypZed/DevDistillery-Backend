package domain

type ProjectWithTasks struct {
	Project Project `json:"project"`
	Tasks   []Task  `json:"tasks"`
}

type TeamWithUsers struct {
	Team  Team   `json:"team"`
	Users []User `json:"users"`
}
