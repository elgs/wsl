package scripts

import (
	"log"

	"github.com/elgs/wsl"
)

func LoadBuiltInScripts(app *wsl.App) {
	// app.Scripts["init"] = Init
	// app.Scripts["signup"] = Signup
	// app.Scripts["login"] = Login
	// app.Scripts["verify-user"] = VerifyUser
	// app.Scripts["logout"] = logoutSql
	// app.Scripts["change-password"] = ChangePassword
	// app.Scripts["reset-password"] = ResetPassword
	// app.Scripts["forget-password-send-code"] = ForgetPasswordSendCode
	// app.Scripts["forget-password-reset-password"] = ForgetPasswordResetPassword
	test, err := wsl.BuildScript(Test)
	if err != nil {
		log.Fatal(err)
	}
	app.Scripts["test"] = test
}
