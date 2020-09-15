-- +migrate Up

create table requests (
    id bigint primary key,
    owner text not null,
    status text not null,
    failure_reason text
);

create table transactions (
    hash text primary key,
    request_id bigint not null references requests(id) on delete cascade,
    body text not null,
    status text not null,
    failure_reason text
);

-- +migrate Down

drop table requests;
drop table transactions;