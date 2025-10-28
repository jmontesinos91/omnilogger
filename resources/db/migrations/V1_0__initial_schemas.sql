CREATE TABLE public.log_messages (
    id integer NOT NULL,
    "message" varchar(255) NULL,
    "lang" varchar(10) NULL,
    UNIQUE(id, "lang")
);

CREATE TABLE public.logs (
    id varchar(36) NOT NULL,
    ip_address varchar(20) NULL,
    client_host varchar(100) NULL,
    "provider" varchar(20) NULL,
    "level" smallint NULL,
    "message" integer NULL,
    "description" varchar(255) NULL,
    "path" varchar(255) NULL,
    "resource" varchar(50) NULL,
    "action" varchar(50) NULL,
    "data" json NULL,
    "old_data" json NULL,
    tenant_id jsonb NULL,
    user_id integer NULL,
    created_at timestamp NOT NULL
);