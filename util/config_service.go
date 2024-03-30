package util

import (
	"sync"
)

type ConfigService struct {
	config Config
	mu     sync.Mutex
}

var configService *ConfigService
var once sync.Once

func GetConfigService() *ConfigService {
	once.Do(func() {
		configService = &ConfigService{}
	})
	return configService
}

func (cs *ConfigService) SetConfig(config Config) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.config = config
}

func (cs *ConfigService) GetConfig() Config {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.config
}
