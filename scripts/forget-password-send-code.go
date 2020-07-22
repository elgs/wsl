package scripts

var ForgetPasswordSendCode = `
-- username or email

set @username := ?;

update USER SET USER.USER_FLAG='forget_password', USER.USER_FLAG_CODE='__forget_password'
WHERE USER.USERNAME=@username OR USER.EMAIL=@username;

select EMAIL FROM USER WHERE (USERNAME=@username OR EMAIL=@username);`
