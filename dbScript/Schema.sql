CREATE TABLE users (
                       id BINARY(16) PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       image VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_username_index ON users (username);
CREATE INDEX users_email_index ON users (email);

CREATE TABLE novel_status (
                              id int PRIMARY KEY AUTO_INCREMENT,
                              name VARCHAR(20) NOT NULL
);

INSERT INTO novel_status (name) VALUES ('Ongoing');
INSERT INTO novel_status (name) VALUES ('Completed');
INSERT INTO novel_status (name) VALUES ('Dropped');

CREATE TABLE novels (
                        id BINARY(16) PRIMARY KEY,
                        title VARCHAR(255) NOT NULL,
                        tagline VARCHAR(255) NOT NULL,
                        description TEXT NOT NULL,
                        author VARCHAR(255) NOT NULL,
                        image VARCHAR(255) NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        total_rating INT NOT NULL DEFAULT 0,
                        rate_count INT NOT NULL DEFAULT 0,
                        views INT NOT NULL DEFAULT 0,
                        clicks INT NOT NULL DEFAULT 0,
                        adult BOOLEAN NOT NULL DEFAULT FALSE,
                        status_id INT NOT NULL DEFAULT 1
);

CREATE INDEX novels_title_index ON novels (title);
CREATE INDEX novels_author_index ON novels (author);
CREATE INDEX novels_status_id_index ON novels (status_id);

CREATE TABLE tags (
                      id BINARY(16) PRIMARY KEY,
                      name VARCHAR(255) NOT NULL,
                      description TEXT NOT NULL,
                      created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX tags_name_index ON tags (name);

CREATE TABLE novel_tags (
                            novel_id BINARY(16) NOT NULL,
                            tag_id BINARY(16) NOT NULL,
                            PRIMARY KEY (novel_id, tag_id)
);

CREATE INDEX novel_tags_novel_id_index ON novel_tags (novel_id);
CREATE INDEX novel_tags_tag_id_index ON novel_tags (tag_id);

CREATE TABLE volumes (
                         id BINARY(16) PRIMARY KEY,
                         novel_id BINARY(16) NOT NULL,
                         title VARCHAR(255) NOT NULL,
                         tagline VARCHAR(255) NOT NULL,
                         description TEXT NOT NULL,
                         image VARCHAR(255) NOT NULL,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         views INT NOT NULL DEFAULT 0,
                         clicks INT NOT NULL DEFAULT 0
);

CREATE INDEX volumes_novel_id_index ON volumes (novel_id);
CREATE INDEX volumes_title_index ON volumes (title);

CREATE TABLE chapters (
                          id BINARY(16) PRIMARY KEY,
                          volume_id BINARY(16) NOT NULL,
                          title VARCHAR(255) NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          views INT NOT NULL DEFAULT 0,
                          clicks INT NOT NULL DEFAULT 0
);

CREATE INDEX chapters_volume_id_index ON chapters (volume_id);
CREATE INDEX chapters_title_index ON chapters (title);

CREATE TABLE comments (
                          id BINARY(16) PRIMARY KEY,
                          to_id BINARY(16) NOT NULL,
                          user_id BINARY(16) NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX comments_to_id_index ON comments (to_id);
CREATE INDEX comments_user_id_index ON comments (user_id);

CREATE TABLE images (
                        id BINARY(16) PRIMARY KEY,
                        user_id BINARY(16) NOT NULL,
                        novel_id BINARY(16) NOT NULL,
                        url VARCHAR(255) NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX images_user_id_index ON images (user_id);
CREATE INDEX images_novel_id_index ON images (novel_id);

CREATE TABLE report_reason (
                               id int PRIMARY KEY AUTO_INCREMENT,
                               reason VARCHAR(50) NOT NULL
);

INSERT INTO report_reason (reason) VALUES ('Spam');
INSERT INTO report_reason (reason) VALUES ('Inappropriate');
INSERT INTO report_reason (reason) VALUES ('Other');


CREATE TABLE reported (
                          id BINARY(16) PRIMARY KEY,
                          user_id BINARY(16) NOT NULL,
                          novel_id BINARY(16) NOT NULL,
                          reason_id INT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX reported_user_id_index ON reported (user_id);
CREATE INDEX reported_novel_id_index ON reported (novel_id);
