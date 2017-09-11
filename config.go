package wsl

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/elgs/gojq"
)

// Config structure
type Config struct {
	HttpAddr   string
	HttpsAddr  string
	Cors       bool
	CertFile   string
	KeyFile    string
	ConfFile   string
	ScriptPath string
	DbType     string
	DbUrl      string
	Scripts    map[string]string
}

func (this *Config) httpEnabled() bool {
	return len(this.HttpAddr) > 0
}

func (this *Config) httpsEnabled() bool {
	return len(this.HttpsAddr) > 0
}

func NewConfig(confFile string) (*Config, error) {
	config := &Config{
		ConfFile: confFile,
	}
	config.Scripts = make(map[string]string)
	err := config.LoadConfig()
	if err != nil {
		return config, err
	}
	err = config.LoadScripts()
	return config, err
}

func (this *Config) LoadConfig() error {
	jqConf, err := gojq.NewFileQuery(this.ConfFile)
	if err != nil {
		return err
	}
	v1, err := jqConf.QueryToString("http_addr")
	if err == nil {
		this.HttpAddr = v1
	}
	v2, err := jqConf.QueryToString("https_addr")
	if err == nil {
		this.HttpsAddr = v2
	}
	v3, err := jqConf.QueryToString("cert_file")
	if err == nil {
		this.CertFile = v3
	} else {
		// default
		this.CertFile = path.Join(path.Dir(this.ConfFile), "cert.pem")
	}
	v4, err := jqConf.QueryToString("key_file")
	if err == nil {
		this.KeyFile = v4
	} else {
		// default
		this.KeyFile = path.Join(path.Dir(this.ConfFile), "key.pem")
	}
	v5, err := jqConf.QueryToString("conf_file")
	if err == nil {
		this.ConfFile = v5
	}
	v6, err := jqConf.QueryToString("script_path")
	if err == nil {
		this.ScriptPath = v6
	}
	v7, err := jqConf.QueryToString("db_type")
	if err == nil {
		this.DbType = v7
	}
	v8, err := jqConf.QueryToString("db_url")
	if err == nil {
		this.DbUrl = v8
	}
	v9, err := jqConf.QueryToBool("cors")
	if err == nil {
		this.Cors = v9
	} else {
		this.Cors = false
	}
	return nil
}

func (this *Config) LoadScripts() error {
	if this.ScriptPath == "" {
		this.ScriptPath = path.Dir(this.ConfFile)
	}
	filepath.Walk(this.ScriptPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
			}
			scriptName := strings.TrimSuffix(strings.ToLower(info.Name()), ".sql")
			this.Scripts[scriptName] = string(data)
		}
		return nil
	})
	return nil
}
