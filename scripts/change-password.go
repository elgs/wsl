package scripts

var ChangePassword = `
-- new password, old password

set @salt := SHA2(RAND(), 512);

UPDATE USER, USER_SESSION SET 
USER.PASSWORD=ENCRYPT(?, CONCAT('\$6\$rounds=5000$',@salt))
WHERE USER.ID=USER_SESSION.USER_ID 
AND USER.PASSWORD=ENCRYPT(?, USER.PASSWORD)
AND USER_SESSION.ID='__session_id';

delete FROM USER_SESSION WHERE USER_ID='__user_id' AND USER_SESSION.ID!='__session_id'`
