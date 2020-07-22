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
	"log"

	"github.com/elgs/wsl"
	"github.com/elgs/wsl/interceptors"
	"github.com/elgs/wsl/scripts"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app, err := wsl.New()
	if err != nil {
		log.Fatal(err)
	}

	// optionally load built in user management interceptors and scripts
	scripts.LoadBuiltInScripts(app)
	interceptors.RegisterBuiltInInterceptors(app)

	// done manully
	// wsld.RegisterGlobalInterceptors(&interceptors.AuthInterceptor{})
	// wsld.RegisterQueryInterceptors("signup", &interceptors.SignupInterceptor{})
	// ...

	// wsld.Scripts["init"] = scripts.Init
	// wsld.Scripts["signup"] = scripts.Signup
	// ...

	app.Start()
	wsl.Hook()
}
```
### wsld.json
```json
{
   "web": {
      "http_addr": "127.0.0.1:1103",
      "https_addr": "127.0.0.1:1443",
      "cors": true,
      "cert_file": "cert.pem",
      "key_file": "key.pem"
   },
   "databases": {
      "main": {
         "db_type": "mysql",
         "db_url": "root:password@tcp(host:3306)/db"
      },
      "audit": {
         "db_type": "mysql",
         "db_url": "root:password@tcp(host:3306)/db"
      }
   },
   "mail": {
      "mail_host": "host:587",
      "mail_username": "mail",
      "mail_password": "password",
      "mail_from": "noreply@host"
   },
   "app": {
      "foo": "bar"
   }
}
```
