package wsl

import "github.com/elgs/gojq"

// Config structure
type Config struct {
	HttpPort  int
	HttpHost  string // "127.0.0.1"
	HttpsPort int
	HttpsHost string
	CertFile  string
	KeyFile   string
	ConfFile  string
}

func (this *Config) LoadConfig(file string) error {
	jqConf, err := gojq.NewFileQuery(file)
	if err != nil {
		//ignore
		return err
	}
	v0, err := jqConf.QueryToInt64("http_port")
	if err == nil {
		this.HttpPort = int(v0)
	}
	v1, err := jqConf.QueryToInt64("https_port")
	if err == nil {
		this.HttpsPort = int(v1)
	}
	v2, err := jqConf.QueryToString("https_host")
	if err == nil {
		this.HttpsHost = v2
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
	return nil
}
