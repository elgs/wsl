package wsl

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"
)

func Hook() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-sigs:
				fmt.Println(sig)
				// cleanup code here
				done <- true
			}
		}
	}()

	<-done
	fmt.Println("Bye!")
}

func extractParamsFromMap(m map[string]interface{}) []interface{} {
	if params, ok := m["params"].([]interface{}); ok {
		return params
	}
	ret := []interface{}{}
	for i := 0; ; i++ {
		if val, ok := m[fmt.Sprint("_", i)]; ok {
			ret = append(ret, val)
		} else {
			break
		}
	}
	return ret
}

func extractScriptParamsFromMap(m map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for k, v := range m {
		if strings.HasPrefix(k, "__") {
			vs := v.(string)
			sqlSafe(&vs)
			ret[k] = v
		}
	}
	return ret
}

func valuesToMap(keyLowerCase bool, values ...map[string][]string) map[string]interface{} {
	ret := map[string]interface{}{}
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
func shouldExport(sql string) bool {
	if !unicode.IsLetter([]rune(sql)[0]) {
		return false
	}
	return strings.ToUpper(sql[0:1]) == sql[0:1]
}

// return whether export the result of this sql statement or not
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

func sqlSafe(s *string) {
	*s = strings.Replace(*s, "'", "''", -1)
	*s = strings.Replace(*s, "--", "", -1)
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

func ConvertInterfaceArrayToStringArray(arrayOfInterfaces []interface{}) ([]string, error) {
	ret := []string{}
	for _, v := range arrayOfInterfaces {
		if s, ok := v.(string); ok {
			ret = append(ret, s)
		} else {
			return nil, errors.New("Failed to convert.")
		}
	}
	return ret, nil
}

var ConvertMapOfInterfacesToMapOfStrings = func(data map[string]interface{}) (map[string]string, error) {
	if data == nil {
		return nil, errors.New("Cannot convert nil.")
	}
	ret := map[string]string{}
	for k, v := range data {
		if v == nil {
			return nil, errors.New("Data contains nil.")
		}
		ret[k] = v.(string)
	}
	return ret, nil
}

var ConvertMapOfStringsToMapOfInterfaces = func(data map[string]string) (map[string]interface{}, error) {
	if data == nil {
		return nil, errors.New("Cannot convert nil.")
	}
	ret := map[string]interface{}{}
	for k, v := range data {
		ret[k] = v
	}
	return ret, nil
}
