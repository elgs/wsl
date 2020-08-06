package interceptors

import "github.com/elgs/wsl"

func RegisterBuiltInInterceptors(app *wsl.WSL) {
	app.RegisterGlobalInterceptors(&AuthInterceptor{})

	app.RegisterQueryInterceptors("signup", &SignupInterceptor{})
	app.RegisterQueryInterceptors("login", &LoginInterceptor{})
	app.RegisterQueryInterceptors("verify-user", &VerifyUserInterceptor{})
	app.RegisterQueryInterceptors("logout", &LogoutInterceptor{})
	app.RegisterQueryInterceptors("session", &SessionInterceptor{})
	app.RegisterQueryInterceptors("forget-password-send-code", &ForgetPasswordSendCodeInterceptor{})
	app.RegisterQueryInterceptors("forget-password-reset-password", &ForgetPasswordResetPasswordInterceptor{})
	app.RegisterQueryInterceptors("reset-password", &ResetPasswordInterceptor{})
	app.RegisterQueryInterceptors("change-password", &ChangePasswordInterceptor{})
}
