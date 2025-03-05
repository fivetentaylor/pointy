-- Create a wrapper function for default access token length
CREATE OR REPLACE FUNCTION generate_default_access_token()
RETURNS VARCHAR LANGUAGE plpgsql AS $$
BEGIN
  RETURN generate_access_link(128); -- default length of 128
END;
$$;

CREATE TABLE one_time_access_tokens (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    token varchar(128) NOT NULL UNIQUE DEFAULT generate_default_access_token(),
    expires_at timestamp with time zone NOT NULL,
    is_used boolean NOT NULL DEFAULT FALSE,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
