BEGIN;

-- Step 1: Add the display_name column as nullable
ALTER TABLE public.users
ADD COLUMN display_name character varying(255);

-- Step 2: Update existing records to set display_name as the first part of the name
UPDATE public.users
SET display_name = split_part(name, ' ', 1);

-- Step 3: Alter display_name to be NOT NULL
ALTER TABLE public.users
ALTER COLUMN display_name SET NOT NULL;

COMMIT;
