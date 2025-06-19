CREATE TABLE IF NOT EXISTS missions
(
    id       CHAR(36) NOT NULL,
    complete BOOLEAN  NOT NULL DEFAULT false,
    cat_id   CHAR(36) NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (cat_id) REFERENCES cats (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS targets
(
    id         CHAR(36)     NOT NULL,
    name       VARCHAR(100) NOT NULL,
    country    VARCHAR(50)  NOT NULL,
    notes      TEXT         NOT NULL,
    complete   BOOLEAN      NOT NULL DEFAULT false,
    mission_id CHAR(36)     NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (mission_id) REFERENCES missions (id) ON DELETE SET NULL
);