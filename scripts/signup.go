package scripts

var Signup = `
-- username, email, password

set @safe_id := REPLACE(UUID(),'-','');
set @safe_id_1 := REPLACE(UUID(),'-','');
set @salt := SHA2(RAND(), 512);
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

INSERT INTO USER SET 
ID=@safe_id, 
USERNAME=?, 
EMAIL=?, 
PASSWORD=ENCRYPT(?, CONCAT('\$6\$rounds=5000$',@salt)), 
MODE='',
CREATED_TIME=@now_utc;

INSERT INTO USER_FLAG SET
ID=@safe_id_1,
USER_ID=@safe_id,
CODE='signup',
VALUE='__signup',
PRIVATE=1,
CREATED_TIME=@now_utc;
`
