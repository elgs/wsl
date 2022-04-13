package jobs

import "github.com/elgs/wsl"

func RegisterBuiltInJobs(app *wsl.App) {
	app.RegisterJob("cleanup", cleanup)
}
