module github.com/elgs/wsl

go 1.18

replace github.com/elgs/gorediscache => ../gorediscache

require (
	github.com/elgs/gorediscache v0.0.0-20220715024807-4eb3b9b68434
	github.com/elgs/gosplitargs v0.0.0-20161028071935-a491c5eeb3c8
	github.com/elgs/gosqljson v0.0.0-20220712125658-2f85b34a6a73
)

require (
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
)
