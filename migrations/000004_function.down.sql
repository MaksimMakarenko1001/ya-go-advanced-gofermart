DROP FUNCTION orders.orders_create(text, json, json, json, json);
DROP FUNCTION orders.orders_update_accruals(json, json, json);

DROP FUNCTION orders.orders_list_accruals_by_order_status(text, integer);

DROP FUNCTION orders.orders_list_withdrawals_by_user_id(integer);
DROP FUNCTION orders.orders_get_user_balance_by_user_id(integer);