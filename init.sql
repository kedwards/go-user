CREATE USER dbuser WITH PASSWORD 'dbpassword';

CREATE DATABASE usermgmt;

GRANT ALL PRIVILEGES ON DATABASE usermgmt TO dbuser;

CREATE TABLE public.users (
	email text NOT NULL,
	password text NOT NULL,
	first_name text NOT NULL,
	last_name text NOT NULL,
	role text NOT NULL,
	username text NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	CONSTRAINT users_email_key UNIQUE (username),
	CONSTRAINT users_pkey PRIMARY KEY (email)
);

INSERT INTO public.users (email,password,first_name,last_name,role,username,created_at,updated_at) VALUES
	 ('kedwards@kevinedwards.ca','password','Kevin','Edwards','Admin','kedwards','2024-01-28 01:26:30.681059+00','2024-01-28 01:26:30.681059+00');

