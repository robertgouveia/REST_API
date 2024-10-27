ALTER TABLE IF EXISTS users
ADD COLUMN role_id INT REFERENCES roles (id) DEFAULT 1;
-- trying to resolve the problem with having users already in the database
UPDATE users
SET
    role_id = (
        SELECT id
        FROM roles
        WHERE
            name = 'user'
    );
-- we know we are definitely setting them to user
ALTER TABLE users ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;