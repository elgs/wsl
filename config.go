package wsl

import (
	"encoding/json"

	"github.com/elgs/optional"
)

type Config struct {
	Web       *Web                `json:"web"`
	Databases map[string]Database `json:"databases"`
	Mail      *Mail               `json:"mail"`
	Opt       map[string]any      `json:"opt"`
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

func NewConfig(confBytes []byte) *optional.Optional[*Config] {
	var config Config
	err := json.Unmarshal(confBytes, &config)
	return optional.New(&config, err)
}
