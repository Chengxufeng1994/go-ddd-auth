INSERT INTO roles(
	id, name, created_at, updated_at)
	VALUES (1, 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO roles(
	id, name, created_at, updated_at)
	VALUES (2, 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO users(
	username, password, role_id, created_at, updated_at)
	VALUES ('admin', '$2a$10$wNuLYXiEsz3Q7QW21JCsV.K6Vnky3BforO6skXeuw6hNwg2el5qSa', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO users(
	username, password, role_id, created_at, updated_at)
	VALUES ('user', '$2a$10$6fED5RKnTVmHnK0gEmLyEew97HEbH8ot68PHRwvTw1zgsUa9vkREG', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);