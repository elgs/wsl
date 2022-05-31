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
	BeforeEach(tx *sql.Tx, context map[string]any, statement *Statement, cumulativeResults map[string]any) (bool, error)
	AfterEach(tx *sql.Tx, context map[string]any, statement *Statement, result any, cumulativeResults map[string]any) error
	After(tx *sql.Tx, context map[string]any, results *any, allResults any) error
}

type DefaultInterceptor struct{}

func (this *DefaultInterceptor) Before(tx *sql.Tx, context map[string]any) error {
	return nil
}

func (this *DefaultInterceptor) After(tx *sql.Tx, context map[string]any, results *any, allResults any) error {
	return nil
}

func (this *DefaultInterceptor) BeforeEach(tx *sql.Tx, context map[string]any, statement *Statement, cumulativeResults map[string]any) (bool, error) {
	return false, nil
}

func (this *DefaultInterceptor) AfterEach(tx *sql.Tx, context map[string]any, statement *Statement, result any, cumulativeResults map[string]any) error {
	return nil
}

func (this *App) RegisterGlobalInterceptors(is ...Interceptor) {
	for _, i := range is {
		this.globalInterceptors = append(this.globalInterceptors, i)
	}
}
