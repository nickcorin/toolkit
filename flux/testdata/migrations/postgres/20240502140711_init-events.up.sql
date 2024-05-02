create table "events" (
    "id" uuid not null,
    "topic" varchar(64) not null,
    "sequence" serial not null,
    "key" varchar(64) not null,
    "timestamp" timestamp not null,

    primary key (id)
);