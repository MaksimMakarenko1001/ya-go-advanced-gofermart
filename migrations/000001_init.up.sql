CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE IF NOT EXISTS orders.users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS orders.orders (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_number TEXT UNIQUE NOT NULL,
    order_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL,

    FOREIGN KEY (user_id) REFERENCES orders.users(id)
);

CREATE TABLE IF NOT EXISTS orders.accruals (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    accrual_status TEXT NOT NULL,
    accrual_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL,
    accrued_at TIMESTAMPTZ,
    accrue_after TIMESTAMPTZ,

    FOREIGN KEY (order_id) REFERENCES orders.orders(id)
);

CREATE TABLE IF NOT EXISTS orders.withdrawals (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL,

    FOREIGN KEY (order_id) REFERENCES orders.orders(id)
);

CREATE TABLE IF NOT EXISTS orders.user_balances (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    accrual_amount BIGINT NOT NULL,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER UNIQUE NOT NULL,

    FOREIGN KEY (user_id) REFERENCES orders.users(id)
);

CREATE OR REPLACE FUNCTION orders.orders_create_accrual(_order_number text, _order json, _accrual json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _ins_order_id integer;
begin
    if exists (select 1 from orders.orders where order_number = _order_number) then
        return json_build_object(
            'ok', false, 
            'already_exists', true,
			'already_exists_by_user_id', (select user_id from orders.orders where order_number = _order_number)
        );
    end if;

    with 
        order_row as (
            select * from json_populate_record(null::orders.orders, _order)
        ),
        accrual_row as (
            select * from json_populate_record(null::orders.accruals, _accrual)
        ),
        order_ins as (
            insert into orders.orders as ins (order_number, order_status, created_at, updated_at, user_id)
            select src.order_number, src.order_status, src.created_at, src.updated_at, src.user_id
                from order_row as src
            returning ins.id
        ),
        accrual_ins as (
            insert into orders.accruals (accrual_status, accrual_amount, created_at, updated_at, order_id)
            select src.accrual_status, src.accrual_amount, src.created_at, src.updated_at, order_ins.id
                from accrual_row as src, order_ins
        )
    select cte.id from order_ins as cte
        into _ins_order_id
    ;

    return json_build_object('ok', true, 'order_id', _ins_order_id);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_create_withdrawal(_order_number text, _order json, _withdrawal json, _user_balance json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _ins_order_id integer;
begin
    if exists (select 1 from orders.orders where order_number = _order_number) then
        return json_build_object(
            'ok', false, 
            'already_exists', true
        );
    end if;

    with 
        order_row as (
            select * from json_populate_record(null::orders.orders, _order)
        ),
        withdrawal_row as (
            select * from json_populate_record(null::orders.withdrawals, _withdrawal)
        ),
        user_balance_row as (
            select * from json_populate_record(null::orders.user_balances, _user_balance)
        ),
        order_ins as (
            insert into orders.orders as ins (order_number, order_status, created_at, updated_at, user_id)
            select src.order_number, src.order_status, src.created_at, src.updated_at, src.user_id
                from order_row as src
            returning ins.id
        ),
        withdrawal_ins as (
            insert into orders.withdrawals (withdrawal_amount, created_at, updated_at, order_id)
            select src.withdrawal_amount, src.created_at, src.updated_at, order_ins.id
                from withdrawal_row as src, order_ins
        ),
        user_balance_upd as (
            update orders.user_balances as upd set
                updated_at = src.updated_at,
                withdrawal_amount = upd.withdrawal_amount + src.withdrawal_amount
            from user_balance_row as src
            where upd.user_id = src.user_id
        )
    select cte.id from order_ins as cte
        into _ins_order_id
    ;

    return json_build_object('ok', true, 'order_id', _ins_order_id);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_list_accruals_by_order_status(_order_status text, _limit integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with cte as (
        select * from orders.orders where order_status = _order_status
            order by created_at
            limit _limit
    )
    select 
        json_agg(
            json_build_object(
                'order', to_json(cte.*),
                'accrual', to_json(a.*)
            )
        )
    into _res
    from orders.accruals as a
        inner join cte
            on a.order_id = cte.id
    ;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_update_accruals(_accruals json, _orders json, _user_balances json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with
        accrual_cte as (
            select * from json_populate_recordset(null::orders.accruals, _accruals)
        ),
        order_cte as (
            select * from json_populate_recordset(null::orders.orders, _orders)
        ),
        user_balance_cte as (
            select * from json_populate_recordset(null::orders.user_balances, _user_balances)
        ),
        accrual_upd as (
            update orders.accruals as upd set 
                updated_at = src.updated_at,
                accrual_status = src.accrual_status,
                accrual_amount = src.accrual_amount,
                accrued_at = src.accrued_at,
                accrue_after = src.accrue_after
            from accrual_cte as src
            where upd.order_id = src.order_id
        ),
        order_upd as (
            update orders.orders as upd set 
                updated_at = src.updated_at,
                order_status = src.order_status
            from order_cte as src
            where upd.order_number = src.order_number
            returning upd.order_number
        ),
        user_balance_upd as (
            update orders.user_balances as upd set
                updated_at = src.updated_at,
                accrual_amount = upd.accrual_amount + src.accrual_amount
            from user_balance_cte as src
            where upd.user_id = src.user_id
        )
    select json_agg(cte.order_number) from order_upd as cte
        into _res
    ;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_list_accruals_by_user_id(_user_id integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with cte as (
        select * from orders.orders where user_id = _user_id
    )
    select 
        json_agg(
            json_build_object(
                'order', to_json(cte.*),
                'accrual', to_json(a.*)
            )
        )
    into _res
    from orders.accruals as a
        inner join cte
            on a.order_id = cte.id
    ;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_list_withdrawals_by_user_id(_user_id integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    with cte as (
        select * from orders.orders where user_id = _user_id
    )
    select 
        json_agg(
            json_build_object(
                'order', to_json(cte.*),
                'withdrawal', to_json(w.*)
            )
        )
    into _res
    from orders.withdrawals as w
        inner join cte
            on w.order_id = cte.id
    ;

    return coalesce(_res, '[]'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.orders_get_user_balance_by_user_id(_user_id integer)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    select to_json(ub.*)
    into _res
    from orders.user_balances as ub
    where ub.user_id = _user_id
    ;

    return coalesce(_res, '{}'::json);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.users_create(_username text, _user json)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _ins_user_id integer;
begin
    if exists (select 1 from orders.users where username = _username) then
        return json_build_object(
            'ok', false, 
            'already_exists', true
        );
    end if;

    with 
        user_row as (
            select * from json_populate_record(null::orders.users, _user)
        ),
        user_ins as (
            insert into orders.users as ins (username, password_hash, created_at, updated_at)
            select src.username, src.password_hash, src.created_at, src.updated_at
                from user_row as src
            returning ins.id
        ),
        user_balance_ins as (
            insert into orders.user_balances as ins (accrual_amount, withdrawal_amount, created_at, updated_at, user_id)
            select 0, 0, src.created_at, src.updated_at, user_ins.id
                from user_row as src, user_ins
        )
    select cte.id from user_ins as cte
        into _ins_user_id
    ;

    return json_build_object('ok', true, 'user_id', _ins_user_id);
end;
$function$
;

CREATE OR REPLACE FUNCTION orders.users_get_by_username(_username text)
 RETURNS json
 LANGUAGE plpgsql
AS $function$
declare
    _res json;
begin
    select to_json(u.*)
    into _res
    from orders.users as u
    where u.username = _username
    ;

    return coalesce(_res, '{}'::json);
end;
$function$
;

CREATE SCHEMA IF NOT EXISTS lock;

CREATE TABLE IF NOT EXISTS lock.locks (
    key TEXT NOT NULL,
    segment TEXT NOT NULL,
    until TIMESTAMPTZ NOT NULL,
    pid TEXT NOT NULL
);

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
