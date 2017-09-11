# WSL - Web SQL Lite
A web interface for executing SQL scripts.

## Getting started
### Installation
```bash
go get -u github.com/elgs/wsl
```

### Simple Example
```golang
package main

import (
	"github.com/elgs/wsl"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	wsld, err := wsl.New("/home/pi/wsld/wsld.json")
	if err != nil {
		log.Println(err)
		return
	}
	wsld.Start()
	wsl.Hook()
}
```
### wsld.json
```json
{
    "http_addr": "127.0.0.1:8080",
    "db_type": "mysql",
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

### Create a pet
Now let's create a new pet with `curl`:

```bash
$ curl "http://127.0.0.1:8080/new_pet?_0=Charlie&_1=1"
```
where `new_pet` is the SQL script name, without the `.sql`, `_0` is for the first parameter in the SQL statement, and `_1` for the second, and so on, if there are more.

The `curl` command above yields the following output:
```
[1]
``` 
which means `1` record is affected.

### List all pets

```bash
$ curl -s "http://127.0.0.1:8080/list_pets"
```

Output as follows:
```json
[[{"age":"1","name":"Charlie"}]]
```

## Advanced

### Full Config File
```json
{
    "http_addr": "127.0.0.1:8080",
    "https_addr": "127.0.0.1:8443",
    "cert_file": "/path/to/cert_file",
    "key_file": "/path/to/key_file",
    "script_path": "/path/to/script_path/",
    "db_type": "mysql",
    "db_url": "root:password@tcp(127.0.0.1:3306)/mydb"
}
```