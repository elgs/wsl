package wsl

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/elgs/optional"
)

func extractParamsFromMap(m map[string]any) []any {
	if params, ok := m["params"].([]any); ok {
		return params
	}
	ret := []any{}
	for i := 0; ; i++ {
		if val, ok := m[fmt.Sprint("_", i)]; ok {
			ret = append(ret, val)
		} else {
			break
		}
	}
	return ret
}

func ExtractScriptParamsFromMap(m map[string]any) map[string]any {
	ret := map[string]any{}
	for k, v := range m {
		if strings.HasPrefix(k, "__") {
			vs := v.(string)
			sqlSafe(&vs)
			ret[k] = v
		}
	}
	return ret
}

func valuesToMap(keyLowerCase bool, values ...map[string][]string) map[string]any {
	ret := map[string]any{}
	for _, vs := range values {
		for k, v := range vs {
			if keyLowerCase {
				ret[strings.ToLower(k)] = v[0]
			} else {
				ret[k] = v[0]
			}
		}
	}
	return ret
}

// true if the first character is uppercase, false otherwise
func ShouldExport(sql string) bool {
	if !unicode.IsLetter([]rune(sql)[0]) {
		return false
	}
	return strings.ToUpper(sql[0:1]) == sql[0:1]
}

// returns whether to export the result of this sql statement or not
func SqlNormalize(sql *string) {
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

func SplitSqlLable(sql string) (label string, s string) {
	sql = strings.TrimSpace(sql)
	if strings.HasPrefix(sql, "#") {
		ss := strings.Fields(sql)
		lenSS := len(ss)
		if lenSS == 0 {
			return "", ""
		} else if lenSS == 1 {
			return ss[0][1:], ""
		}
		return ss[0][1:], strings.TrimSpace(sql[len(ss[0]):])
	}
	return "", sql
}

func sqlSafe(s *string) {
	*s = strings.Replace(*s, "'", "''", -1)
	*s = strings.Replace(*s, "--", "", -1)
}

func IsQuery(sql string) bool {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	if strings.HasPrefix(sqlUpper, "SELECT") ||
		strings.HasPrefix(sqlUpper, "SHOW") ||
		strings.HasPrefix(sqlUpper, "DESCRIBE") ||
		strings.HasPrefix(sqlUpper, "EXPLAIN") {
		return true
	}
	return false
}

func ConvertArray[T any, U any](arrayOfInterfaces []T) *optional.Optional[[]U] {
	ret := []U{}
	for _, v := range arrayOfInterfaces {
		if s, ok := any(v).(U); ok {
			ret = append(ret, s)
		} else {
			return optional.New[[]U](nil, errors.New("Failed to convert."))
		}
	}
	return optional.New(&ret, nil)
}

func ConvertMap[T any, U any](data map[string]T) *optional.Optional[map[string]U] {
	if data == nil {
		return optional.New[map[string]U](nil, errors.New("Cannot convert nil."))
	}
	ret := map[string]U{}
	for k, v := range data {
		if s, ok := any(v).(U); ok {
			ret[k] = s
		} else {
			return optional.New[map[string]U](nil, errors.New("Failed to convert."))
		}
	}
	return optional.New(&ret, nil)
}
