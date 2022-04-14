module github.com/elgs/wsl

go 1.18

replace github.com/elgs/optional => ../optional

require (
	github.com/elgs/gosplitargs v0.0.0-20161028071935-a491c5eeb3c8
	github.com/elgs/gosqljson v0.0.0-20160403005647-027aa4915315
	github.com/elgs/gostrgen v0.0.0-20161222160715-9d61ae07eeae
	github.com/elgs/optional v0.0.0-20220414001220-578858faeba6
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
)

require github.com/stretchr/testify v1.7.1 // indirect
