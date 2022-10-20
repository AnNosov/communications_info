package entity

var (
	CorrectSMSProviders = [...]string{"Topolo", "Rond", "Kildy"}
)

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}
