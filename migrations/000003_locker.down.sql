DROP FUNCTION lock.locks_acquire(text, text, timestamp with time zone, text);
DROP FUNCTION lock.locks_release(text, text, text);

DROP TABLE IF EXISTS locks;

DROP SCHEMA IF EXISTS lock;