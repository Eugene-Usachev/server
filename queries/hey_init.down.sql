DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS posts;
DROP INDEX IF EXISTS parent_user_idx;

DROP TABLE IF EXISTS posts_likes;
DROP INDEX IF EXISTS posts_likes_user_idx;
DROP INDEX IF EXISTS posts_likes_parent_post_idx;

DROP TABLE IF EXISTS posts_dislikes;
DROP INDEX IF EXISTS posts_dislikes_user_idx;
DROP INDEX IF EXISTS posts_dislikes_parent_post_idx;

DROP TABLE IF EXISTS surveys;
DROP INDEX IF EXISTS parent_posts_idx;

DROP TABLE IF EXISTS surveys_voices;
DROP INDEX IF EXISTS surveys_voices_user_idx;

DROP TABLE IF EXISTS comments;
DROP INDEX IF EXISTS parent_post_idx;

DROP TABLE IF EXISTS chats;

DROP TABLE IF EXISTS messages;
DROP INDEX IF EXISTS parent_chat_idx;

DROP TABLE IF EXISTS musics;
