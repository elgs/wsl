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

func (this *Config) httpEnabled() bool {
	return len(this.HttpAddr) > 0
}

func (this *Config) httpsEnabled() bool {
	return len(this.HttpsAddr) > 0 && len(this.CertFile) > 0 && len(this.KeyFile) > 0
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
	} else {
		return err
	}
	v2, err := jqConf.QueryToString("https_addr")
	if err == nil {
		this.HttpsAddr = v2
	} else {
		return err
	}
	v3, err := jqConf.QueryToString("cert_file")
	if err == nil {
		this.CertFile = v3
	} else {
		return err
	}
	v4, err := jqConf.QueryToString("key_file")
	if err == nil {
		this.KeyFile = v4
	} else {
		return err
	}
	v5, err := jqConf.QueryToString("conf_file")
	if err == nil {
		this.ConfFile = v5
	} else {
		return err
	}
	v6, err := jqConf.QueryToString("db_url")
	if err == nil {
		this.DbUrl = v6
	} else {
		return err
	}
	return nil
}
