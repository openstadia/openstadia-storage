package configuration

import "github.com/openstadia/openstadia-storage/models"

type ConfigStore struct {
	config *models.Configuration
}

func CreateConfigStore(config *models.Configuration) ConfigStore {
	return ConfigStore{config: config}
}

func (s *ConfigStore) GetMinioSettings() models.MinioSettings {
	return s.config.MinioSettings
}

func (s *ConfigStore) GetHubSettings() models.HubSettings {
	return s.config.HubSettings
}

func (s *ConfigStore) GetOperationSettings() models.OperationSettings {
	return s.config.OperationRelated
}
