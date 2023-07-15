-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_event
(
    ts       timestamptz NOT NULL default now(),
    order_id bigint      NOT NULL,
    status   text        NOT NULL
);
create index order_status_event_ts_idx
    on order_status_event (ts);
create index order_status_event_order_id_idx
    on order_status_event (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists order_status_event;
-- +goose StatementEnd
