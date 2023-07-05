DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS joined_channels;
DROP TABLE IF EXISTS channels;
DROP TABLE IF EXISTS users;

# gl: global
# pr: private
# se: server
# ro: room
# pe: personal
# fr: friend

CREATE TABLE users (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE channels (
    id INT AUTO_INCREMENT NOT NULL,
    owner_id INT,
    parent_id INT,
    name VARCHAR(255) NOT NULL,
    type ENUM('gl', 'pr', 'se', 'ro', 'pe', 'fr') NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES channels(id),
    FOREIGN KEY (owner_id) REFERENCES users(id),
    PRIMARY KEY (id),
    CONSTRAINT chk_parent_id CHECK((id > 1) AND (parent_id IS NOT NULL)),
    CONSTRAINT chk_global CHECK((id > 1) AND (type != 'gl'))
);

CREATE TABLE joined_channels (
    user_id INT NOT NULL,
    channel_id INT NOT NULL,
    can_write BOOL NOT NULL,
    joined_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    PRIMARY KEY (user_id, channel_id)
);

CREATE TABLE messages (
    id INT AUTO_INCREMENT NOT NULL,
    user_id INT NOT NULL,
    channel_id INT NOT NULL,
    text VARCHAR(1024) NOT NULL,
    sent_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    PRIMARY KEY (id)
);
