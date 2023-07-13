CREATE USER gophermart
    PASSWORD 'P@ssw0rd';

CREATE DATABASE gophermart
    OWNER 'gophermart'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

CREATE USER accrual
    PASSWORD 'P@ssw0rd';

CREATE DATABASE accrual
    OWNER 'accrual'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);
