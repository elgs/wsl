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
### wsld.json
```json
{
    "http_addr":"127.0.0.1:8080",
    "db_type":"mysql",
    "db_url": "root:password@tcp(127.0.0.1:3306)/mydb"
}
```

Assume we have the `pet` table in `mydb` defined as follows:

```sql
CREATE TABLE `pet` (
  `name` varchar(50) NOT NULL,
  `age` int(11) NOT NULL,
  PRIMARY KEY (`name`)
)
```