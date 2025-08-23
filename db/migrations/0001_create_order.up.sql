CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(10),
    internal_signature TEXT,
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255),
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS payment (
    transaction VARCHAR(255) PRIMARY KEY,
    order_uid VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(20),
    provider VARCHAR(255),
    amount INT,
    payment_dt TIMESTAMP,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255),
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(255),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT,
    FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);