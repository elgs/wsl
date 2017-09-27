package wsl

import "database/sql"

// Interceptor provides a chance for application to gain more controls over
// the parameters before, the result after, and the error when the query is
// executed. An example would be to provide additional input parameters for
// the query, or convert the result to other formats.
type Interceptor interface {
	Before(tx *sql.Tx, script *string, params map[string]string, headers map[string]string) error
	After(tx *sql.Tx, result *[]interface{}) error
	OnError(err *error) error
}

// GlobalInterceptors is the registry of a list of Interceptors that is
// guaranteed to be executed in its list order on each query.
var GlobalInterceptors []Interceptor

// LocalInterceptors is the registry for each named query (qID) of a list of
// Interceptors that is guaranteed to be executed in its list order on each
// query.
var LocalInterceptors = make(map[string][]Interceptor)
