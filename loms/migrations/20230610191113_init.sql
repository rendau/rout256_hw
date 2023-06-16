-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock
(
    warehouse_id bigint NOT NULL,
    sku          bigint NOT NULL,
    cnt          bigint NOT NULL,

    CONSTRAINT stock_pk
        PRIMARY KEY (warehouse_id, sku)
);
create index stock_warehouse_id_idx
    on stock (warehouse_id);
create index stock_sku_idx
    on stock (sku);

CREATE TABLE IF NOT EXISTS stock_reserve
(
    order_id     bigint NOT NULL,
    warehouse_id bigint NOT NULL,
    sku          bigint NOT NULL,
    cnt          bigint NOT NULL,

    CONSTRAINT stock_reserve_pk
        PRIMARY KEY (order_id, warehouse_id, sku)
);
create index stock_reserve_order_id_idx
    on stock_reserve (order_id);
create index stock_reserve_warehouse_id_idx
    on stock_reserve (warehouse_id);
create index stock_reserve_sku_idx
    on stock_reserve (sku);

CREATE TABLE IF NOT EXISTS ord
(
    id      bigserial
        primary key,
    user_id bigint NOT NULL,
    status  text   not null
);
create index ord_user_id_idx
    on ord (user_id);
create index ord_status_idx
    on ord (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ord;
drop table if exists stock_reserve;
drop table if exists stock;
-- +goose StatementEnd
