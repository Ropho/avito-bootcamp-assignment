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

CREATE TABLE IF NOT EXISTS subscriptions (
    user_id         uuid     NOT NULL,
    house_id        integer  NOT NULL,
    email           text     NOT NULL,
    PRIMARY KEY (user_id, house_id),
    FOREIGN KEY (user_id)   REFERENCES users  ("uuid")   ON DELETE CASCADE,
    FOREIGN KEY (house_id)  REFERENCES houses (house_id) ON DELETE CASCADE
);

-- Тестовые данные для интеграционных тестов

INSERT INTO users ("uuid", email, encr_pass, salt, user_type) VALUES 
('ab91568c-6a95-49b1-9029-718df981e2df', 'moderator@gmail.com', '3527948363b878f415880d5a3df65ca1c757a83b9848c0e9036e177b2ad05b95bc444538e2e3c1655f4d1b4647f7689bfa6a
ecf68b139c6dccf21b08e7565a1', '76004dc8aa6a787dba09073532f222a4fb7d313217c173f5fb616e410561ef2f', 'moderator');


INSERT INTO houses (address, year, developer, created_at, updated_at) VALUES 
('Лесная улица, 7, Москва, 125196', 2000, 'Мэрия города', '2024-08-08 03:02:35.26279', '2024-08-08 03:02:35.26279'), -- id = 1 (used for create flat tests)
('Долгопрудный, МФТИ', 2000, 'Мэрия города', '2024-08-08 03:02:35.26279', '2024-08-08 03:02:35.26279'); -- id = 2 (used for house/2 tests)


INSERT INTO flats (house_id, price, rooms_number, status) VALUES
(1, 1, 1, 'approved'),
(1, 2, 2, 'on moderation'),
(1, 3, 3, 'created'),
(1, 4, 4, 'declined'),
(2, 1, 1, 'approved'),
(2, 2, 2, 'on moderation'),
(2, 3, 3, 'created'),
(2, 4, 4, 'declined');