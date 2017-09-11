# WSL - Web SQL Lite

```golang
package main

import (
	"github.com/elgs/wsl"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	wsld := &wsl.WSL{
		Config: wsl.NewConfig("/home/pi/wsld/wsld.json"),
	}
	wsld.Start()
	wsl.Hook()
}
```