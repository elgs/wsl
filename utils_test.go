package wsl

import "testing"

func TestExtractSQLParameter(t *testing.T) {
	testCases := map[string]string{
		"set @var0 := ?":       "var0",
		" set  @var1  :=  ?  ": "var1",
		"set @var2:=?":         "var2",
		"set var3 := ?":        "",
	}

	for k, v := range testCases {
		got := ExtractSQLParameter(k)
		if got != v {
			t.Errorf(`%s; wanted "%s", got "%s"`, k, v, got)
		}
	}
}

func TestSplitSqlLable(t *testing.T) {
	testCases := map[string]string{
		"-- @label:insert":    "insert",
		" --@label: insert ":  "insert",
		" --@label : insert ": "insert",
	}

	for k, v := range testCases {
		label, _ := SplitSqlLable(k)
		if label != v {
			t.Errorf(`%s; wanted "%s", got "%s"`, k, v, label)
		}
	}
}
