CREATE TYPE attendance_type AS ENUM ('CHECKIN', 'CHECKOUT');

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

CREATE TABLE attendance_records (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(50) NOT NULL,
    details TEXT,
    event_timestamp VARCHAR(50) NOT NULL,
    event_date VARCHAR(50) NOT NULL,
    event_hours VARCHAR(50) NOT NULL,
    record_type attendance_type NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);