CREATE SCHEMA IF NOT EXISTS lock;

CREATE TABLE IF NOT EXISTS lock.locks (
    key TEXT NOT NULL,
    segment TEXT NOT NULL,
    until TIMESTAMPTZ NOT NULL,
    pid TEXT NOT NULL
)

ALTER TABLE lock.locks
    ADD CONSTRAINT locks_pkey PRIMARY KEY (key, segment);

CREATE OR REPLACE FUNCTION lock.locks_acquire(_key text, _segment text, _until timestamp with time zone, _pid text)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
begin
    perform pg_advisory_xact_lock(hashtext('lock_'||_key||_segment));

    if exists(select 1 from lock.locks as l where l.key = _key and l.segment = _segment and l.until > now() and l.pid != _pid) then
        -- other proccess was first, leave
        return json_build_object('ok', false);
    end if;

    insert into lock.locks (key, segment, until, pid)
        select  _key, _segment, _until, _pid
        on conflict (key, segment) do update set
            until = excluded.until,
            pid = excluded.pid
    ;
    
    return json_build_object('ok', true);
end;
$function$
;

CREATE OR REPLACE FUNCTION lock.locks_release(_key text, _segment text, _pid text)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
begin
    perform pg_advisory_xact_lock(hashtext('lock_'||_key||_segment));

    if exists(select 1 from lock.locks as l where l.key = _key and l.segment = _segment and l.until > now() and l.pid != _pid) then
        -- other proccess was first, leave
        return json_build_object('ok', false);
    end if;

    update lock.locks set
        until = now()
        where   key = _key and 
                segment = _segment and
                until > now() 
    ;       
    
    return json_build_object('ok', true);
end;
$function$
;
