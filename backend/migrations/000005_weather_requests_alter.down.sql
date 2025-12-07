ALTER TABLE weather_requests
    ALTER COLUMN pressure TYPE NUMERIC(5,2),
    ALTER COLUMN humidity TYPE NUMERIC(5,2),
    ALTER COLUMN wind_speed TYPE NUMERIC(5,2);
