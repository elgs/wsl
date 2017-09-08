package wsl

import "github.com/elgs/gojq"

// Config structure
type Config struct {
	HttpAddr  string
	HttpsAddr string
	CertFile  string
	KeyFile   string
	ConfFile  string
	DbUrl     string
}

func (this *Config) LoadConfig(file string) error {
	jqConf, err := gojq.NewFileQuery(file)
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
	v6, err := jqConf.QueryToString("db_url")
	if err == nil {
		this.DbUrl = v6
	}
	return nil
}
