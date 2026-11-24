
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE servo (
    id INTEGER PRIMARY KEY,
    status_servo INTEGER
);

-- +migrate StatementEnd
