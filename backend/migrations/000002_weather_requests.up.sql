CREATE TABLE IF NOT EXISTS weather_requests(
	city_id     VARCHAR(30)  PRIMARY KEY,
	temperature NUMERIC(5,2) NOT NULL,
	feels_like  NUMERIC(5,2) NOT NULL,
	description VARCHAR(50)  NOT NULL,
	humidity	NUMERIC(5,2) NOT NULL,
	pressure 	NUMERIC(5,2) NOT NULL,
	wind_speed  NUMERIC(5,2) NOT NULL,
	
	created_at  TIMESTAMP    NOT NULL DEFAULT now()
);
