-- Fill NULL values in `system_content` with a default value (e.g., an empty string) before making it NOT NULL
UPDATE public.prompts
SET system_content = ''
WHERE system_content IS NULL;

-- Make `system_content` column NOT NULL
ALTER TABLE public.prompts ALTER COLUMN system_content SET NOT NULL;

-- Rename `system_content` back to `prompt_text`
ALTER TABLE public.prompts RENAME COLUMN system_content TO prompt_text;

-- Remove the `content` column
ALTER TABLE public.prompts DROP COLUMN content;
