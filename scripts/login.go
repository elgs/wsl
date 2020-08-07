package scripts

var Login = `
-- username or email, password
set @username := ?;
set @password := ?;

set @safe_id := REPLACE(UUID(),'-','');
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

insert INTO USER_SESSION
SELECT @safe_id, ID,USERNAME,EMAIL,@now_utc,@now_utc,'__client_ip','__client_ip'
FROM USER WHERE (USERNAME=@username OR EMAIL=@username) 
AND PASSWORD=ENCRYPT(@password, PASSWORD);

select 
USER_SESSION.ID AS SESSION_ID, USER.EMAIL, USER_FLAG.CODE, USER_FLAG.VALUE
FROM USER 
INNER JOIN USER_SESSION ON USER.ID=USER_SESSION.USER_ID
LEFT OUTER JOIN USER_FLAG ON USER.ID=USER_FLAG.USER_ID
WHERE 
USER_SESSION.ID=@safe_id
AND USER_SESSION.LOGIN_IP='__client_ip';`
