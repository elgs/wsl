package jobs

import (
	"database/sql"
	"fmt"

	"github.com/elgs/gosplitargs"
	"github.com/elgs/gosqljson"
	"github.com/elgs/wsl"
)

func executeSQL(db *sql.DB, scripts string, sqlParams *[]any, before func(), after func(map[string]any)) func() {
	if sqlParams == nil {
		sqlParams = &[]any{}
	}
	return func() {
		if before != nil {
			before()
		}
		exportedResults := map[string]any{}
		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			return
		}

		scriptsArray, err := gosplitargs.SplitArgs(scripts, ";", true)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}

		totalCount := 0
		for index, s := range scriptsArray {
			label, s := wsl.SplitSqlLable(s)
			wsl.SqlNormalize(&s)
			if len(s) == 0 {
				continue
			}

			count, err := gosplitargs.CountSeparators(s, "\\?")
			totalCount += count
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
				return
			}
			if len(*sqlParams) < totalCount {
				tx.Rollback()
				fmt.Println(fmt.Sprint(s, "Incorrect param count. Expected: ", totalCount, " actual: ", len(*sqlParams)))
				return
			}

			localSqlParams := (*sqlParams)[totalCount-count : totalCount]

			resultKey := label
			if resultKey == "" {
				resultKey = fmt.Sprint(index)
			}

			export := wsl.ShouldExport(s)

			if wsl.IsQuery(s) {
				result, err := gosqljson.QueryTxToMap(tx, "lower", s, localSqlParams...)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
				if export {
					exportedResults[resultKey] = result
				}
			} else {
				result, err := gosqljson.ExecTx(tx, s, localSqlParams...)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
				if export {
					exportedResults[resultKey] = result
				}
			}
		}
		tx.Commit()
		if after != nil {
			after(exportedResults)
		}
	}
}
