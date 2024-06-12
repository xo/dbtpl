-- table with manual primary key
-- generate insert only (no update, save, upsert, delete)
CREATE TABLE a_manual_table (
  a_text VARCHAR(255)
);

-- table with sequence
CREATE TABLE a_sequence (
  a_seq SERIAL,
  CONSTRAINT a_sequence_pkey PRIMARY KEY (a_seq)
);

CREATE TABLE a_sequence_multi (
  a_seq SERIAL,
  a_text VARCHAR(255),
  CONSTRAINT a_sequence_multi_pkey PRIMARY KEY (a_seq)
);

-- table with primary key
CREATE TABLE a_primary (
  a_key INTEGER PRIMARY KEY
);

CREATE TABLE a_primary_multi (
  a_key INTEGER PRIMARY KEY,
  a_text VARCHAR(255)
);

-- table with composite primary key
CREATE TABLE a_primary_composite (
  a_key1 INTEGER,
  a_key2 INTEGER,
  PRIMARY KEY (a_key1, a_key2)
);

-- table with foreign key
CREATE TABLE a_foreign_key (
  a_key INTEGER REFERENCES a_primary (a_key)
);

-- table with composite foreign key
CREATE TABLE a_foreign_key_composite (
  a_key1 INTEGER,
  a_key2 INTEGER,
  FOREIGN KEY (a_key1, a_key2) REFERENCES a_primary_composite (a_key1, a_key2)
);

-- tables with the same foreign key's name 'foreign_key_1'
-- case for issue #396
CREATE TABLE a_same_fk_name_1 (
    a_fkey INTEGER CONSTRAINT foreign_key_1 REFERENCES a_primary (a_key)
);

CREATE TABLE a_same_fk_name_2 (
    a_fkey INTEGER CONSTRAINT foreign_key_1 REFERENCES a_primary (a_key)
);


-- table with index
CREATE TABLE a_index (
  a_key INTEGER
);

CREATE INDEX a_index_idx ON a_index (a_key);

-- table with composite index
CREATE TABLE a_index_composite (
  a_key1 INTEGER,
  a_key2 INTEGER
);

CREATE INDEX a_index_composite_idx ON a_index_composite (a_key1, a_key2);

-- table with unique index
CREATE TABLE a_unique_index (
  a_key INTEGER UNIQUE
);

-- table with composite unique index
CREATE TABLE a_unique_index_composite (
  a_key1 INTEGER,
  a_key2 INTEGER,
  UNIQUE (a_key1, a_key2)
);

-- table with camelCase naming
CREATE TABLE "aTableNamedWithCamelCase" (
  a_key VARCHAR(255) 
);

-- enum type
CREATE TYPE a_enum AS ENUM (
  'ONE',
  'TWO'
);

/*

a_enum
bigint
bigserial
bit
bit varying
bool
boolean
bpchar
bytea
char
character
character varying
date
decimal
double precision
inet
int
integer
interval
json
jsonb
money
numeric
real
serial
smallint
smallserial
text
time
timestamp
timestamptz
timetz
uuid
varchar
xml

*/

