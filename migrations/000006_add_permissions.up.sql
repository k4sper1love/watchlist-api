CREATE TABLE IF NOT EXISTS permissions
(
    id   BIGSERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_permissions
(
    user_id        BIGINT NOT NULL,
    permissions_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, permissions_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (permissions_id) REFERENCES permissions (id) ON DELETE CASCADE
);

-- Insert default permissions into the permissions table
INSERT INTO permissions (code)
VALUES ('film:create'),
       ('collection:create');