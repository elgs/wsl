package wsl

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func extractScriptParams(scriptArray *[]string) (*[]string, error) {
	ret := []string{}
	for _, script := range *scriptArray {
		_ = script
		// if script matches `set @variable := ?`
		if true {
			key := ""
			ret = append(ret, key)
		}
	}
	return &ret, nil
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
			var value any
			if len(v) == 0 {
				value = nil
			} else if len(v) == 1 {
				value = v[0]
			} else {
				value = v
			}
			if keyLowerCase {
				ret[strings.ToLower(k)] = value
			} else {
				ret[k] = value
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

func ExtractSQLParameter(statement string) string {
	r := regexp.MustCompile(`(?i)\s*set\s+@(.+?)\s*\:\=\s*\?\s*`)
	m := r.FindStringSubmatch(statement)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func SplitSqlLabel(sqlString string) (label string, s string) {
	sqlString = strings.TrimSpace(sqlString) + "\n"
	labelAndSql := strings.SplitN(sqlString, "\n", 2)
	labelPart := labelAndSql[0]
	sqlPart := labelAndSql[1]
	r := regexp.MustCompile(`(?i)\s*\-\-\s*@label\s*\:\s*(.+)\s*`)
	m := r.FindStringSubmatch(labelPart)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1]), sqlPart
	}
	return "", sqlString
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

func ConvertArray[T any, U any](arrayOfInterfaces []T) (*[]U, error) {
	ret := []U{}
	for _, v := range arrayOfInterfaces {
		if s, ok := any(v).(U); ok {
			ret = append(ret, s)
		} else {
			return nil, errors.New("Failed to convert.")
		}
	}
	return &ret, nil
}

func ConvertMap[T any, U any](data map[string]T) (map[string]U, error) {
	if data == nil {
		return nil, errors.New("Cannot convert nil.")
	}
	ret := map[string]U{}
	for k, v := range data {
		if s, ok := any(v).(U); ok {
			ret[k] = s
		} else {
			return nil, errors.New("Failed to convert.")
		}
	}
	return ret, nil
}
