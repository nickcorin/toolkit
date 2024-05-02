create table if not exists "cursors" (
    "id" uuid not null,
    "name" varchar(255) not null,
    "sequence" serial not null,
    "created_at" timestamp not null,
    "updated_at" timestamp not null,

    primary key (id)
);