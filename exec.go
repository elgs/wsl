package wsl

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/elgs/gosplitargs"
	"github.com/elgs/gosqljson"
)

func (this *WSL) Exec(db *sql.DB, script string, params map[string]string) ([]interface{}, error) {
	var ret []interface{}

	array := false
	theCase := ""

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	scriptParams := extractScriptParamsFromMap(params)
	for k, v := range scriptParams {
		script = strings.Replace(script, k, v, -1)
	}

	scriptsArray, err := gosplitargs.SplitArgs(script, ";", true)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

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
			return nil, errors.New(fmt.Sprintln("Incorrect param count. Expected: ", totalCount+count, " actual: ", len(sqlParams)))
		}
		isQ := isQuery(s)
		if isQ {
			if array {
				header, data, err := gosqljson.QueryTxToArray(tx, theCase, s, sqlParams[totalCount:totalCount+count]...)
				data = append([][]string{header}, data...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				ret = append(ret, data)
			} else {
				data, err := gosqljson.QueryTxToMap(tx, theCase, s, sqlParams[totalCount:totalCount+count]...)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				ret = append(ret, data)
			}
		} else {
			rowsAffected, err := gosqljson.ExecTx(tx, s, sqlParams[totalCount:totalCount+count]...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			ret = append(ret, rowsAffected)
		}
		totalCount += count
	}

	tx.Commit()
	return ret, nil
}
