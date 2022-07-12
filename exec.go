package wsl

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/elgs/gosqljson"
)

func (this *App) exec(db *sql.DB, script *Script, params map[string]any, context map[string]any) (any, error) {

	context["script"] = &script
	context["params"] = params
	context["app"] = this

	exportedResults := map[string]any{}
	cumulativeResults := map[string]any{}
	var result any

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	for _, gi := range *this.GlobalInterceptors {
		err := gi.Before(tx, context)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, li := range *script.Interceptors {
		err := li.Before(tx, context)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// log.Println(script)
	format := ""
	if v, ok := params["format"].(string); ok {
		format = v
	}

	for _, statement := range *script.Statements {
		if len(statement.Text) == 0 {
			continue
		}
		SqlNormalize(&statement.Text)

		// double underscore
		scriptParams := ExtractScriptParamsFromMap(params)
		for k, v := range scriptParams {
			statement.Text = strings.Replace(statement.Text, k, v.(string), -1)
		}

		sqlParams := []any{}
		if statement.Param != "" {
			if val, ok := params[statement.Param]; ok {
				sqlParams = append(sqlParams, val)
			} else {
				tx.Rollback()
				return nil, fmt.Errorf("Parameter %v not provided.", statement.Param)
			}
		}

		skipSql := false
		for _, li := range *script.Interceptors {
			skip, err := li.BeforeEach(tx, context, &statement, cumulativeResults)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			if skip {
				skipSql = true
			}
		}

		if skipSql {
			continue
		}

		if statement.IsQuery {
			if format == "array" {
				header, data, err := gosqljson.QueryToArray(tx, gosqljson.Lower, statement.Text, sqlParams...)
				data = append([][]string{header}, data...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				cumulativeResults[statement.Label] = data
				if statement.ShouldExport {
					exportedResults[statement.Label] = data
				}
				result = data
			} else {
				result, err = gosqljson.QueryToMap(tx, gosqljson.Lower, statement.Text, sqlParams...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				cumulativeResults[statement.Label] = result
				if statement.ShouldExport {
					exportedResults[statement.Label] = result
				}
			}
		} else {
			result, err = gosqljson.Exec(tx, statement.Text, sqlParams...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			cumulativeResults[statement.Label] = result
			if statement.ShouldExport {
				exportedResults[statement.Label] = result
			}
		}

		for _, li := range *script.Interceptors {
			err := li.AfterEach(tx, context, &statement, cumulativeResults, result)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// var ret any = exportedResults
	// if len(exportedResults) == 1 {
	// 	ret = exportedResults[0]
	// }

	for _, li := range *script.Interceptors {
		err := li.After(tx, context, exportedResults, cumulativeResults)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, gi := range *this.GlobalInterceptors {
		err := gi.After(tx, context, exportedResults, cumulativeResults)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return exportedResults, nil
}
