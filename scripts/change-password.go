package scripts

var ChangePassword = `
-- new password, old password
set @newPassword := ?;
set @oldPassword := ?;

set @salt := SHA2(RAND(), 512);

#update_password
UPDATE USER SET 
USER.PASSWORD=ENCRYPT(@newPassword, CONCAT('\$6\$rounds=5000$',@salt))
WHERE USER.PASSWORD=ENCRYPT(@oldPassword, USER.PASSWORD)
AND EXISTS(
	SELECT 1 FROM USER_SESSION WHERE
	USER.ID=USER_SESSION.USER_ID 
	AND USER_SESSION.ID='__session_id'
);

#delete_session
delete FROM USER_SESSION WHERE USER_ID='__user_id' AND USER_SESSION.ID!='__session_id';
`
