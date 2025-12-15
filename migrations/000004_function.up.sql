CREATE OR REPLACE FUNCTION orders.orders_create(_order_number text, _order json, _accrual json, _withdrawal json, _user_balance json)
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
        accrual_ins as (
            insert into orders.accruals (accrual_status, accrual_amount, created_at, updated_at, order_id)
            select src.accrual_status, src.accrual_amount, src.created_at, src.updated_at, order_ins.id
                from accrual_row as src, order_ins
        ),
        withdrawal_ins as (
            insert into orders.withdrawals (withdrawal_amount, created_at, updated_at, order_id)
            select src.withdrawal_amount, src.created_at, src.updated_at, order_ins.id
                from withdrawal_row as src, order_ins
        ),
        --*** TODO delete that when auth will get ready ***--
        user_balance_ins as (
            insert into orders.user_balances as ins (accrual_amount, withdrawal_amount, created_at, updated_at, user_id)
            select src.accrual_amount, src.withdrawal_amount, src.updated_at, src.updated_at, src.user_id
                from user_balance_row as src
            on conflict (user_id) do update set
                accrual_amount = ins.accrual_amount + excluded.accrual_amount,
                withdrawal_amount = ins.withdrawal_amount + excluded.withdrawal_amount,
                updated_at = excluded.updated_at
        )
       --*** TODO delete that when auth will get ready ***--
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
                accrued_at = src.accrued_at
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
        user_balance_ins as (
            insert into orders.user_balances as ins (accrual_amount, withdrawal_amount, created_at, updated_at, user_id)
            select src.accrual_amount, src.withdrawal_amount, src.updated_at, src.updated_at, src.user_id
                from user_balance_cte as src
            on conflict (user_id) do update set
                accrual_amount = ins.accrual_amount + excluded.accrual_amount,
                withdrawal_amount = ins.withdrawal_amount + excluded.withdrawal_amount,
                updated_at = excluded.updated_at
        )
        --*** TODO replace user_balance_ins cte with code below when auth will get ready ***--
        -- user_balance_upd as (
        --     update orders.user_balances as upd set
        --         updated_at = src.updated_at,
        --         accrual_amount = upd.accrual_amount + src.accrual_amount,
        --         withdrawal_amount = upd.withdrawal_amount + src.withdrawal_amount
        --     from user_balance_cte as src
        --     where upd.user_id = src.user_id
        -- )
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