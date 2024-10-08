package monitors_services

import (
	"fmt"

	"github.com/google/uuid"
	monitors_model "github.com/montinger-com/montinger-server/app/monitors/models"
	monitors_repository "github.com/montinger-com/montinger-server/app/monitors/repositories"
	"github.com/montinger-com/montinger-server/lib/cache"
	"github.com/montinger-com/montinger-server/lib/db"
	"github.com/montinger-com/montinger-server/lib/exceptions"
	"github.com/montinger-com/montinger-server/lib/utilities"
	"github.com/rashintha/logger"
)

type MonitorsService struct {
	monitorsRepo *monitors_repository.MonitorsRepository
}

func NewMonitorsService() *MonitorsService {

	return &MonitorsService{
		monitorsRepo: monitors_repository.NewMonitorsRepository(db.MongoClient),
	}
}

func (s *MonitorsService) GetAll() ([]*monitors_model.Monitor, error) {
	return s.monitorsRepo.GetAll()
}

func (s *MonitorsService) Create(monitor *monitors_model.Monitor) (*monitors_model.MonitorCreateResponse, error) {
	monitor.Status = "active"
	id := uuid.New().String()

	token, err := utilities.GenerateSecret(256)
	if err != nil {
		logger.Errorf("Error generating token: %v", err.Error())
		return nil, err
	}

	monitorCreateResponse := &monitors_model.MonitorCreateResponse{
		ID:     id,
		Name:   monitor.Name,
		Type:   monitor.Type,
		Status: monitor.Status,
		Token:  token,
	}

	err = cache.Set(id, monitorCreateResponse, 3600)
	if err != nil {
		logger.Errorf("Error setting cache: %v", err.Error())
		return nil, err
	}

	monitorCreateResponse.Token = fmt.Sprintf("%v&%v", id, token)

	return monitorCreateResponse, nil
}

func (s *MonitorsService) Register(monitor *monitors_model.MonitorRegisterDTO) (*monitors_model.MonitorRegisterResponse, error) {
	monitorCacheInterface, err := cache.Get(monitor.ID)
	if err != nil {
		logger.Errorf("Error getting cache: %v", err.Error())
		return nil, err
	}

	monitorCacheMap, ok := monitorCacheInterface.(map[string]interface{})
	if !ok {
		logger.Errorln("Error casting cache to map.")
		return nil, exceptions.InvalidToken
	}

	if monitorCacheMap["token"].(string) != monitor.Token {
		logger.Errorln("Invalid token.")
		return nil, exceptions.InvalidToken
	}

	uuid := uuid.New()
	secret, err := utilities.GenerateSecret(128)
	if err != nil {
		logger.Errorf("Error generating api key: %v", err.Error())
		return nil, err
	}

	monitorCreate := &monitors_model.Monitor{
		Name:   monitorCacheMap["name"].(string),
		Type:   monitorCacheMap["type"].(string),
		Status: monitorCacheMap["status"].(string),
		APIKey: fmt.Sprintf("%v-%v", uuid, secret),
	}

	err = s.monitorsRepo.Create(monitorCreate)
	if err != nil {
		logger.Errorf("Error creating monitor: %v", err.Error())
		return nil, err
	}

	err = cache.Delete(monitor.ID)
	if err != nil {
		logger.Errorf("Error deleting cache: %v", err.Error())
	}

	monitorRegisterResponse := &monitors_model.MonitorRegisterResponse{
		ID:     monitorCreate.ID,
		Name:   monitorCreate.Name,
		Type:   monitorCreate.Type,
		Status: monitorCreate.Status,
		APIKey: monitorCreate.APIKey,
	}

	return monitorRegisterResponse, nil
}

func (s *MonitorsService) Push(id string, monitor *monitors_model.MonitorPushDTO, apiKey string) error {
	monitorDB, err := s.monitorsRepo.GetByID(id)
	if err != nil {
		logger.Errorf("Error getting monitor: %v", err.Error())
		return err
	}

	if monitorDB.APIKey != apiKey {
		logger.Errorln("Invalid API key.")
		return exceptions.InvalidAPIKey
	}

	monitorDB.LastData = monitors_model.LastData{
		CPUUsage:    monitor.LastData.CPUUsage,
		MemoryUsage: monitor.LastData.MemoryUsage,
	}

	err = s.monitorsRepo.Update(monitorDB)
	if err != nil {
		logger.Errorf("Error updating monitor: %v", err.Error())
		return err
	}

	return nil
}