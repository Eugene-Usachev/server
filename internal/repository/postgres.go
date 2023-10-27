package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	UserPass string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(ctx context.Context, maxAttempts uint8, cfg Config, logger *PostgresLogger) (pool *pgxpool.Pool, err error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.UserName, cfg.UserPass, cfg.Host, cfg.Port, cfg.DBName)
	err = doWithTries(func() error {
		ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx1, url)
		config, _ := pgxpool.ParseConfig(url)
		config.ConnConfig.Tracer = logger
		pool, err = pgxpool.NewWithConfig(ctx1, config)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 15*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	logger.logger.Info("creating tables for postgres")
	err = createTablesAndIndexes(ctx, pool)
	if err != nil {
		logger.logger.Fatal(err)
	} else {
		logger.logger.Info("all tables for postgres created")
	}

	return pool, nil
}

func createTablesAndIndexes(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id                    serial         NOT NULL PRIMARY KEY,
			login                 varchar(32)    UNIQUE NOT NULL,
			password              varchar(64),
			email                 varchar(32)    UNIQUE NOT NULL,
			name                  varchar(16)    NOT NULL,
			surname               varchar(32)    NOT NULL,
			friends               integer []     default ARRAY[0]::int[],
			subscribers           integer []     default ARRAY[0]::int[],
			avatar                varchar(64)    default '',
			birthday              varchar(32)    default '',
			favourites_books      text           default '',
			favourites_films      text           default '',
			favourites_games      text           default '',
			favourites_meals      text           default '',
			description           varchar(256)   default '',
			family_status         SMALLINT       default -1,
			place_of_residence    varchar(64)    default '',
			attitude_to_smocking  SMALLINT       default -1,
			attitude_to_sport     SMALLINT       default -1,
			attitude_to_alcohol   SMALLINT       default -1,
			dreams                varchar(256)   default '',
			chat_lists            text           default '[{"name":"favourites","chats":[]},{"name":"friends","chats":[]},{"name":"subscribers","chats":[]},{"name":"nobody","chats":[]}]',
			all_chats             integer []     default ARRAY[]::integer[] []
		);
		
		CREATE TABLE IF NOT EXISTS posts (
			id             serial        NOT NULL PRIMARY KEY,
			parent_user_id int           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			likes          int           default 0,
			dislikes       int           default 0,
			data           varchar(512),
			date           int8,
			files          text []       default ARRAY[]::text[] [],
			have_a_survey  bool
		);
		CREATE INDEX IF NOT EXISTS parent_user_idx ON posts(parent_user_id);
		
		CREATE TABLE IF NOT EXISTS posts_likes (
			parent_post_id  int         NOT NULL REFERENCES posts(id) ON DELETE CASCADE PRIMARY KEY,
			user_id         int         NOT NULL REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS posts_likes_user_idx ON posts_likes(user_id);
		CREATE INDEX IF NOT EXISTS posts_likes_parent_post_idx ON posts_likes(parent_post_id);
		
		CREATE TABLE IF NOT EXISTS posts_dislikes (
			parent_post_id  int         NOT NULL REFERENCES posts(id) ON DELETE CASCADE PRIMARY KEY,
			user_id         int         NOT NULL REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS posts_dislikes_user_idx ON posts_dislikes(user_id);
		CREATE INDEX IF NOT EXISTS posts_dislikes_parent_post_idx ON posts_dislikes(parent_post_id);
		
		CREATE TABLE IF NOT EXISTS surveys (
			parent_post_id  int         NOT NULL REFERENCES posts(id) ON DELETE CASCADE PRIMARY KEY,
			data            text []     default ARRAY[]::text[] [],
			sl0v            int         default 0,
			sl1v            int         default 0,
			sl2v            int         default 0,
			sl3v            int         default 0,
			sl4v            int         default 0,
			sl5v            int         default 0,
			sl6v            int         default 0,
			sl7v            int         default 0,
			sl8v            int         default 0,
			sl9v            int         default 0,
			background      int 		default 0,
			is_multiVoices  bool        default false
		);
		CREATE INDEX IF NOT EXISTS parent_posts_idx ON surveys(parent_post_id);
		
		CREATE TABLE IF NOT EXISTS surveys_voices (
			parent_survey_id             int8        NOT NULL PRIMARY KEY REFERENCES surveys(parent_post_id) ON DELETE CASCADE,
			user_id                      int         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			voices                       int2        default 0
		);
		CREATE INDEX IF NOT EXISTS surveys_voices_user_idx ON surveys_voices(user_id);
		
		CREATE TABLE IF NOT EXISTS comments (
			id                  serial        NOT NULL PRIMARY KEY,
			parent_post_id      int           NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			data                varchar(128),
			date                int8,
			parent_user_id      int           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			likes               int           default 0,
			likes_by            integer []    default ARRAY[]::integer[] [],
			dislikes            int           default 0,
			dislikes_by         integer []    default ARRAY[]::integer[] [],
			files               text []       default ARRAY[]::text[] [],
			parent_comment_id   int
		);
		CREATE INDEX IF NOT EXISTS parent_post_idx ON comments(parent_post_id);
		
		CREATE TABLE IF NOT EXISTS chats (
			id      serial      NOT NULL PRIMARY KEY,
			name    varchar(64) NOT NULL,
			avatar  varchar(64) default '',
			members integer []  default ARRAY[]::integer[] []
		);
		
		CREATE TABLE IF NOT EXISTS messages (
			id                serial        NOT NULL PRIMARY KEY,
			parent_chat_id    int           NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
			parent_user_id	  int           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			data              varchar(256),
			date              int8			default 0,
			files             text []       default ARRAY[]::text[] [],
			message_parent_id int
		);
		CREATE INDEX IF NOT EXISTS parent_chat_idx ON messages(parent_chat_id);
		
		CREATE TABLE IF NOT EXISTS musics (
			id                      serial       NOT NULL PRIMARY KEY,
			parent_user_id          int          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			author                  varchar(32),
			title                   varchar(64),
			number_of_eavesdroppers int          default 0
		);
	`)
	return err
}

func deleteAllTablesAndIndexes(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
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
	`)
	if err != nil {
		log.Println(err)
	}
}

