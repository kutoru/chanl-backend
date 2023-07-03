# 1
INSERT INTO users (name, password, created_at)
values ('Kut', '1234', now());

# 2
INSERT INTO users (name, password, created_at)
values ('Oru', '1234', now());



# 1
INSERT INTO channels (owner_id, name, type, created_at)
value (1, 'Global', 'gl', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 1, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 1, 0, now());



# 2
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 1, "Kut's private channel", 'pr', now());

# 3
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 1, "Kut's personal channel", 'pe', now());

# 4
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 1, "Oru's private channel", 'pr', now());

# 5
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 1, "Oru's personal channel", 'pe', now());

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
value (1, 2, "Kut's server", 'se', now());

# 7
INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 4, "Oru's server", 'se', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 6, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 7, 0, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 7, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 6, 0, now());