-- table with all field types and all nullable field types
CREATE TABLE a_bit_of_everything (
  a_enum a_enum NOT NULL,
  a_enum_nullable a_enum,
  a_bigint BIGINT NOT NULL,
  a_bigint_nullable BIGINT,
  a_bigserial BIGSERIAL NOT NULL,
  a_bigserial_nullable BIGSERIAL,
  a_bit BIT NOT NULL,
  a_bit_nullable BIT,
  a_bit_varying BIT VARYING NOT NULL,
  a_bit_varying_nullable BIT VARYING,
  a_bool BOOL NOT NULL,
  a_bool_nullable BOOL,
  a_boolean BOOLEAN NOT NULL,
  a_boolean_nullable BOOLEAN,
  a_bpchar BPCHAR NOT NULL,
  a_bpchar_nullable BPCHAR,
  a_bytea BYTEA NOT NULL,
  a_bytea_nullable BYTEA,
  a_char CHAR NOT NULL,
  a_char_nullable CHAR,
  a_character CHARACTER NOT NULL,
  a_character_nullable CHARACTER,
  a_character_varying CHARACTER VARYING NOT NULL,
  a_character_varying_nullable CHARACTER VARYING,
  a_date DATE NOT NULL,
  a_date_nullable DATE,
  a_decimal DECIMAL NOT NULL,
  a_decimal_nullable DECIMAL,
  a_double_precision DOUBLE PRECISION NOT NULL,
  a_double_precision_nullable DOUBLE PRECISION,
  a_inet INET NOT NULL,
  a_inet_nullable INET,
  a_int INT NOT NULL,
  a_int_nullable INT,
  a_integer INTEGER NOT NULL,
  a_integer_nullable INTEGER,
  a_interval INTERVAL NOT NULL,
  a_interval_nullable INTERVAL,
  a_json JSON NOT NULL,
  a_json_nullable JSON,
  a_jsonb JSONB NOT NULL,
  a_jsonb_nullable JSONB,
  a_money MONEY NOT NULL,
  a_money_nullable MONEY,
  a_numeric NUMERIC NOT NULL,
  a_numeric_nullable NUMERIC,
  a_real REAL NOT NULL,
  a_real_nullable REAL,
  a_serial SERIAL NOT NULL,
  a_serial_nullable SERIAL,
  a_smallint SMALLINT NOT NULL,
  a_smallint_nullable SMALLINT,
  a_smallserial SMALLSERIAL NOT NULL,
  a_smallserial_nullable SMALLSERIAL,
  a_text TEXT NOT NULL,
  a_text_nullable TEXT,
  a_time TIME NOT NULL,
  a_time_nullable TIME,
  a_timestamp TIMESTAMP NOT NULL,
  a_timestamp_nullable TIMESTAMP,
  a_timestamptz TIMESTAMPTZ NOT NULL,
  a_timestamptz_nullable TIMESTAMPTZ,
  a_timetz TIMETZ NOT NULL,
  a_timetz_nullable TIMETZ,
  a_uuid UUID NOT NULL,
  a_uuid_nullable UUID,
  a_varchar VARCHAR NOT NULL,
  a_varchar_nullable VARCHAR,
  a_xml XML NOT NULL,
  a_xml_nullable XML
);

-- views
CREATE VIEW a_view_of_everything AS
  SELECT * FROM a_bit_of_everything;

CREATE VIEW a_view_of_everything_some AS
  SELECT a_bool, a_text FROM a_bit_of_everything;

-- procs
CREATE PROCEDURE a_0_in_0_out() AS $$
BEGIN
END;
$$ LANGUAGE plpgsql;

CREATE PROCEDURE a_1_in_0_out(a_param INTEGER) AS $$
BEGIN
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION a_0_in_1_out(OUT a_return INTEGER) AS $$
BEGIN
  a_return := 10;
END;
$$ LANGUAGE plpgsql;

-- funcs
CREATE FUNCTION a_1_in_1_out(a_param INTEGER, OUT a_return INTEGER) AS $$
BEGIN
  a_return := a_param;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION a_2_in_2_out(param_one INTEGER, param_two INTEGER, OUT return_one INTEGER, OUT return_two INTEGER) AS $$
BEGIN
  return_one := param_one;
  return_two := param_two;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION a_func_0_in() RETURNS INTEGER AS $$
BEGIN
  RETURN 10;
END;
$$ LANGUAGE plpgsql;

-- provided to keep schema output in parity with other databases
CREATE FUNCTION a_func_1_in(a_param INTEGER) RETURNS INTEGER AS $$
BEGIN
  RETURN a_param;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION a_func_2_in(param_one INTEGER, param_two INTEGER) RETURNS INTEGER AS $$
BEGIN
  RETURN param_one + param_two;
END;
$$ LANGUAGE plpgsql;

-- overloading (postgres only)

CREATE PROCEDURE a_overloaded(param_one INTEGER) AS $$
BEGIN
END;
$$ LANGUAGE plpgsql;

CREATE PROCEDURE a_overloaded(param_one INTEGER, param_two INTEGER) AS $$
BEGIN
END;
$$ LANGUAGE plpgsql;
