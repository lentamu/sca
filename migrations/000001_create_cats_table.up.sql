CREATE TABLE IF NOT EXISTS cats
(
    id                  CHAR(36)      NOT NULL,
    name                VARCHAR(255)  NOT NULL UNIQUE,
    years_of_experience INT           NOT NULL,
    breed               VARCHAR(100)  NOT NULL,
    salary              DECIMAL(8, 2) NOT NULL,
    PRIMARY KEY (id)
);