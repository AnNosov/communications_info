package repo

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/filereader"
)

func checkEmailProvider(provider string) bool {
	for _, val := range entity.CorrectEmailProviders {
		if val == provider {
			return true
		}
	}
	return false
}

func makeEmailProfile(emailInfoArray []string) (emailData entity.EmailData, err error) {
	emailData.Country = emailInfoArray[0]
	emailData.Provider = emailInfoArray[1]
	emailData.DeliveryTime, err = strconv.Atoi(emailInfoArray[2])
	if err != nil {
		return entity.EmailData{}, fmt.Errorf("voicecall - makeEmailProfile: %w", err)
	}
	return emailData, nil
}

func GetEmailInfo(filePath, separator string) ([]entity.EmailData, error) {

	buf, err := filereader.GetCsvContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("emailinfo - GetEmailInfo: %w", err)
	}

	lines := strings.Split(string(buf), "\n")
	emailInfo := make([]entity.EmailData, 0)

	for _, line := range lines {
		info := strings.Split(line, separator)
		if len(info) != 3 {
			log.Println("GetEmailInfo: incorrect length of line: ", info)
			continue
		}

		if !checkCountry(info[0]) {
			log.Println("GetEmailInfo: incorrect country")
			continue
		}

		if !checkEmailProvider(info[1]) {
			log.Println("GetEmailInfo: incorrect provider")
			continue
		}
		voiceProfile, err := makeEmailProfile(info)
		if err != nil {
			log.Println("GetEmailInfo: ", err.Error())
			continue
		}
		emailInfo = append(emailInfo, voiceProfile)
	}
	return emailInfo, nil
}

func SortEmailbyTimeForCountry(emailDataArray []entity.EmailData) []entity.EmailData {
	for i := 0; i < len(emailDataArray); i++ {
		for j := 0; j < len(emailDataArray)-1-i; j++ {
			if emailDataArray[j].DeliveryTime < emailDataArray[j+1].DeliveryTime {
				emailDataArray[j], emailDataArray[j+1] = emailDataArray[j+1], emailDataArray[j]
			}
		}
	}
	return emailDataArray
}

func GetFastAndSlowEmailProvidersByCountry(emailDataArray []entity.EmailData, country string) ([]entity.EmailData, []entity.EmailData) {
	// only for sorted slice
	var fastEmailProviders, slowEmailProviders []entity.EmailData
	for _, val := range emailDataArray {
		if val.Country == country {
			fastEmailProviders = append(fastEmailProviders, val)
			if len(fastEmailProviders) == 3 {
				break
			}
		}
	}

	for i := len(emailDataArray) - 1; i >= 0; i-- {
		if emailDataArray[i].Country == country {
			slowEmailProviders = append(slowEmailProviders, emailDataArray[i])
			if len(slowEmailProviders) == 3 {
				break
			}
		}
	}
	return fastEmailProviders, slowEmailProviders
}

func GetEmailCountryList(emailDataArray []entity.EmailData) map[string][][]entity.EmailData {
	countryList := make(map[string][][]entity.EmailData)

	for _, val := range emailDataArray {
		if _, ok := countryList[val.Country]; ok {
			continue
		}
		countryList[val.Country] = nil
	}
	return countryList
}
