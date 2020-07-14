package wsl

import (
	"database/sql"
)

// Interceptor provides a chance for applications to gain more controls over
// the parameters before, the result after, and the error when the query is
// executed. An example would be to provide additional input parameters for
// the query, or convert the result to other formats.
type Interceptor interface {
	Before(tx *sql.Tx, context map[string]interface{}) error
	BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, cumulativeResults interface{}) (bool, error)
	AfterEach(tx *sql.Tx, context map[string]interface{}, result interface{}, cumulativeResults interface{}, scriptIndex int) error
	After(tx *sql.Tx, context map[string]interface{}, results interface{}, allResults interface{}) error
	OnError(err *error) error
}

type DefaultInterceptor struct{}

func (this *DefaultInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {
	params := context["params"].(map[string]interface{})
	params["case"] = "lower"
	return nil
}

func (this *DefaultInterceptor) After(tx *sql.Tx, context map[string]interface{}, results interface{}, allResults interface{}) error {
	return nil
}

func (this *DefaultInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, cumulativeResults interface{}) (bool, error) {
	return false, nil
}

func (this *DefaultInterceptor) AfterEach(tx *sql.Tx, context map[string]interface{}, result interface{}, cumulativeResults interface{}, scriptIndex int) error {
	return nil
}

func (this *DefaultInterceptor) OnError(err *error) error {
	return *err
}

func RegisterQueryInterceptors(queryId string, is ...Interceptor) {
	for _, i := range is {
		if _, ok := queryInterceptors[queryId]; ok {
			queryInterceptors[queryId] = append(queryInterceptors[queryId], i)
		} else {
			queryInterceptors[queryId] = []Interceptor{i}
		}
	}
}

func RegisterGlobalInterceptors(is ...Interceptor) {
	for _, i := range is {
		globalInterceptors = append(globalInterceptors, i)
	}
}

// globalInterceptors is the registry of a list of Interceptors that is
// guaranteed to be executed in its list order on each query.
var globalInterceptors []Interceptor

// queryInterceptors is the registry for each named query (qID) of a list of
// Interceptors that is guaranteed to be executed in its list order on each
// query.
var queryInterceptors = make(map[string][]Interceptor)
