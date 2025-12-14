CREATE TABLE IF NOT EXISTS orders.orders (
    id SERIAL PRIMARY KEY,
    order_number TEXT UNIQUE NOT NULL,
    order_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS orders.accruals (
    id SERIAL PRIMARY KEY,
    accrual_status TEXT NOT NULL,
    accrual_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL,
    accrued_at TIMESTAMPTZ,

    FOREIGN KEY (order_id) REFERENCES orders.orders(id)
);

CREATE TABLE IF NOT EXISTS orders.withdrawals (
    id SERIAL PRIMARY KEY,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL,

    FOREIGN KEY (order_id) REFERENCES orders.orders(id)
);

CREATE TABLE IF NOT EXISTS orders.user_balances (
    id SERIAL PRIMARY KEY,
    accrual_amount BIGINT NOT NULL,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER UNIQUE NOT NULL
);
