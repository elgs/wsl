package scripts

var Signup = `
-- username, email, password

set @safe_id := REPLACE(UUID(),'-','');
set @salt := SHA2(RAND(), 512);
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

INSERT INTO USER SET 
ID=@safe_id, 
USERNAME=?, 
EMAIL=?, 
PASSWORD=ENCRYPT(?, CONCAT('\$6\$rounds=5000$',@salt)), 
USER_FLAG='signup',
USER_FLAG_CODE='__signup',
MODE='',
CREATED_TIME=@now_utc;`
