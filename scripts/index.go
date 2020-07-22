package scripts

import "github.com/elgs/wsl"

func LoadBuiltInScripts(app *wsl.WSL) {
	app.Scripts["init"] = Init
	app.Scripts["signup"] = Signup
	app.Scripts["login"] = Login
	app.Scripts["verify-user"] = VerifyUser
	app.Scripts["logout"] = logoutSql
	app.Scripts["change-password"] = ChangePassword
	app.Scripts["reset-password"] = ResetPassword
	app.Scripts["forget-password-send-code"] = ForgetPasswordSendCode
	app.Scripts["forget-password-reset-password"] = ForgetPasswordResetPassword
}
