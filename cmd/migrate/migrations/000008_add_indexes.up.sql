CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE EXTENSION IF NOT EXISTS btree_gin;
-- comments index
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);
-- post title index
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
-- post tags index
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);
-- user index
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING gin (username);
-- post index based on user id
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts USING gin (user_id);
-- comments index based on post id
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments USING gin (post_id);