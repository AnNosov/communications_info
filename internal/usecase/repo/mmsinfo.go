package repo

import (
	"encoding/json"
	"fmt"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/http/client"
)

func checkMMSProvider(provider string) bool {
	for _, val := range entity.CorrectMMSProviders {
		if val == provider {
			return true
		}
	}
	return false
}

func deleteFailMMS(data []entity.MMSData, n int) []entity.MMSData {
	return append(data[:n], data[n+1:]...)
}

func GetMMSInfo(host, port, methodPath string) ([]entity.MMSData, error) {
	var data []entity.MMSData
	buf, err := client.GetHttpResponse(host, port, methodPath)
	if err != nil {
		return data, fmt.Errorf("mmsinfo - GetMMSInfo: %w", err)
	}
	if err := json.Unmarshal(buf, &data); err != nil {
		return data, fmt.Errorf("mmsinfo - GetMMSInfo: %w", err)
	}
	for num, mms := range data {
		if !checkCountry(mms.Country) || !checkMMSProvider(mms.Provider) {
			deleteFailMMS(data, num)
		}
	}
	return data, nil
}

func SortMMSbyCountry(mmsDataArray []entity.MMSData) []entity.MMSData {
	var ln int
	for i := 0; i < len(mmsDataArray)-1; i++ {
		for j := 0; j < len(mmsDataArray)-1-i; j++ {
			runeArr1 := []rune(mmsDataArray[j].Country)
			runeArr2 := []rune(mmsDataArray[j+1].Country)
			if len(runeArr1) >= len(runeArr2) {
				ln = len(runeArr2)
			} else {
				ln = len(runeArr1)
			}
			for k := 0; k < ln; k++ {
				if runeArr1[k] > runeArr2[k] {
					mmsDataArray[j], mmsDataArray[j+1] = mmsDataArray[j+1], mmsDataArray[j]
					break
				} else if runeArr1[k] == runeArr2[k] {
					continue
				} else {
					break
				}
			}
		}
	}
	return mmsDataArray
}

func SortMMSbyProvider(mmsDataArray []entity.MMSData) []entity.MMSData {
	var ln int
	for i := 0; i < len(mmsDataArray)-1; i++ {
		for j := 0; j < len(mmsDataArray)-1-i; j++ {
			runeArr1 := []rune(mmsDataArray[j].Provider)
			runeArr2 := []rune(mmsDataArray[j+1].Provider)

			if len(runeArr1) >= len(runeArr2) {
				ln = len(runeArr2)
			} else {
				ln = len(runeArr1)
			}
			for k := 0; k < ln; k++ {
				if runeArr1[k] > runeArr2[k] {
					mmsDataArray[j], mmsDataArray[j+1] = mmsDataArray[j+1], mmsDataArray[j]
					break
				} else if runeArr1[k] == runeArr2[k] {
					continue
				} else {
					break
				}
			}
		}
	}
	return mmsDataArray
}
