package plugins

import "database/sql"
import "log"
import "github.com/elgs/wsl"

type DummyInterceptor struct {
	*wsl.DefaultInterceptor
	Message string
}

func (this *DummyInterceptor) Before(tx *sql.Tx, script *string, params map[string]string, headers map[string]string) error {
	log.Println("Before", this.Message)
	// log.Println(headers)
	// log.Println(params["__client_ip"])
	return nil
}
func (this *DummyInterceptor) After(tx *sql.Tx, result *[]interface{}) error {
	log.Println("After", this.Message)
	return nil
}
func (this *DummyInterceptor) OnError(err *error) error {
	log.Println("error: ", this.Message, *err)
	return nil
}
