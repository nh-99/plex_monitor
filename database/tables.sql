CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS weather_measurements (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	wind_speed numeric not null,
	wind_speed_3_sec_avg numeric,
	wind_speed_2_min_avg numeric,
	wind_speed_10_min_avg numeric,
	wind_direction numeric not null,
	wind_direction_3_sec_avg numeric,
	wind_direction_2_min_avg numeric,
	wind_direction_10_min_avg numeric,
	air_temperature numeric not null,
	avg_temperature_today numeric not null,
	heat_index numeric not null,
	dew_point numeric not null,
	wind_chill numeric not null,
	relative_humidity numeric not null,
	barometric_pressure numeric not null,
	density_altitude numeric not null,
	latitude numeric not null,
	longitude numeric not null,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by uuid,
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_by uuid,
	deleted_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_by uuid
);