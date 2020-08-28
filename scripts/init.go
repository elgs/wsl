package scripts

var Init = `
-- root email, root password
set @rootEmail=?;
set @rootPassword=?;

set @safe_id := REPLACE(UUID(),'-','');
set @salt := SHA2(RAND(), 512);
set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

create TABLE IF NOT EXISTS USER (
  ID char(32) COLLATE utf8mb4_general_ci NOT NULL,
  USERNAME varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  EMAIL varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  PASSWORD varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  MODE varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
  CREATED_TIME datetime NOT NULL,
  PRIMARY KEY (ID),
  UNIQUE KEY USERNAME_UNIQUE (USERNAME),
  UNIQUE KEY EMAIL_UNIQUE (EMAIL),
  KEY idx_USER_CREATED_TIME (CREATED_TIME)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

create TABLE USER_SESSION (
	ID char(32)COLLATE utf8mb4_general_ci NOT NULL,
	USER_ID char(32) COLLATE utf8mb4_general_ci NOT NULL,
	USERNAME varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
	EMAIL varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
	LOGIN_TIME datetime NOT NULL,
	LAST_SEEN_TIME datetime NOT NULL,
	LOGIN_IP varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
	LAST_SEEN_IP varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
	ORIGIN VARCHAR(200) COLLATE utf8mb4_general_ci NOT NULL,
	PRIMARY KEY (ID),
	KEY idx_USER_SESSION_USERNAME (USERNAME),
	KEY idx_USER_SESSION_EMAIL (EMAIL),
	KEY idx_USER_SESSION_LOGIN_TIME (LOGIN_TIME),
	KEY idx_USER_SESSION_LOGIN_IP (LOGIN_IP),
	KEY idx_USER_SESSION_LAST_SEEN_TIME (LAST_SEEN_TIME),
	KEY idx_USER_SESSION_LAST_SEEN_IP (LAST_SEEN_IP)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

create TABLE USER_FLAG (
	ID char(32) COLLATE utf8mb4_general_ci NOT NULL,
	USER_ID char(32) COLLATE utf8mb4_general_ci NOT NULL,
	CODE varchar(50) COLLATE utf8mb4_general_ci NOT NULL,
	VALUE varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
	PRIVATE tinyint(1) NOT NULL,
	CREATED_TIME datetime NOT NULL,
	PRIMARY KEY (ID),
	UNIQUE KEY USER_ID (USER_ID,CODE)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

#init
INSERT INTO USER SET 
ID=@safe_id, 
USERNAME='root', 
EMAIL=@rootEmail, 
PASSWORD=ENCRYPT(@rootPassword, CONCAT('\$6\$rounds=5000$',@salt)), 
MODE='root',
CREATED_TIME=@now_utc;
 `
