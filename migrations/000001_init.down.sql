DROP FUNCTION orders.orders_create(text, json, json, json, json);
DROP FUNCTION orders.orders_update_accruals(json, json, json);
DROP FUNCTION orders.orders_list_accruals_by_order_status(text, integer);
DROP FUNCTION orders.orders_list_withdrawals_by_user_id(integer);
DROP FUNCTION orders.orders_get_user_balance_by_user_id(integer);

DROP TABLE IF EXISTS orders.user_balances;
DROP TABLE IF EXISTS orders.withdrawals;
DROP TABLE IF EXISTS orders.accruals;
DROP TABLE IF EXISTS orders.orders;

DROP SCHEMA IF EXISTS orders;

DROP FUNCTION lock.locks_acquire(text, text, timestamp with time zone, text);
DROP FUNCTION lock.locks_release(text, text, text);

DROP TABLE IF EXISTS locks;

DROP SCHEMA IF EXISTS lock;