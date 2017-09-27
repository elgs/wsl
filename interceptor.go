package wsl

import (
	"database/sql"
	"log"
)

// Interceptor provides a chance for application to gain more controls over
// the parameters before, the result after, and the error when the query is
// executed. An example would be to provide additional input parameters for
// the query, or convert the result to other formats.
type Interceptor interface {
	Before(tx *sql.Tx, script *string, params map[string]string, headers map[string]string) error
	After(tx *sql.Tx, result *[]interface{}) error
	OnError(err *error) error
}

type DefaultInterceptor struct{}

func (this *DefaultInterceptor) Before(tx *sql.Tx, script *string, params map[string]string, headers map[string]string) error {
	log.Println("Default:Before")
	return nil
}
func (this *DefaultInterceptor) After(tx *sql.Tx, result *[]interface{}) error {
	log.Println("Default:After")
	return nil
}
func (this *DefaultInterceptor) OnError(err *error) error {
	log.Println("Default:Error")
	return nil
}

func RegisterToQuery(qID string, i Interceptor) {
	if _, ok := queryInterceptors[qID]; ok {
		queryInterceptors[qID] = append(queryInterceptors[qID], i)
	} else {
		queryInterceptors[qID] = []Interceptor{i}
	}
}

func RegisterToGlobal(i Interceptor) {
	globalInterceptors = append(globalInterceptors, i)
}

// globalInterceptors is the registry of a list of Interceptors that is
// guaranteed to be executed in its list order on each query.
var globalInterceptors []Interceptor

// querynterceptors is the registry for each named query (qID) of a list of
// Interceptors that is guaranteed to be executed in its list order on each
// query.
var queryInterceptors = make(map[string][]Interceptor)
