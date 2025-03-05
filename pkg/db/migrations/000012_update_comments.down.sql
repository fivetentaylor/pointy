ALTER TABLE comments
DROP COLUMN selection,
DROP COLUMN selection_start,
DROP COLUMN selection_finish,
DROP COLUMN notes;

DELETE FROM public.users WHERE id = '00000000-0000-0000-0000-000000000000';
