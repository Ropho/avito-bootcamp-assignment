/* /////// */
/* USERS */
/* /////// */
CREATE TABLE IF NOT EXISTS users (
    "uuid"      uuid        PRIMARY KEY,
    email       text        UNIQUE NOT NULL,
    encr_pass   text        NOT NULL,
    salt        text        UNIQUE NOT NULL,
    user_type   text        NOT NULL
);



/* /////// */
/* HOUSES */
/* /////// */
CREATE TABLE IF NOT EXISTS houses (
    house_id    serial  PRIMARY KEY,
    address     text    UNIQUE NOT NULL,
    year        integer NOT NULL,
    developer   text,
    created_at  timestamp NOT NULL,
    updated_at  timestamp NOT NULL
);


/* /////// */
/* FLATS */
/* /////// */
CREATE TABLE IF NOT EXISTS flats (
    flat_id         serial  NOT NULL,
    house_id        integer NOT NULL,
    price           integer NOT NULL,
    rooms_number    integer NOT NULL,
    status          text    NOT NULL,

    PRIMARY KEY (flat_id, house_id),
    FOREIGN KEY (house_id) REFERENCES houses (house_id) ON DELETE CASCADE
);