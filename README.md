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

### SQL Scripts
We create a bunch of SQL scrips in the same directory as `wsld.json`, in this case `/home/pi/wsld/`.

`new_pet.sql`
```sql
INSERT INTO pet (name, age) VALUES(?,?);
```

`list_pets.sql`
```sql
SELECT * FROM pet;
```

Assume we have the `pet` table in `mydb` defined as follows:

```sql
CREATE TABLE `pet` (
  `name` varchar(50) NOT NULL,
  `age` int(11) NOT NULL,
  PRIMARY KEY (`name`)
)
```