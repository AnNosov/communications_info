package repo

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/pkg/filereader"
)

func checkVoiceProvider(provider string) bool {
	for _, val := range entity.CorrectVoiceCallProviders {
		if val == provider {
			return true
		}
	}
	return false
}

func makeVoiceProfile(voiceInfoArray []string) (voiceData entity.VoiceCallData, err error) {
	voiceData.Country = voiceInfoArray[0]
	voiceData.Bandwidth = voiceInfoArray[1]
	voiceData.ResponseTime = voiceInfoArray[2]
	voiceData.Provider = voiceInfoArray[3]
	float64Val, err := strconv.ParseFloat(voiceInfoArray[4], 32)
	if err != nil {
		return entity.VoiceCallData{}, fmt.Errorf("voicecall - makeVoiceProfile: %w", err)
	}
	voiceData.ConnectionStability = float32(float64Val)
	voiceData.TTFB, err = strconv.Atoi(voiceInfoArray[5])
	if err != nil {
		return entity.VoiceCallData{}, fmt.Errorf("voicecall - makeVoicProfile: %w", err)
	}
	voiceData.VoicePurity, err = strconv.Atoi(voiceInfoArray[6])
	if err != nil {
		return entity.VoiceCallData{}, fmt.Errorf("voicecall - makeVoicProfile: %w", err)
	}
	voiceData.MedianOfCallsTime, err = strconv.Atoi(voiceInfoArray[7])
	if err != nil {
		return entity.VoiceCallData{}, fmt.Errorf("voicecall - makeVoicProfile: %w", err)
	}
	return voiceData, nil
}

func GetVoiceInfo(filePath, separator string) ([]entity.VoiceCallData, error) {

	buf, err := filereader.GetCsvContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("voicecall - GetVoiceInfo: %w", err)
	}

	lines := strings.Split(string(buf), "\n")
	voiceInfo := make([]entity.VoiceCallData, 0)

	for _, line := range lines {
		info := strings.Split(line, separator)
		if len(info) != 8 {
			log.Println("GetVoiceInfo: incorrect length of line: ", info)
			continue
		}

		if !checkCountry(info[0]) {
			log.Println("GetVoiceInfo: incorrect country: ", info[0])
			continue
		}

		if !checkVoiceProvider(info[3]) {
			log.Println("GetVoiceInfo: incorrect provider: ", info[3])
			continue
		}
		voiceProfile, err := makeVoiceProfile(info)
		if err != nil {
			log.Println("GetVoiceInfo: ", err.Error())
			continue
		}
		voiceInfo = append(voiceInfo, voiceProfile)
	}
	return voiceInfo, nil
}
