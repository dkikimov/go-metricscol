package postgres

const CreateTable = `CREATE TABLE IF NOT EXISTS metrics(
	id serial PRIMARY KEY,
	name VARCHAR NOT NULL,
	type VARCHAR NOT NULL,
	value numeric NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS metrics_type ON metrics(type);`
