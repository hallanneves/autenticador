package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Mysql struct definindo chave
type Mysql struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

// Config struct definindo chave
type Config struct {
	APIKey    string  `json:"api_key"`
	MySQLPool []Mysql `json:"mysql"`
}

// ConfigConecta struct com conf dos servidores e chaves
var ConfigConecta Config
var globalMutex sync.RWMutex

// LerConfig retorna config
func LerConfig(ConfigFile string) error {
	globalMutex.Lock()
	if len(ConfigConecta.APIKey) > 0 {
		globalMutex.Unlock()
		return nil
	}
	jsonFile, err := os.Open(ConfigFile)
	defer jsonFile.Close()
	if err != nil {
		globalMutex.Unlock()
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		globalMutex.Unlock()
		return err
	}

	json.Unmarshal(byteValue, &ConfigConecta)
	globalMutex.Unlock()
	return nil
}
