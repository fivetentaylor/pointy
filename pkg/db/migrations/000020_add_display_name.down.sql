BEGIN;

-- Remove the display_name column
ALTER TABLE public.users
DROP COLUMN display_name;

COMMIT;
