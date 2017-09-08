package wsl

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/elgs/gosplitargs"
	"github.com/elgs/gosqljson"
)

func exec(tx *sql.Tx, db *sql.DB, script string, scriptParams map[string]string, params map[string]string) ([]interface{}, error) {
	var ret []interface{}

	innerTrans := false
	if tx == nil {
		var err error
		tx, err = db.Begin()
		innerTrans = true
		if err != nil {
			return nil, err
		}
	}

	for k, v := range scriptParams {
		script = strings.Replace(script, k, v, -1)
	}

	scriptsArray, err := gosplitargs.SplitArgs(script, ";", true)
	if err != nil {
		if innerTrans {
			tx.Rollback()
		}
		return nil, err
	}

	array := false
	theCase := ""

	params1 := extractParamsFromMap(params)
	totalCount := 0
	for _, s := range scriptsArray {
		sqlNormalize(&s)
		if len(s) == 0 {
			continue
		}
		count, err := gosplitargs.CountSeparators(s, "\\?")
		if err != nil {
			if innerTrans {
				tx.Rollback()
			}
			return nil, err
		}
		if len(params1) < totalCount+count {
			if innerTrans {
				tx.Rollback()
			}
			return nil, errors.New(fmt.Sprintln("Incorrect param count. Expected: ", totalCount+count, " actual: ", len(params1)))
		}
		isQ := isQuery(s)
		if isQ {
			if array {
				header, data, err := gosqljson.QueryTxToArray(tx, theCase, s, params1[totalCount:totalCount+count]...)
				data = append([][]string{header}, data...)
				if err != nil {
					if innerTrans {
						tx.Rollback()
					}
					return nil, err
				}
				ret = append(ret, data)
			} else {
				data, err := gosqljson.QueryTxToMap(tx, theCase, s, params1[totalCount:totalCount+count]...)
				if err != nil {
					if innerTrans {
						tx.Rollback()
					}
					return nil, err
				}
				ret = append(ret, data)
			}
		} else {
			rowsAffected, err := gosqljson.ExecTx(tx, s, params1[totalCount:totalCount+count]...)
			if err != nil {
				if innerTrans {
					tx.Rollback()
				}
				return nil, err
			}
			ret = append(ret, rowsAffected)
		}
		totalCount += count
	}

	if innerTrans {
		tx.Commit()
	}
	return ret, nil
}

func extractParamsFromMap(map[string]string) []interface{} {
	return []interface{}{}
}

func sqlNormalize(sql *string) {
	*sql = strings.TrimSpace(*sql)
	var ret string
	lines := strings.Split(*sql, "\n")
	for _, line := range lines {
		lineTrimmed := strings.TrimSpace(line)
		if lineTrimmed != "" && !strings.HasPrefix(lineTrimmed, "-- ") {
			ret += line + "\n"
		}
	}
	*sql = ret
}

func isQuery(sql string) bool {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if strings.HasPrefix(sqlUpper, "SELECT") ||
		strings.HasPrefix(sqlUpper, "SHOW") ||
		strings.HasPrefix(sqlUpper, "DESCRIBE") ||
		strings.HasPrefix(sqlUpper, "EXPLAIN") {
		return true
	}
	return false
}
