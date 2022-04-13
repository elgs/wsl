package wsl

import (
	"database/sql"
)

// Interceptor provides a chance for applications to gain more controls over
// the parameters before, the result after, and the error when the query is
// executed. An example would be to provide additional input parameters for
// the query, or convert the result to other formats.
type Interceptor interface {
	Before(tx *sql.Tx, context map[string]any) error
	BeforeEach(tx *sql.Tx, context map[string]any, script *string, sqlParams *[]any, scriptIndex int, scriptLabel string, cumulativeResults map[string]any) (bool, error)
	AfterEach(tx *sql.Tx, context map[string]any, scriptIndex int, scriptLabel string, result any, cumulativeResults map[string]any) error
	After(tx *sql.Tx, context map[string]any, results *any, allResults any) error
	OnError(err *error) error
}

type DefaultInterceptor struct{}

func (this *DefaultInterceptor) Before(tx *sql.Tx, context map[string]any) error {
	return nil
}

func (this *DefaultInterceptor) After(tx *sql.Tx, context map[string]any, results *any, allResults any) error {
	return nil
}

func (this *DefaultInterceptor) BeforeEach(tx *sql.Tx, context map[string]any, script *string, sqlParams *[]any, scriptIndex int, scriptLabel string, cumulativeResults map[string]any) (bool, error) {
	return false, nil
}

func (this *DefaultInterceptor) AfterEach(tx *sql.Tx, context map[string]any, scriptIndex int, scriptLabel string, result any, cumulativeResults map[string]any) error {
	return nil
}

func (this *DefaultInterceptor) OnError(err *error) error {
	return *err
}

func (this *App) RegisterQueryInterceptors(queryId string, is ...Interceptor) {
	for _, i := range is {
		if _, ok := this.queryInterceptors[queryId]; ok {
			this.queryInterceptors[queryId] = append(this.queryInterceptors[queryId], i)
		} else {
			this.queryInterceptors[queryId] = []Interceptor{i}
		}
	}
}

func (this *App) RegisterGlobalInterceptors(is ...Interceptor) {
	for _, i := range is {
		this.globalInterceptors = append(this.globalInterceptors, i)
	}
}
