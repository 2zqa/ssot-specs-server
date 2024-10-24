-- This program was written in 2023, devices cannot possibly be updated before then.
ALTER TABLE devices
ADD CONSTRAINT check_last_updated_at
CHECK (last_updated_at >= '2023-01-01 00:00:00+00'::timestamp with time zone
       AND last_updated_at <= NOW());