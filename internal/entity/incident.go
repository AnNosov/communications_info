package entity

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы: active и closed
}
