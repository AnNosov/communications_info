package repo

import (
	"fmt"
	"log"
	"strings"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/filereader"
)

func checkSMSProvider(provider string) bool {
	for _, val := range entity.CorrectSMSProviders {
		if val == provider {
			return true
		}
	}
	return false
}

func makeSMSProfile(smsInfoArray []string) (smsData entity.SMSData) {
	smsData.Country = smsInfoArray[0]
	smsData.Bandwidth = smsInfoArray[1]
	smsData.ResponseTime = smsInfoArray[2]
	smsData.Provider = smsInfoArray[3]
	return
}

func GetSMSInfo(filePath, separator string) ([]entity.SMSData, error) {

	buf, err := filereader.GetCsvContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("smsinfo - getSMSInfo: %w", err)
	}

	lines := strings.Split(string(buf), "\n")
	smsInfo := make([]entity.SMSData, 0)

	for _, line := range lines {
		info := strings.Split(line, separator)
		if len(info) != 4 {
			log.Println("getSMSInfo: incorrect length of line")
			continue
		}

		if !checkCountry(info[0]) {
			log.Println("getSMSInfo: incorrect country")
			continue
		}

		if !checkSMSProvider(info[3]) {
			log.Println("getSMSInfo: incorrect provider")
			continue
		}

		smsInfo = append(smsInfo, makeSMSProfile(info))
	}
	return smsInfo, nil
}

func SortSMSbyCountry(smsDataArray []entity.SMSData) []entity.SMSData {
	var ln int
	for i := 0; i < len(smsDataArray)-1; i++ {
		for j := 0; j < len(smsDataArray)-1-i; j++ {
			runeArr1 := []rune(smsDataArray[j].Country)
			runeArr2 := []rune(smsDataArray[j+1].Country)
			if len(runeArr1) >= len(runeArr2) {
				ln = len(runeArr2)
			} else {
				ln = len(runeArr1)
			}

			for k := 0; k < ln-1; k++ {
				if runeArr1[k] > runeArr2[k] {
					smsDataArray[j], smsDataArray[j+1] = smsDataArray[j+1], smsDataArray[j]
					break
				} else if runeArr1[k] == runeArr2[k] {
					continue
				} else {
					break
				}
			}
		}
	}
	return smsDataArray
}

func SortSMSbyProvider(smsDataArray []entity.SMSData) []entity.SMSData {
	var ln int
	for i := 0; i < len(smsDataArray)-1; i++ {
		for j := 0; j < len(smsDataArray)-1-i; j++ {
			runeArr1 := []rune(smsDataArray[j].Provider)
			runeArr2 := []rune(smsDataArray[j+1].Provider)

			if len(runeArr1) >= len(runeArr2) {
				ln = len(runeArr2)
			} else {
				ln = len(runeArr1)
			}
			for k := 0; k < ln-1; k++ {
				if runeArr1[k] > runeArr2[k] {
					smsDataArray[j], smsDataArray[j+1] = smsDataArray[j+1], smsDataArray[j]
					break
				} else if runeArr1[k] == runeArr2[k] {
					continue
				} else {
					break
				}
			}

		}
	}
	return smsDataArray
}
