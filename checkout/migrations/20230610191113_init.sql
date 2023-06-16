-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart
(
    id      bigserial NOT NULL
        primary key,
    user_id bigint    NOT NULL
);
create index cart_user_id_idx
    on cart using hash (user_id);

CREATE TABLE IF NOT EXISTS cart_item
(
    cart_id bigint NOT NULL,
    sku     bigint NOT NULL,
    cnt     bigint NOT NULL,

    CONSTRAINT cart_item_unique UNIQUE (cart_id, sku)
);
create index cart_item_cart_id_idx
    on cart_item using hash (cart_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists cart_item;
drop table if exists cart;
-- +goose StatementEnd
