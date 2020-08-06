package scripts

var ResetPassword = `
-- username or email, new password

set @salt := SHA2(RAND(), 512);
set @username=?;

select ID INTO @uid FROM USER WHERE (USER.USERNAME=@username OR USER.EMAIL=@username);

UPDATE USER SET
USER.PASSWORD=ENCRYPT(?, CONCAT('\$6\$rounds=5000$',@salt))
WHERE USER.ID=@uid;

delete FROM USER_SESSION WHERE USER_SESSION.USER_ID=@uid AND USER_SESSION.ID!='__session_id';
`
