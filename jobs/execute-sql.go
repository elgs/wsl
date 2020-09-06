package jobs

import (
	"database/sql"
	"fmt"

	"github.com/elgs/gosplitargs"
	"github.com/elgs/gosqljson"
	"github.com/elgs/wsl"
)

func executeSQL(db *sql.DB, scripts string, sqlParams *[]interface{}) func() {
	if sqlParams == nil {
		sqlParams = &[]interface{}{}
	}
	return func() {
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
		for _, s := range scriptsArray {
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

			if wsl.IsQuery(s) {
				_, err := gosqljson.QueryTxToMap(tx, "lower", s, localSqlParams...)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
			} else {
				_, err := gosqljson.ExecTx(tx, s, localSqlParams...)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
			}
		}
		tx.Commit()
	}
}
