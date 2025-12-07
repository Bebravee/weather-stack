ALTER TABLE weather_requests
    ALTER COLUMN pressure TYPE NUMERIC(7,2),
    ALTER COLUMN humidity TYPE NUMERIC(7,2),
    ALTER COLUMN wind_speed TYPE NUMERIC(7,2);
