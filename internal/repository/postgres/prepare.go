package postgres

const CreateTable = `CREATE TABLE IF NOT EXISTS metrics(
	name VARCHAR PRIMARY KEY,
	type VARCHAR NOT NULL,
	value double precision,
	delta bigint
);

CREATE INDEX IF NOT EXISTS metrics_type ON metrics(type);`
