package scripts

var ForgetPasswordResetPassword = `
-- username or email, new password, verification code
set @username := ?;
set @salt := SHA2(RAND(), 512);

UPDATE USER SET 
USER.PASSWORD=ENCRYPT(?, CONCAT('\$6\$rounds=5000$',@salt)),
USER.USER_FLAG='',
USER.USER_FLAG_CODE=''
WHERE (USER.USERNAME=@username OR USER.EMAIL=@username) 
AND USER.USER_FLAG='forget_password'
AND USER.USER_FLAG_CODE=?;`
