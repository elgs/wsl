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
	return len(this.HttpsAddr) > 0 && len(this.CertFile) > 0 && len(this.KeyFile) > 0
}

func NewConfig(confFile string) *Config {
	config := &Config{
		ConfFile: confFile,
	}
	config.Scripts = make(map[string]string)
	config.LoadConfig()
	config.LoadScripts()
	return config
}

func (this *Config) LoadConfig() error {
	jqConf, err := gojq.NewFileQuery(this.ConfFile)
	if err != nil {
		//ignore
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
	}
	v4, err := jqConf.QueryToString("key_file")
	if err == nil {
		this.KeyFile = v4
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
