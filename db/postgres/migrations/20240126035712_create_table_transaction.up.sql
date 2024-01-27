CREATE TABLE transaction (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "to" text NOT NULL,
    "from" text NOT NULL,
    amount NUMERIC(14, 2) NOT NULL,
    date timestamp NOT NULL,
    FOREIGN KEY ("to") REFERENCES wallet (id)
        ON UPDATE CASCADE
        ON DELETE NO ACTION,
    FOREIGN KEY ("from") REFERENCES wallet (id)
        ON UPDATE CASCADE
        ON DELETE NO ACTION
);
