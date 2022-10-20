package entity

const (
	MinWorkload = 9
	MaxWorkload = 16
)

var WorkSpeedHour float32 = 60.0 / 18.0

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}
