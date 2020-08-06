package scripts

var ForgetPasswordResetPassword = `
-- username or email, new password, verification code
set @username := ?;
set @newPassword := ?;
set @vCode := ?;
set @salt := SHA2(RAND(), 512);

UPDATE USER SET 
USER.PASSWORD=ENCRYPT(@newPassword, CONCAT('\$6\$rounds=5000$',@salt))
WHERE (USER.USERNAME=@username OR USER.EMAIL=@username) 
AND EXISTS (
	SELECT 1 FROM USER_FLAG WHERE USER.ID=USER_FLAG.USER_ID
	AND USER_FLAG.CODE='forget_password'
	AND USER_FLAG.VALUE=@vCode
);

DELETE FROM USER_FLAG
WHERE USER_FLAG.CODE='forget_password'
AND USER_FLAG.VALUE=@vCode
AND EXISTS (
	SELECT 1 FROM USER WHERE USER_FLAG.USER_ID=USER.ID
	AND (USER.USERNAME=@username OR USER.EMAIL=@username) 
);`
