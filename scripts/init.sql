CREATE TABLE IF NOT EXISTS sites (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    interval INTEGER NOT NULL
);

INSERT INTO sites (name, url, interval) 
VALUES ('Onet', 'https://onet.pl', 60);