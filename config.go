package wsl

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Web       *Web                `json:"web"`
	Databases map[string]Database `json:"databases"`
	Mail      *Mail               `json:"mail"`
}

type Web struct {
	HttpAddr  string `json:"http_addr"`
	HttpsAddr string `json:"https_addr"`
	Cors      bool   `json:"cors"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
}

type Database struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Mail struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

func NewConfig(configFilePath string) (*Config, error) {
	appFileBytes, err := ioutil.ReadFile(configFilePath)
	var config Config
	err = json.Unmarshal(appFileBytes, &config)
	return &config, err
}

// 			this.Web.KeyFile = path.Join(path.Dir(this.ConfFile), "key.pem")
