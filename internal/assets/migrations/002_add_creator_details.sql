-- +migrate Up

alter table payments add column creator_details jsonb;

-- +migrate Down

alter table payments drop column creator_details;


