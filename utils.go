package wsl

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func extractParamsFromMap(m map[string]string) []interface{} {
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

func extractScriptParamsFromMap(m map[string]string) map[string]string {
	ret := map[string]string{}
	for k, v := range m {
		if strings.HasPrefix(k, "__") {
			ret[k] = v
		}
	}
	return ret
}

func valuesToMap(values ...map[string][]string) map[string]string {
	ret := map[string]string{}
	for _, vs := range values {
		for k, v := range vs {
			ret[k] = v[0]
		}
	}
	return ret
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
	sqlSafe(sql)
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
