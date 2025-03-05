-- Rename the `prompt_text` column to `system_content`
ALTER TABLE public.prompts RENAME COLUMN prompt_text TO system_content;

-- Make `system_content` column nullable
ALTER TABLE public.prompts ALTER COLUMN system_content DROP NOT NULL;

-- Add a new `content` column with JSON data type
ALTER TABLE public.prompts ADD COLUMN content_json JSON;
