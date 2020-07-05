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

func (this *WSL) exec(qID string, db *sql.DB, scripts string, params map[string]interface{}, context map[string]interface{}) (interface{}, error) {

	context["scripts"] = &scripts
	context["params"] = params
	context["app"] = this

	queryResult := []interface{}{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	for _, gi := range globalInterceptors {
		err := gi.Before(tx, context)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, li := range queryInterceptors[qID] {
		err := li.Before(tx, context)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// log.Println(script)
	if scripts != "" {
		format := ""
		theCase := ""
		if v, ok := params["format"].(string); ok {
			format = v
		}
		if v, ok := params["case"].(string); ok {
			theCase = v
		}

		// double underscore
		scriptParams := extractScriptParamsFromMap(params)
		for k, v := range scriptParams {
			scripts = strings.Replace(scripts, k, v.(string), -1)
		}

		scriptsArray, err := gosplitargs.SplitArgs(scripts, ";", true)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// single underscore
		sqlParams := extractParamsFromMap(params)
		totalCount := 0
		for index, s := range scriptsArray {

			sqlNormalize(&s)
			if len(s) == 0 {
				continue
			}
			count, err := gosplitargs.CountSeparators(s, "\\?")
			totalCount += count
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			if len(sqlParams) < totalCount {
				tx.Rollback()
				return nil, errors.New(fmt.Sprint(s, "Incorrect param count. Expected: ", totalCount, " actual: ", len(sqlParams)))
			}

			skipSql := false
			for _, li := range queryInterceptors[qID] {
				skip, err := li.BeforeEach(tx, context, &s, sqlParams[totalCount-count:totalCount], index)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				if skip == true {
					skipSql = true
				}
			}

			if skipSql {
				continue
			}

			export := shouldExport(s)
			if isQuery(s) {
				if format == "array" {
					header, data, err := gosqljson.QueryTxToArray(tx, theCase, s, sqlParams[totalCount-count:totalCount]...)
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
					data, err := gosqljson.QueryTxToMap(tx, theCase, s, sqlParams[totalCount-count:totalCount]...)
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
				rowsAffected, err := gosqljson.ExecTx(tx, s, sqlParams[totalCount-count:totalCount]...)
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

			for index, li := range queryInterceptors[qID] {
				err := li.AfterEach(tx, context, queryResult, index)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	var results interface{}
	if len(queryResult) == 0 {
		results = []interface{}{}
	} else if len(queryResult) == 1 {
		results = queryResult[0]
	} else {
		results = queryResult
	}

	for _, li := range queryInterceptors[qID] {
		err := li.After(tx, context, results)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, gi := range globalInterceptors {
		err := gi.After(tx, context, results)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return results, nil
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
