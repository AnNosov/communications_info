package usecase

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/AnNosov/communications_info/config"
	"github.com/AnNosov/communications_info/internal/entity"
	"github.com/AnNosov/communications_info/internal/usecase/repo"
	"golang.org/x/sync/errgroup"
)

type CommunicationAction struct {
	cfg   *config.Config
	cache *repo.ResultInfoCache
}

func New(cfg *config.Config) *CommunicationAction {
	cache := repo.New(5 * time.Minute)
	return &CommunicationAction{
		cfg:   cfg,
		cache: cache,
	}
}

func getSMSDatas(cfg config.FilePath) ([][]entity.SMSData, error) {
	var smsResultData [][]entity.SMSData
	smsDataArr, err := repo.GetSMSInfo(cfg.SmsFilePath, cfg.SmsFileSeparator)
	if err != nil {
		return nil, fmt.Errorf("communications - getSMSDatas: %w", err)
	}
	for i := 0; i < len(smsDataArr); i++ {
		(&smsDataArr[i]).Country = entity.Countries[smsDataArr[i].Country]
		if smsDataArr[i].Country == "" {
			return nil, fmt.Errorf("communications - getSMSDatas: country not found")
		}
	}
	smsByProvider := repo.SortSMSbyProvider(smsDataArr)
	smsResultData = append(smsResultData, smsByProvider)
	smsByCountry := make([]entity.SMSData, len(smsDataArr))
	copy(smsByCountry, smsDataArr)                     // через копию, иначе первоначальный слайс меняется и записывается в итоговый объект 2 раза один и тот же слайс
	smsByCountry = repo.SortSMSbyCountry(smsByCountry) // сортируемый скопированный слайс !!!
	smsResultData = append(smsResultData, smsByCountry)
	return smsResultData, nil
}

func getMMSDatas(cfg config.HttpClient) ([][]entity.MMSData, error) {
	var mmsResultData [][]entity.MMSData
	mmsDataArr, err := repo.GetMMSInfo(cfg.MMSHost, cfg.MMSPort, cfg.MethodPath)
	if err != nil {
		return nil, fmt.Errorf("communications - getMMSDatas: %w", err)
	}
	for i := 0; i < len(mmsDataArr); i++ {
		(&mmsDataArr[i]).Country = entity.Countries[mmsDataArr[i].Country]
		if mmsDataArr[i].Country == "" {
			return nil, fmt.Errorf("communications - getMMSDatas: country not found")
		}
	}
	log.Println(repo.SortMMSbyCountry(mmsDataArr))
	log.Println(repo.SortMMSbyProvider(mmsDataArr))
	mmsByCountry := repo.SortMMSbyCountry(mmsDataArr)
	mmsResultData = append(mmsResultData, mmsByCountry)
	mmsByProvider := make([]entity.MMSData, len(mmsDataArr))
	copy(mmsByProvider, mmsDataArr)
	mmsByProvider = repo.SortMMSbyProvider(mmsByProvider)
	mmsResultData = append(mmsResultData, mmsByProvider)
	return mmsResultData, nil
}

func getEmailDatas(cfg config.FilePath) (map[string][][]entity.EmailData, error) {

	emailDataArr, err := repo.GetEmailInfo(cfg.EmailFilePath, cfg.EmailFileSeparator)
	if err != nil {
		return nil, fmt.Errorf("communications - getEmailDatas: %w", err)
	}
	sortedEmailDataArr := repo.SortEmailbyTimeForCountry(emailDataArr)

	mmsResultMap := repo.GetEmailCountryList(emailDataArr)
	for key := range mmsResultMap {
		fastEmailProviders, slowEmailProviders := repo.GetFastAndSlowEmailProvidersByCountry(sortedEmailDataArr, key)
		mmsResultMap[key] = append(mmsResultMap[key], fastEmailProviders)
		mmsResultMap[key] = append(mmsResultMap[key], slowEmailProviders)
	}
	return mmsResultMap, nil
}

func getBillingDatas(cfg config.FilePath) (entity.BillingData, error) {
	return repo.GetBillingInfo(cfg.BillingFilePath, cfg.BillingFileSeparator)
}

func getVoiceDatas(cfg config.FilePath) ([]entity.VoiceCallData, error) {
	return repo.GetVoiceInfo(cfg.VoiceFilePath, cfg.VoiceFileSeparator)
}

func getSupportDatas(cfg config.HttpClient) ([]int, error) {
	response := make([]int, 0)
	supportDataArr, err := repo.GetSupportInfo(cfg.SupportHost, cfg.SupportPort, cfg.SupportPath)
	if err != nil {
		return nil, fmt.Errorf("communications - getSupportDatas: %w", err)
	}
	status, responseTime := repo.GetSupportWorkStatus(supportDataArr)
	response = append(response, status, responseTime)
	return response, nil
}

func getIncidentDatas(cfg config.HttpClient) ([]entity.IncidentData, error) {
	incidentDataArr, err := repo.GetIncidentInfo(cfg.IncidentHost, cfg.IncidentPort, cfg.IncidentPath)
	if err != nil {
		return nil, fmt.Errorf("communications - getIncidentDatas: %w", err)
	}
	return repo.SortIncidetsByActive(incidentDataArr), nil
}

func getResultSet(cfg config.Config, ctx context.Context) (entity.ResultSetT, error) {

	group, _ := errgroup.WithContext(ctx)
	var resultSet entity.ResultSetT

	group.Go(func() (err error) {
		resultSet.SMS, err = getSMSDatas(cfg.FilePath)
		return
	})

	group.Go(func() (err error) {
		resultSet.MMS, err = getMMSDatas(cfg.HttpClient)
		return
	})

	group.Go(func() (err error) {
		resultSet.Email, err = getEmailDatas(cfg.FilePath)
		return
	})

	group.Go(func() (err error) {
		resultSet.Support, err = getSupportDatas(cfg.HttpClient)
		return
	})

	group.Go(func() (err error) {
		resultSet.VoiceCall, err = getVoiceDatas(cfg.FilePath)
		return
	})

	group.Go(func() (err error) {
		resultSet.Billing, err = getBillingDatas(cfg.FilePath)
		return
	})

	group.Go(func() (err error) {
		resultSet.Incidents, err = getIncidentDatas(cfg.HttpClient)
		return
	})

	if err := group.Wait(); err != nil {
		return entity.ResultSetT{}, err
	}

	return resultSet, nil
}

func (cAction *CommunicationAction) checkCache() (entity.ResultSetT, error) {
	var data entity.ResultSetT
	var err error

	data, err = cAction.cache.Get()
	if err != nil || reflect.DeepEqual(data, entity.ResultSetT{}) { // обновляем кеш, если истекло время или объект пустой
		cAction.cache.Lock()

		data, err = getResultSet(*cAction.cfg, context.Background())
		if err != nil {
			return entity.ResultSetT{}, err
		}
		cAction.cache.Unlock()
		cAction.cache.Set(data, 5*time.Minute)
	}

	return data, nil
}

func (cAction *CommunicationAction) GetCommunicationResult() entity.ResultT {

	var result entity.ResultT
	var err error
	result.Status = true

	result.Data, err = cAction.checkCache()
	if err != nil {
		result.Status = false
		result.Error = "Error on collect data"
		log.Println(err)
	}
	return result
}
