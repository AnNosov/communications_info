package repo

import (
	"encoding/json"
	"fmt"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/http/client"
)

func GetIncidentInfo(host, port, methodPath string) ([]entity.IncidentData, error) {
	var data []entity.IncidentData
	buf, err := client.GetHttpResponse(host, port, methodPath)
	if err != nil {
		return data, fmt.Errorf("supportinfo - GetSupportInfo: %w", err)
	}
	if err := json.Unmarshal(buf, &data); err != nil {
		return data, fmt.Errorf("supportinfo - GetSupportInfo: %w", err)
	}

	return data, nil
}

func SortIncidetsByActive(data []entity.IncidentData) []entity.IncidentData {

	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j].Status == "Active" && data[j+1].Status != "Active" {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
	return data
}
