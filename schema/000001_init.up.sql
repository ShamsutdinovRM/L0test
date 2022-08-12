CREATE TABLE orders
(
    order_uid          TEXT PRIMARY KEY,
    track_number       TEXT      NOT NULL,
    entry              TEXT      NOT NULL,
    locale             TEXT      NOT NULL,
    internal_signature TEXT      NOT NULL,
    customer_id        TEXT      NOT NULL,
    delivery_service   TEXT      NOT NULL,
    shardkey           TEXT      NOT NULL,
    sm_id              INTEGER   NOT NULL,
    date_created       TIMESTAMP NOT NULL,
    oof_shard          TEXT      NOT NULL
);

CREATE TABLE delivery
(
    order_uid TEXT REFERENCES orders (order_uid) ON DELETE CASCADE,
    name      TEXT NOT NULL,
    phone     TEXT NOT NULL,
    zip       TEXT NOT NULL,
    city      TEXT NOT NULL,
    address   TEXT NOT NULL,
    region    TEXT NOT NULL,
    email     TEXT NOT NULL
);

CREATE TABLE payment
(
    order_uid     TEXT REFERENCES orders (order_uid) ON DELETE CASCADE,
    transaction   TEXT    NOT NULL,
    request_id    TEXT    NOT NULL,
    currency      TEXT    NOT NULL,
    provider      TEXT    NOT NULL,
    amount        INTEGER NOT NULL,
    payment_dt    INTEGER NOT NULL,
    bank          TEXT    NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total   INTEGER NOT NULL,
    custom_fee    INTEGER NOT NULL
);

CREATE TABLE items
(
    order_uid    TEXT REFERENCES orders (order_uid) ON DELETE CASCADE,
    chrt_id      INTEGER NOT NULL,
    track_number TEXT    NOT NULL,
    price        INTEGER NOT NULL,
    rid          TEXT    NOT NULL,
    name         TEXT    NOT NULL,
    sale         INTEGER NOT NULL,
    size         TEXT    NOT NULL,
    total_price  INTEGER NOT NULL,
    nm_id        INTEGER NOT NULL,
    brand        TEXT    NOT NULL,
    status       INTEGER NOT NULL
);

INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,
                    sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 'test', 'meest', '9', '99', '2021-11-26T06:22:19Z',
        '1');
INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ('b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot',
        'test@gmail.com');
INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost,
                     goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', 'b563feb7b2b84b6test', '', 'USD', 'wbpay', '1817', '1637907727', 'alpha', '1500', '317',
        '0');
INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, "size", total_price, nm_id, brand, status)
VALUES ('b563feb7b2b84b6test', '9934930', 'WBILMTESTTRACK', '453', 'ab4219087a764ae0btest', 'Mascaras', '30', '0',
        '317', '2389212', 'Vivienne Sabo', '202');