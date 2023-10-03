CREATE TABLE users
(
    id          BINARY(16) PRIMARY KEY,
    username    CHAR(32) UNIQUE NOT NULL,
    displayname VARCHAR(64),
    password    BINARY(60)      NOT NULL,
    email       VARCHAR(255) UNIQUE      DEFAULT NULL,
    image       VARCHAR(255)    NOT NULL DEFAULT '',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_username_index ON users (username);
CREATE INDEX users_email_index ON users (email);

CREATE TABLE novel_status
(
    id   INT PRIMARY KEY AUTO_INCREMENT,
    name CHAR(9) NOT NULL
);

INSERT INTO novel_status (name)
VALUES ('Ongoing');
INSERT INTO novel_status (name)
VALUES ('Completed');
INSERT INTO novel_status (name)
VALUES ('Dropped');

CREATE TABLE visibility
(
    id   INT PRIMARY KEY AUTO_INCREMENT,
    name CHAR(3) UNIQUE NOT NULL
);

INSERT INTO visibility (name)
VALUES ('PRI');
INSERT INTO visibility (name)
VALUES ('PUB');

CREATE TABLE novels
(
    id           BINARY(16) PRIMARY KEY,
    title        VARCHAR(255)  NOT NULL,
    tagline      VARCHAR(255)  NOT NULL,
    description  VARCHAR(5000) NOT NULL,
    author       BINARY(16)    NOT NULL,
    image        VARCHAR(255)  NOT NULL,
    language     CHAR(3)       NOT NULL,
    created_at   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    total_rating INT           NOT NULL DEFAULT 0,
    rate_count   INT           NOT NULL DEFAULT 0,
    views        INT           NOT NULL DEFAULT 0,
    clicks       INT           NOT NULL DEFAULT 0,
    adult        BOOLEAN       NOT NULL DEFAULT FALSE,
    status_id    INT           NOT NULL DEFAULT 1,
    visibility   INT           NOT NULL DEFAULT 1
);

CREATE INDEX novels_title_index ON novels (title);
CREATE INDEX novels_author_index ON novels (author);
CREATE INDEX novels_status_id_index ON novels (status);

CREATE TABLE tags
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(50)  NOT NULL,
    description VARCHAR(300) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX tags_name_index ON tags (name);

CREATE TABLE novel_tags
(
    novel_id BINARY(16) NOT NULL,
    tag_id   INT        NOT NULL,
    PRIMARY KEY (novel_id, tag_id)
);

CREATE INDEX novel_tags_novel_id_index ON novel_tags (novel_id);
CREATE INDEX novel_tags_tag_id_index ON novel_tags (tag_id);

CREATE TABLE volumes
(
    id          BINARY(16) PRIMARY KEY,
    novel_id    BINARY(16)    NOT NULL,
    title       VARCHAR(255)  NOT NULL,
    tagline     VARCHAR(255)  NOT NULL,
    description VARCHAR(5000) NOT NULL,
    image       VARCHAR(255)  NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    views       INT           NOT NULL DEFAULT 0,
    visibility  INT           NOT NULL DEFAULT 1
);

CREATE INDEX volumes_novel_id_index ON volumes (novel_id);
CREATE INDEX volumes_title_index ON volumes (title);

CREATE TABLE chapters
(
    id         BINARY(16) PRIMARY KEY,
    volume_id  BINARY(16)   NOT NULL,
    title      VARCHAR(255) NOT NULL,
    content    MEDIUMTEXT   NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    views      INT          NOT NULL DEFAULT 0,
    visibility INT          NOT NULL DEFAULT 1
);

CREATE INDEX chapters_volume_id_index ON chapters (volume_id);
CREATE INDEX chapters_title_index ON chapters (title);

CREATE TABLE comments
(
    id         BINARY(16) PRIMARY KEY,
    to_id      BINARY(16)    NOT NULL,
    user_id    BINARY(16)    NOT NULL,
    content    VARCHAR(5000) NOT NULL,
    created_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX comments_to_id_index ON comments (to_id);
CREATE INDEX comments_user_id_index ON comments (user_id);

CREATE TABLE images
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    user_id    BINARY(16)   NOT NULL,
    novel_id   BINARY(16)   NOT NULL,
    url        VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX images_user_id_index ON images (user_id);
CREATE INDEX images_novel_id_index ON images (novel_id);

CREATE TABLE follows_user
(
    from_id BINARY(16) NOT NULL,
    to_id   BINARY(16) NOT NULL,
    PRIMARY KEY (from_id, to_id)
);

CREATE TABLE follows_novel
(
    user_id BINARY(16) NOT NULL,
    novel_id BINARY(16) NOT NULL,
    PRIMARY KEY (user_id, novel_id)
);

CREATE TABLE report_reason
(
    id     int PRIMARY KEY AUTO_INCREMENT,
    reason VARCHAR(50) NOT NULL
);

INSERT INTO report_reason (reason)
VALUES ('Spam');
INSERT INTO report_reason (reason)
VALUES ('Inappropriate');
INSERT INTO report_reason (reason)
VALUES ('Other');


CREATE TABLE reports
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    user_id    BINARY(16) NOT NULL,
    to_id      BINARY(16) NOT NULL,
    reason_id  INT        NOT NULL,
    created_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX reported_user_id_index ON reports (user_id);
CREATE INDEX reported_novel_id_index ON reports (to_id);

CREATE TABLE sessions
(
    id          BINARY(16) PRIMARY KEY,
    user_id     BINARY(16)   NOT NULL,
    expires_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device_name VARCHAR(255) NOT NULL
);

CREATE INDEX sessions_user_id_index ON sessions (user_id);