ALTER TABLE public.waitlist_users
ADD COLUMN allow_access boolean NOT NULL DEFAULT false;