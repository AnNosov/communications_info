package repo

import (
	"encoding/json"
	"fmt"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/http/client"
)

func GetSupportInfo(host, port, methodPath string) ([]entity.SupportData, error) {
	var data []entity.SupportData
	buf, err := client.GetHttpResponse(host, port, methodPath)
	if err != nil {
		return data, fmt.Errorf("supportinfo - GetSupportInfo: %w", err)
	}
	if err := json.Unmarshal(buf, &data); err != nil {
		return data, fmt.Errorf("supportinfo - GetSupportInfo: %w", err)
	}

	return data, nil
}

func GetSupportWorkStatus(data []entity.SupportData) (int, int) {
	var sumActiveTickets int
	for _, val := range data {
		sumActiveTickets += val.ActiveTickets
	}
	responseTime := float32(sumActiveTickets) * entity.WorkSpeedHour
	if sumActiveTickets < entity.MinWorkload {
		return 1, int(responseTime)
	} else if sumActiveTickets > entity.MaxWorkload {
		return 3, int(responseTime)
	}
	return 2, int(responseTime)
}
