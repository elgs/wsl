package jobs

import "github.com/elgs/wsl"

func RegisterBuiltInJobs(app *wsl.WSL) {
	app.RegisterJob("cleanup", cleanup)
}
