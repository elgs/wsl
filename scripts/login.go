package scripts

var Login = `
-- username or email, username or email, password

set @safe_id := REPLACE(UUID(),'-','');
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');
set @username=?;

insert INTO USER_SESSION
SELECT @safe_id, ID,USERNAME,EMAIL,@now_utc,@now_utc,'__client_ip'
FROM USER WHERE (USERNAME=@username OR EMAIL=@username) 
AND PASSWORD=ENCRYPT(?, PASSWORD);

select USER_SESSION.ID AS SESSION_ID, USER.EMAIL, USER.USER_FLAG, USER.USER_FLAG_CODE 
FROM USER INNER JOIN USER_SESSION ON USER.ID=USER_SESSION.USER_ID
WHERE USER_SESSION.ID=@safe_id
AND USER_SESSION.IP='__client_ip';`
