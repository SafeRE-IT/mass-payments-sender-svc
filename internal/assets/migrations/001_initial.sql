-- +migrate Up

create table requests (
    id bigint primary key,
    owner text not null,
    status text not null,
    failure_reason text,
    lockup_until timestamp without time zone
);

create table payments (
    id bigserial primary key,
    request_id bigint not null references requests(id) on delete cascade,
    status text not null,
    failure_reason text,
    amount bigint not null,
    destination text not null,
    destination_type text not null,
    tx_body text
);

-- +migrate Down

drop table requests;
drop table payments;