CREATE TABLE IF NOT EXISTS order.orders (
    id SERIAL PRIMARY KEY,
    order_number TEXT UNIQUE NOT NULL,
    order_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL,
);

CREATE TABLE IF NOT EXISTS order.accruals (
    id SERIAL PRIMARY KEY,
    accrual_status TEXT NOT NULL,
    accrual_amount BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL

    FOREIGN KEY (order_id) REFERENCES order.orders(id)
);

CREATE TABLE IF NOT EXISTS order.withdrawals (
    id SERIAL PRIMARY KEY,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    order_id INTEGER UNIQUE NOT NULL

    FOREIGN KEY (order_id) REFERENCES order.orders(id)
);

CREATE TABLE IF NOT EXISTS order.user_balances (
    id SERIAL PRIMARY KEY,
    accrual_amount BIGINT NOT NULL,
    withdrawal_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER UNIQUE NOT NULL
);
