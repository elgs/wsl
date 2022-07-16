package wsl

import (
	"database/sql"
)

type Interceptor interface {
	Before(tx *sql.Tx, context *Context) error
	After(tx *sql.Tx, context *Context, exportedResults any, cumulativeResults any) error
	BeforeEach(tx *sql.Tx, context *Context, statement *Statement, cumulativeResults map[string]any) (bool, error)
	AfterEach(tx *sql.Tx, context *Context, statement *Statement, cumulativeResults map[string]any, result any) error
}

type DefaultInterceptor struct {
}

func (this *DefaultInterceptor) Before(tx *sql.Tx, context *Context) error {
	return nil
}

func (this *DefaultInterceptor) After(tx *sql.Tx, context *Context, exportedResults any, cumulativeResults any) error {
	return nil
}

func (this *DefaultInterceptor) BeforeEach(tx *sql.Tx, context *Context, statement *Statement, cumulativeResults map[string]any) (bool, error) {
	return false, nil
}

func (this *DefaultInterceptor) AfterEach(tx *sql.Tx, context *Context, statement *Statement, cumulativeResults map[string]any, result any) error {
	return nil
}

func (this *App) RegisterGlobalInterceptors(is ...Interceptor) {
	for _, i := range is {
		*this.GlobalInterceptors = append(*this.GlobalInterceptors, i)
	}
}

func (this *App) RegisterScriptInterceptors(scriptID string, is ...Interceptor) {
	interceptors := this.Interceptors[scriptID]
	if interceptors == nil {
		this.Interceptors[scriptID] = &[]Interceptor{}
	}
	for _, i := range is {
		*this.Interceptors[scriptID] = append(*this.Interceptors[scriptID], i)
	}
}
