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
	if v, err := jqConf.QueryToString("web.http_addr"); err == nil {
		this.Web.HttpAddr = v
	}

	if v, err := jqConf.QueryToString("web.https_addr"); err == nil {
		this.Web.HttpsAddr = v
	} else {
		this.Web.HttpsAddr = "127.0.0.1:2443"
	}

	if v, err := jqConf.QueryToString("web.cert_file"); err == nil {
		this.Web.CertFile = v
	} else {
		this.Web.CertFile = path.Join(path.Dir(this.ConfFile), "cert.pem")
	}

	if v, err := jqConf.QueryToString("web.key_file"); err == nil {
		this.Web.KeyFile = v
	} else {
		this.Web.KeyFile = path.Join(path.Dir(this.ConfFile), "key.pem")
	}

	if v, err := jqConf.QueryToMap("databases"); err == nil {
		this.Databases = v
	}

	if v, err := jqConf.QueryToBool("web.cors"); err == nil {
		this.Web.Cors = v
	} else {
		this.Web.Cors = true
	}

	if v, err := jqConf.QueryToString("mail.mail_host"); err == nil {
		this.Mail.MailHost = v
	}

	if v, err := jqConf.QueryToString("mail.mail_username"); err == nil {
		this.Mail.MailUsername = v
	}

	if v, err := jqConf.QueryToString("mail.mail_password"); err == nil {
		this.Mail.MailPassword = v
	}

	if v, err := jqConf.QueryToMap("app"); err == nil {
		this.App = v
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