func regenerateIndexes(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
		DROP INDEX IF EXISTS parent_user_idx;
		CREATE INDEX IF NOT EXISTS parent_user_idx ON posts(parent_user_id);
		
		DROP INDEX IF EXISTS posts_likes_user_idx;
		DROP INDEX IF EXISTS posts_likes_parent_post_idx;
		CREATE INDEX IF NOT EXISTS posts_likes_user_idx ON posts_likes(user_id);
		CREATE INDEX IF NOT EXISTS posts_likes_parent_post_idx ON posts_likes(parent_post_id);
		
		DROP INDEX IF EXISTS posts_dislikes_user_idx;
		DROP INDEX IF EXISTS posts_dislikes_parent_post_idx;
		CREATE INDEX IF NOT EXISTS posts_dislikes_user_idx ON posts_dislikes(user_id);
		CREATE INDEX IF NOT EXISTS posts_dislikes_parent_post_idx ON posts_dislikes(parent_post_id);
		
		DROP INDEX IF EXISTS parent_posts_idx;
		CREATE INDEX IF NOT EXISTS parent_posts_idx ON surveys(parent_post_id);
		
		DROP INDEX IF EXISTS surveys_voices_user_idx;
		CREATE INDEX IF NOT EXISTS surveys_voices_user_idx ON surveys_voices(user_id);
		
		DROP INDEX IF EXISTS parent_post_idx;
		CREATE INDEX IF NOT EXISTS parent_post_idx ON comments(parent_post_id);
		
		DROP INDEX IF EXISTS parent_chat_idx;
		CREATE INDEX IF NOT EXISTS parent_chat_idx ON messages(parent_chat_id);
	`)
	if err != nil {
		log.Println(err)
	}
}

func doWithTries(fn func() error, attempts uint8, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return err
}
