package wsl

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/elgs/gojq"
)

type ConfigWeb struct {
	HttpAddr  string
	HttpsAddr string
	Cors      bool
	CertFile  string
	KeyFile   string
}

type ConfigMail struct {
	MailHost     string
	MailUsername string
	MailPassword string
	MailFrom     string
}

// Config structure
type Config struct {
	ConfFile  string
	Web       *ConfigWeb
	Databases map[string]interface{}
	Mail      *ConfigMail
	App       map[string]interface{}
}

func (this *Config) httpEnabled() bool {
	return len(this.Web.HttpAddr) > 0
}

func (this *Config) httpsEnabled() bool {
	return len(this.Web.HttpsAddr) > 0
}

func NewConfig(confFile string) (*Config, error) {
	config := &Config{
		ConfFile:  confFile,
		Web:       &ConfigWeb{},
		Databases: map[string]interface{}{},
		Mail:      &ConfigMail{},
		App:       map[string]interface{}{},
	}
	err := config.LoadConfig()
	if err != nil {
		return config, err
	}
	return config, err
}

func (this *Config) LoadConfig() error {
	jqConf, err := gojq.NewFileQuery(this.ConfFile)
	if err != nil {
		return err
	}
	v1, err := jqConf.QueryToString("web.http_addr")
	if err == nil {
		this.Web.HttpAddr = v1
	}
	v2, err := jqConf.QueryToString("web.https_addr")
	if err == nil {
		this.Web.HttpsAddr = v2
	}
	v3, err := jqConf.QueryToString("web.cert_file")
	if err == nil {
		this.Web.CertFile = v3
	} else {
		// default
		this.Web.CertFile = path.Join(path.Dir(this.ConfFile), "cert.pem")
	}
	v4, err := jqConf.QueryToString("web.key_file")
	if err == nil {
		this.Web.KeyFile = v4
	} else {
		// default
		this.Web.KeyFile = path.Join(path.Dir(this.ConfFile), "key.pem")
	}
	v7, err := jqConf.QueryToMap("databases")
	if err == nil {
		this.Databases = v7
	}
	v8, err := jqConf.QueryToBool("web.cors")
	if err == nil {
		this.Web.Cors = v8
	} else {
		this.Web.Cors = false
	}
	v9, err := jqConf.QueryToString("mail.mail_host")
	if err == nil {
		this.Mail.MailHost = v9
	}
	v10, err := jqConf.QueryToString("mail.mail_username")
	if err == nil {
		this.Mail.MailUsername = v10
	}
	v11, err := jqConf.QueryToString("mail.mail_password")
	if err == nil {
		this.Mail.MailPassword = v11
	}
	v12, err := jqConf.QueryToMap("app")
	if err == nil {
		this.App = v12
	}
	// fmt.Println(this.Web)
	// fmt.Println(this.Db)
	// fmt.Println(this.Mail)
	// fmt.Println(this.App)
	return nil
}

func (this *WSL) LoadScripts(scriptName string) error {
	scriptPath := path.Dir(this.Config.ConfFile)

	return filepath.Walk(scriptPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
			}
			scriptName := strings.TrimSuffix(info.Name(), ".sql")
			this.Scripts[scriptName] = string(data)
			if info.Name() == scriptName {
				return io.EOF
			}
		}
		return nil
	})
}
