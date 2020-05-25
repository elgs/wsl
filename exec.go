package wsl

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/elgs/gosplitargs"
	"github.com/elgs/gosqljson"
)

func (this *WSL) exec(qID string, db *sql.DB, script string, params map[string]string, context map[string]interface{}) (map[string]interface{}, error) {
	queryResult := []interface{}{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	for _, gi := range globalInterceptors {
		err := gi.Before(tx, &script, params, context, this)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, li := range queryInterceptors[qID] {
		err := li.Before(tx, &script, params, context, this)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// log.Println(script)
	if script != "" {
		format := params["format"]
		theCase := params["case"]

		// double underscore
		scriptParams := extractScriptParamsFromMap(params)
		for k, v := range scriptParams {
			script = strings.Replace(script, k, v, -1)
		}

		scriptsArray, err := gosplitargs.SplitArgs(script, ";", true)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// single underscore
		sqlParams := extractParamsFromMap(params)
		totalCount := 0
		for _, s := range scriptsArray {
			sqlNormalize(&s)
			if len(s) == 0 {
				continue
			}
			count, err := gosplitargs.CountSeparators(s, "\\?")
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			if len(sqlParams) < totalCount+count {
				tx.Rollback()
				return nil, errors.New(fmt.Sprint(s, "Incorrect param count. Expected: ", totalCount+count, " actual: ", len(sqlParams)))
			}
			export := shouldExport(s)
			if isQuery(s) {
				if format == "array" {
					header, data, err := gosqljson.QueryTxToArray(tx, theCase, s, sqlParams[totalCount:totalCount+count]...)
					data = append([][]string{header}, data...)
					if err != nil {
						tx.Rollback()
						ierr := this.interceptError(qID, &err)
						if ierr != nil {
							log.Println(ierr)
						}
						return nil, err
					}
					if export {
						queryResult = append(queryResult, data)
					}
				} else {
					data, err := gosqljson.QueryTxToMap(tx, theCase, s, sqlParams[totalCount:totalCount+count]...)
					if err != nil {
						tx.Rollback()
						ierr := this.interceptError(qID, &err)
						if ierr != nil {
							log.Println(ierr)
						}
						return nil, err
					}
					if export {
						queryResult = append(queryResult, data)
					}
				}
			} else {
				rowsAffected, err := gosqljson.ExecTx(tx, s, sqlParams[totalCount:totalCount+count]...)
				if err != nil {
					tx.Rollback()
					ierr := this.interceptError(qID, &err)
					if ierr != nil {
						log.Println(ierr)
					}
					return nil, err
				}
				if export {
					queryResult = append(queryResult, rowsAffected)
				}
			}
			totalCount += count
		}
	}

	result := map[string]interface{}{
		"data": queryResult,
	}

	for _, li := range queryInterceptors[qID] {
		err := li.After(tx, params, result, context, this)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, gi := range globalInterceptors {
		err := gi.After(tx, params, result, context, this)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return result, nil
}

func (this *WSL) interceptError(qID string, err *error) error {
	for _, li := range queryInterceptors[qID] {
		err := li.OnError(err)
		if err != nil {
			return err
		}
	}

	for _, gi := range globalInterceptors {
		err := gi.OnError(err)
		if err != nil {
			return err
		}
	}
	return nil
}
