create table if not exists "test" (
    id uuid not null,
    created_at timestamp not null,
    updated_at timestamp not null,

    primary key ("id")
);

