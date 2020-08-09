package scripts

var ForgetPasswordSendCode = `
-- username or email
set @username := ?;
set @uid := NULL;
set @email := NULL;

set @safe_id := REPLACE(UUID(),'-','');
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

select ID, EMAIL INTO @uid, @email FROM USER WHERE (USER.USERNAME=@username OR USER.EMAIL=@username);

#get_uid_email
select @uid AS UID, @email AS EMAIL;

#change_password
replace INTO USER_FLAG SET
ID=@safe_id,
USER_ID=@uid,
CODE='forget_password',
VALUE='__forget_password',
PRIVATE=1,
CREATED_TIME=@now_utc;`
