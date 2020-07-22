package scripts

var Init = `
CREATE TABLE IF NOT EXISTS USER (
  ID char(32) COLLATE utf8_bin NOT NULL,
  USERNAME varchar(50) COLLATE utf8_bin NOT NULL,
  EMAIL varchar(50) CHARACTER SET utf8 NOT NULL,
  PASSWORD varchar(255) CHARACTER SET utf8 NOT NULL,
  USER_FLAG varchar(50) COLLATE utf8_bin NOT NULL,
  USER_FLAG_CODE varchar(50) COLLATE utf8_bin NOT NULL,
  MODE varchar(50) COLLATE utf8_bin NOT NULL,
  TIME_CREATED varchar(50) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (ID),
  UNIQUE KEY USERNAME_UNIQUE (USERNAME),
  UNIQUE KEY EMAIL_UNIQUE (EMAIL),
  KEY idx_USER_TIME_CREATED (TIME_CREATED)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS USER_SESSION (
  ID char(32) COLLATE utf8_bin NOT NULL,
  USER_ID char(32) COLLATE utf8_bin NOT NULL COMMENT ' ',
  USERNAME varchar(50) COLLATE utf8_bin NOT NULL,
  EMAIL varchar(50) COLLATE utf8_bin NOT NULL,
  LOGIN_TIME datetime NOT NULL,
  IP varchar(50) COLLATE utf8_bin NOT NULL,
  SESSION_FLAG varchar(50) COLLATE utf8_bin NOT NULL,
  SESSION_FLAG_CODE varchar(50) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (ID),
  KEY idx_USER_SESSION_USERNAME (USERNAME),
  KEY idx_USER_SESSION_EMAIL (EMAIL),
  KEY idx_USER_SESSION_TIME_LOGIN (LOGIN_TIME),
  KEY idx_USER_SESSION_IP (IP)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;`
