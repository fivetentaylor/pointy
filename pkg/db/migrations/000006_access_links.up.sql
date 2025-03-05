-- Create the function to generate a custom access link
CREATE OR REPLACE FUNCTION generate_access_link(length INT)
RETURNS VARCHAR LANGUAGE plpgsql AS $$
DECLARE
  chars VARCHAR := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  result VARCHAR := '';
  i INT;
BEGIN
  FOR i IN 1..length LOOP
    result := result || substr(chars, floor(random() * length(chars) + 1)::INT, 1);
  END LOOP;
  RETURN result;
END;
$$;

-- Create a wrapper function for default invite link length
CREATE OR REPLACE FUNCTION generate_default_invite_link()
RETURNS VARCHAR LANGUAGE plpgsql AS $$
BEGIN
  RETURN generate_access_link(8); -- default length of 8
END;
$$;

CREATE TABLE shared_document_links (
    id SERIAL PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id),
    inviter_id UUID NOT NULL REFERENCES users(id),
    invitee_email VARCHAR(255) NOT NULL,
    invite_link VARCHAR(8) UNIQUE NOT NULL DEFAULT generate_default_invite_link(),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    is_active BOOLEAN DEFAULT TRUE
);

-- Create an index on the access_link column
CREATE UNIQUE INDEX idx_invite_link ON shared_document_links(invite_link);
