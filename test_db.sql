# the passwords are 1234

# 1
INSERT INTO users (name, password, created_at)
values ('Kuto', '$2a$10$n1OXWkpa9B/l7w0HNF15DuK9y2PBO5Dv0r/rqzsOaQ80kdNCXQMCW', now());

# 2
INSERT INTO users (name, password, created_at)
values ('Toru', '$2a$10$rQRIJmr9MWzKK0wnzzbYjuM13zLl.5kc4whSQvDywvktMeXivuU.m', now());



# 1
INSERT INTO channels (owner_id, name, type, created_at)
value (1, 'Global', 'gl', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 1, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 1, 1, now());



# 2
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 1, "Kuto's private channel", 'pr', now());

# 3
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 1, "Kuto's personal channel", 'pe', now());

# 4
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 1, "Toru's private channel", 'pr', now());

# 5
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 1, "Toru's personal channel", 'pe', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 2, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 3, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 4, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 5, 1, now());



# 6
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 2, "Kuto's server", 'se', now());

# 7
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 4, "Toru's server", 'se', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 6, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 7, 0, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 7, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 6, 0, now());
