create table subscription
(
    id bigserial,
    service_name varchar(255) not null,
    price int not null,
    user_id varchar(36) not null,
    start_date DATE not null,
    end_date DATE not null,
    is_deleted boolean not null default false,
    primary key (id)
)