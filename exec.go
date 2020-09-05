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

	sqlParams := extractParamsFromMap(params)

	context["scripts"] = &scripts
	context["params"] = params
	context["sqlParams"] = &sqlParams
	context["app"] = this

	params["case"] = "lower"

	exportedResults := map[string]interface{}{}
	cumulativeResults := map[string]interface{}{}
	var result interface{}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	for _, gi := range this.globalInterceptors {
		err := gi.Before(tx, context)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, li := range this.queryInterceptors[qID] {
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

		scriptsArray, err := gosplitargs.SplitArgs(scripts, ";", true)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		totalCount := 0
		for index, s := range scriptsArray {
			label, s := SplitSqlLable(s)
			SqlNormalize(&s)
			if len(s) == 0 {
				continue
			}

			// double underscore
			scriptParams := ExtractScriptParamsFromMap(params)
			for k, v := range scriptParams {
				s = strings.Replace(s, k, v.(string), -1)
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
			localSqlParams := sqlParams[totalCount-count : totalCount]
			for _, li := range this.queryInterceptors[qID] {
				skip, err := li.BeforeEach(tx, context, &s, &localSqlParams, index, label, cumulativeResults)
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

			resultKey := label
			if resultKey == "" {
				resultKey = fmt.Sprint(index)
			}

			export := ShouldExport(s)
			if IsQuery(s) {
				if format == "array" {
					header, data, err := gosqljson.QueryTxToArray(tx, theCase, s, localSqlParams...)
					data = append([][]string{header}, data...)
					if err != nil {
						tx.Rollback()
						ierr := this.interceptError(qID, &err)
						if ierr != nil {
							log.Println(ierr)
						}
						return nil, err
					}
					cumulativeResults[resultKey] = data
					if export {
						exportedResults[resultKey] = data
					}
					result = data
				} else {
					result, err = gosqljson.QueryTxToMap(tx, theCase, s, localSqlParams...)
					if err != nil {
						tx.Rollback()
						ierr := this.interceptError(qID, &err)
						if ierr != nil {
							log.Println(ierr)
						}
						return nil, err
					}
					cumulativeResults[resultKey] = result
					if export {
						exportedResults[resultKey] = result
					}
				}
			} else {
				result, err = gosqljson.ExecTx(tx, s, localSqlParams...)
				if err != nil {
					tx.Rollback()
					ierr := this.interceptError(qID, &err)
					if ierr != nil {
						log.Println(ierr)
					}
					return nil, err
				}
				cumulativeResults[resultKey] = result
				if export {
					exportedResults[resultKey] = result
				}
			}

			for _, li := range this.queryInterceptors[qID] {
				err := li.AfterEach(tx, context, index, label, result, cumulativeResults)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	var ret interface{} = exportedResults
	// if len(exportedResults) == 1 {
	// 	ret = exportedResults[0]
	// }

	for _, li := range this.queryInterceptors[qID] {
		err := li.After(tx, context, &ret, cumulativeResults)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, gi := range this.globalInterceptors {
		err := gi.After(tx, context, &ret, cumulativeResults)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return ret, nil
}

func (this *WSL) interceptError(qID string, err *error) error {
	for _, li := range this.queryInterceptors[qID] {
		err := li.OnError(err)
		if err != nil {
			return err
		}
	}

	for _, gi := range this.globalInterceptors {
		err := gi.OnError(err)
		if err != nil {
			return err
		}
	}
	return nil
}
