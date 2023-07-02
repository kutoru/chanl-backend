INSERT INTO users (name, password, created_at)
values ('Kut', '1234', now());

INSERT INTO users (name, password, created_at)
values ('Oru', '1234', now());

INSERT INTO channels (owner_id, name, type, created_at)
value (1, 'Global', 'gl', now());

INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 1, "Kut's private channel", 'pr', now());

INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 1, "Oru's private channel", 'pr', now());

INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (1, 2, "Kut's server", 'se', now());

INSERT INTO channels (owner_id, parent_id, name, type, created_at)
value (2, 3, "Oru's server", 'se', now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 1, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 2, 1, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 4, 0, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (1, 5, 0, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 1, 0, now());

INSERT INTO joined_channels (user_id, channel_id, can_write, joined_at)
value (2, 3, 1, now());
