CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE users;
CREATE TABLE IF NOT EXISTS users(
	id          UUID         PRIMARY KEY,
	name        VARCHAR(100) NOT NULL,
	sex         VARCHAR(5)   NOT NULL CHECK(sex IN('male', 'female')),
	age         INTEGER      NOT NULL CHECK(age > 0),
	city_n      VARCHAR(70)  NOT NULL DEFAULT '',
	city_w      VARCHAR(70)  NOT NULL DEFAULT '',
	drop_time   VARCHAR(10)  NOT NULL DEFAULT '10/00',
	t_comfort   INTEGER		 NOT NULL DEFAULT 24,
	t_tol	    INTEGER		 NOT NULL DEFAULT 18,
	t_puh	    INTEGER		 NOT NULL DEFAULT 6,
	temp1		INTEGER		 NOT NULL DEFAULT 0,
	temp2   	VARCHAR(50)  NOT NULL DEFAULT '',

	password    VARCHAR(100) NOT NULL DEFAULT 'bot_user',
	
	telegram_id INTEGER       NOT NULL DEFAULT 0,
	
	created_at  TIMESTAMP    NOT NULL DEFAULT now()
);
