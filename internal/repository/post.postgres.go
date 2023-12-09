package repository

import (
	"GoServer/Entities"
	"context"
	"errors"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/jackc/pgx/v5"
	"strconv"
)

type PostPostgres struct {
	dataBases *DataBases
}

func NewPostPostgres(dataBases *DataBases) *PostPostgres {
	return &PostPostgres{
		dataBases: dataBases,
	}
}

/*region posts*/

// CreatePost(ctx context.Context, id uint, postDTO Entities.CreatePostDTO, surveyDTO Entities.CreateSurveyDTO, date int64) (uint, error
func (repository *PostPostgres) CreatePost(ctx context.Context, id uint, postDTO Entities.CreatePostDTO, surveyDTO Entities.CreateSurveyDTO, date int64) (uint, error) {
	var (
		postId uint
	)
	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `INSERT INTO posts (parent_user_id, data, date, files, have_a_survey) VALUES ($1, $2, $5, $3, $4) RETURNING id`,
		id, postDTO.Data, postDTO.Files, postDTO.HaveASurvey, date)
	err := row.Scan(&postId)
	if err != nil {
		return 0, err
	}

	if postDTO.HaveASurvey {
		_, err = repository.dataBases.Postgres.pool.Exec(ctx, `INSERT INTO surveys (parent_post_id, data, background, is_multivoices) VALUES ($1, $2, $3, $4)`,
			postId, surveyDTO.Data, surveyDTO.Background, surveyDTO.IsMultiVoices)
	}

	return postId, err
}

// GetPostsByUserId  really slow function, because it is postgres
// Uses on Cloud Virtual Machine, where I have no access to other DBMS. (Sorry, I am poor student now and can't pay for normal Cloud Virtual Machine)
func (repository *PostPostgres) GetPostsByUserId(ctx context.Context, authorId uint, offset uint, userId uint) ([]Entities.GetPostDTO, []Entities.GetSurveyDTO, error) {
	rows, err := repository.dataBases.Postgres.pool.Query(ctx, `SELECT id, likes, dislikes, data, date, files, have_a_survey
		FROM posts WHERE parent_user_id = $1 ORDER BY id DESC LIMIT 20 OFFSET $2`, authorId, offset)
	defer rows.Close()
	if err != nil {
		return []Entities.GetPostDTO{}, []Entities.GetSurveyDTO{}, err
	}

	var (
		canVote                   = userId > 0
		clientIdInt               = int32(userId)
		surveys                   = []Entities.GetSurveyDTO{}
		necessarySurveys          []uint
		necessaryLikesAndDislikes        = make(map[uint]struct{}, 20)
		posts                            = []Entities.GetPostDTO{}
		voices                    uint16 = 0
	)

	var (
		id         uint
		data       string
		date       int64
		likes      uint
		dislikes   uint
		files      []string
		haveSurvey bool
	)

	for rows.Next() {
		err = rows.Scan(&id, &likes, &dislikes, &data, &date, &files, &haveSurvey)
		if err != nil {
			repository.dataBases.Postgres.logger.Error("GetPostsByUserId rows.Scan posts, error:", err.Error())
			continue
		}
		if canVote {
			necessaryLikesAndDislikes[id] = struct{}{}
		}

		if haveSurvey {
			necessarySurveys = append(necessarySurveys, id)
		}

		posts = append(posts, Entities.GetPostDTO{
			Id:         id,
			Likes:      likes,
			Dislikes:   dislikes,
			Data:       data,
			Date:       date,
			Files:      files,
			LikeStatus: 0,
			HaveSurvey: haveSurvey,
		})
	}

	// id is 8 byte, "," is 1 byte so (8+1)*20- 1 (-1 because the last item has no ",") = 181
	buf := make([]byte, 0, 179)

	if canVote {
		necessaryLikesAndDislikesLength := len(necessaryLikesAndDislikes)
		if necessaryLikesAndDislikesLength < 1 {
			return posts, []Entities.GetSurveyDTO{}, nil
		}
		i := 0
		for id = range necessaryLikesAndDislikes {
			if i == necessaryLikesAndDislikesLength-1 {
				buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
				i = 0
			} else {
				i += 1
				buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
				buf = append(buf, ',')
			}
		}
		query := fastbytes.B2S(append(append(fastbytes.S2B(`SELECT parent_post_id FROM posts_likes WHERE user_id = $1 AND parent_post_id IN (`), buf...), ')'))
		rows, err = repository.dataBases.Postgres.pool.Query(ctx, query, clientIdInt)
		if err != nil {
			return posts, []Entities.GetSurveyDTO{}, err
		}

		var parentPostId uint

		for rows.Next() {
			err = rows.Scan(&parentPostId)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					continue
				}
				repository.dataBases.Postgres.logger.Error("GetPostsByUserId rows.Scan likes, error:", err.Error())
				continue
			}

			for i = 0; i < len(posts); i++ {
				if posts[i].Id == parentPostId {
					posts[i].LikeStatus = 1
					break
				}
			}
			delete(necessaryLikesAndDislikes, parentPostId)
		}

		buf = buf[:0]

		necessaryLikesAndDislikesLength = len(necessaryLikesAndDislikes)
		i = 0

		if len(necessaryLikesAndDislikes) > 0 {
			for id = range necessaryLikesAndDislikes {
				if i == necessaryLikesAndDislikesLength-1 {
					buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
					i = 0
				} else {
					buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
					buf = append(buf, ',')
					i += 1
				}
			}
			query = fastbytes.B2S(append(append(fastbytes.S2B(`SELECT parent_post_id FROM posts_dislikes WHERE user_id = $1 AND parent_post_id IN (`), buf...), ')'))
			rows, err = repository.dataBases.Postgres.pool.Query(ctx, query, clientIdInt)
			if err != nil {
				return posts, []Entities.GetSurveyDTO{}, err
			}

			for rows.Next() {
				err = rows.Scan(&parentPostId)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						continue
					}
					repository.dataBases.Postgres.logger.Error("GetPostsByUserId rows.Scan dislikes, error:", err.Error())
					continue
				}

				for i = 0; i < len(posts); i++ {
					if posts[i].Id == parentPostId {
						posts[i].LikeStatus = -1
						break
					}
				}
			}
			buf = buf[:0]
		}
	}

	necessarySurveysL := len(necessarySurveys)
	if necessarySurveysL < 1 {
		return posts, []Entities.GetSurveyDTO{}, nil
	}
	for i, v := range necessarySurveys {
		if i == necessarySurveysL-1 {
			buf = append(buf, fastbytes.S2B(strconv.Itoa(int(v)))...)
		} else {
			buf = append(buf, fastbytes.S2B(strconv.Itoa(int(v)))...)
			buf = append(buf, ',')
		}
	}
	query := fastbytes.B2S(append(append(fastbytes.S2B(`SELECT parent_post_id, data, sl0v, sl1v, sl2v, sl3v, sl4v, sl5v, sl6v, sl7v, sl8v, sl9v, background, is_multiVoices 
		FROM surveys WHERE parent_post_id IN (`), buf...), ')'))
	rows, err = repository.dataBases.Postgres.pool.Query(ctx, query)
	if err != nil {
		return posts, []Entities.GetSurveyDTO{}, err
	}

	for rows.Next() {
		var (
			parentPostId  uint
			surveyData    []string
			sl0v          uint
			sl1v          uint
			sl2v          uint
			sl3v          uint
			sl4v          uint
			sl5v          uint
			sl6v          uint
			sl7v          uint
			sl8v          uint
			sl9v          uint
			background    uint8
			isMultiVoices bool
		)
		err = rows.Scan(&parentPostId, &surveyData, &sl0v, &sl1v, &sl2v, &sl3v, &sl4v, &sl5v, &sl6v, &sl7v, &sl8v, &sl9v, &background, &isMultiVoices)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			repository.dataBases.Postgres.logger.Error("GetPostsByUserId rows.Scan surveys, error:", err.Error())
			continue
		}
		surveys = append(surveys, Entities.GetSurveyDTO{
			ParentPostId:  parentPostId,
			Data:          surveyData,
			SL0V:          sl0v,
			SL1V:          sl1v,
			SL2V:          sl2v,
			SL3V:          sl3v,
			SL4V:          sl4v,
			SL5V:          sl5v,
			SL6V:          sl6v,
			SL7V:          sl7v,
			SL8V:          sl8v,
			SL9V:          sl9v,
			Background:    background,
			IsMultiVoices: isMultiVoices,
		})
	}

	if canVote {
		var parentPostId uint
		query = fastbytes.B2S(append(append(fastbytes.S2B(`SELECT voices, parent_survey_id FROM surveys_voices WHERE user_id = $1 AND parent_survey_id IN (`), buf...), ')'))
		rows, err = repository.dataBases.Postgres.pool.Query(ctx, query, clientIdInt)
		if err != nil {
			return posts, surveys, err
		}

		for rows.Next() {
			err = rows.Scan(&voices, &parentPostId)
			if err != nil {
				repository.dataBases.Postgres.logger.Error("GetPostsByUserId rows.Scan voices, error:", err.Error())
				return posts, surveys, err
			}

			for i := 0; i < len(surveys); i++ {
				s := &surveys[i]
				if s.ParentPostId == parentPostId {
					s.VotedFor = voices
					break
				}
			}
		}
	}

	return posts, surveys, nil
}

func (repository *PostPostgres) LikePost(ctx context.Context, userId, postId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("LikePost can't start transaction", err.Error())
		return err
	}

	rows, err := tx.Exec(ctx, `DELETE FROM posts_dislikes WHERE parent_post_id=$1 AND user_id=$2 RETURNING True`, postId, userId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			_ = tx.Rollback(ctx)
			repository.dataBases.Postgres.logger.Error("LikePost can't delete dislike", err.Error())
			return err
		}
	} else {
		if rows.RowsAffected() == 1 {
			_, err = tx.Exec(ctx, `UPDATE posts SET dislikes=dislikes-1 WHERE id=$1`, postId)
			if err != nil {
				_ = tx.Rollback(ctx)
				repository.dataBases.Postgres.logger.Error("LikePost can't update likes", err.Error())
				return err
			}
		}
	}

	_, err = tx.Exec(ctx, `INSERT INTO posts_likes (parent_post_id, user_id) VALUES ($1, $2)`, postId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("LikePost can't insert like", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE posts SET likes = likes+1 WHERE id = $1`, postId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("LikePost can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) UnlikePost(ctx context.Context, userId, postId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("UnlikePost can't start transaction", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM posts_likes WHERE parent_post_id = $1 AND user_id = $2`, postId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UnlikePost can't delete like", err.Error())
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE posts SET likes = likes-1 WHERE id = $1`, postId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UnlikePost can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) DislikePost(ctx context.Context, userId, postId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("DislikePost can't start transaction", err.Error())
		return err
	}

	rows, err := tx.Exec(ctx, `DELETE FROM posts_likes WHERE parent_post_id = $1 AND user_id = $2`, postId, userId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			_ = tx.Rollback(ctx)
			repository.dataBases.Postgres.logger.Error("DislikePost can't delete like", err.Error())
			return err
		}
	} else {
		if rows.RowsAffected() == 1 {
			_, err = tx.Exec(ctx, `UPDATE posts SET likes = likes-1 WHERE id=$1`, postId)
			if err != nil {
				_ = tx.Rollback(ctx)
				repository.dataBases.Postgres.logger.Error("DislikePost can't update likes", err.Error())
				return err
			}
		}
	}

	_, err = tx.Exec(ctx, `INSERT INTO posts_dislikes (parent_post_id, user_id) VALUES ($1, $2)`, postId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("DislikePost can't insert dislike", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE posts SET dislikes = dislikes+1 WHERE id = $1`, postId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("DislikePost can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) UndislikePost(ctx context.Context, userId, postId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("UndislikePost can't start transaction", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM posts_dislikes WHERE parent_post_id = $1 AND user_id = $2`, postId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UndislikePost can't delete dislike", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE posts SET dislikes = dislikes-1 WHERE id = $1`, postId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UndislikePost can't update dislikes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) DeletePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.dataBases.Postgres.pool.Exec(ctx, "DELETE FROM posts WHERE id = $1 AND parent_user_id=$2", postId, userId)
	return err
}

/*endregion*/

/*region comments*/

func (repository *PostPostgres) GetCommentsByPostId(ctx context.Context, postId uint, offset uint, userId uint) ([]Entities.Comment, error) {
	canVote := userId > 0
	necessaryLikesAndDislikes := make(map[uint]struct{}, 20)

	rows, err := repository.dataBases.Postgres.pool.Query(ctx, `
        SELECT id, parent_post_id, data, date, parent_user_id, likes, dislikes, files, parent_comment_id
        FROM comments
        WHERE parent_post_id = $1
        ORDER BY id DESC 
        LIMIT 20 OFFSET $2
    `, postId, offset)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't get comments", err.Error())
		return nil, err
	}

	var comments = []Entities.Comment{}
	for rows.Next() {
		var comment Entities.Comment
		err = rows.Scan(&comment.Id, &comment.ParentPostId, &comment.Data, &comment.Date, &comment.ParentUserId, &comment.Likes, &comment.Dislikes, &comment.Files, &comment.ParentCommentId)
		if err != nil {
			repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't scan comments", err.Error())
			return nil, err
		}
		comment.LikesStatus = 0
		comments = append(comments, comment)
		if canVote {
			necessaryLikesAndDislikes[comment.Id] = struct{}{}
		}
	}

	if canVote {
		necessaryLikesAndDislikesLength := len(necessaryLikesAndDislikes)
		if necessaryLikesAndDislikesLength < 1 {
			return comments, nil
		}

		// 20 * (8 + 1) - 1 = 179
		buf := make([]byte, 0, 179)
		i := 0
		for id := range necessaryLikesAndDislikes {
			if i == necessaryLikesAndDislikesLength-1 {
				buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
				i = 0
			} else {
				i += 1
				buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
				buf = append(buf, ',')
			}
		}

		query := fastbytes.B2S(append(append(fastbytes.S2B(`SELECT parent_comment_id FROM comments_likes
                         WHERE user_id = $1 AND parent_comment_id IN (`), buf...), ')'))
		rows, err = repository.dataBases.Postgres.pool.Query(ctx, query, userId)
		defer rows.Close()

		if err != nil {
			repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't get comments likes", err.Error())
			return comments, err
		}

		for rows.Next() {
			var commentId uint
			err = rows.Scan(&commentId)
			if err != nil {
				repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't scan comments likes", err.Error())
				return comments, err
			}
			for i = 0; i < len(comments); i++ {
				if comments[i].Id == commentId {
					delete(necessaryLikesAndDislikes, commentId)
					comments[i].LikesStatus = 1
				}
			}
		}

		necessaryLikesAndDislikesLength = len(necessaryLikesAndDislikes)

		if necessaryLikesAndDislikesLength > 0 {
			buf = buf[:0]

			i = 0
			for id := range necessaryLikesAndDislikes {
				if i == necessaryLikesAndDislikesLength-1 {
					buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
					i = 0
				} else {
					i += 1
					buf = append(buf, fastbytes.S2B(strconv.Itoa(int(id)))...)
					buf = append(buf, ',')
				}
			}

			query = fastbytes.B2S(append(append(fastbytes.S2B(`SELECT parent_comment_id FROM comments_dislikes WHERE user_id = $1 AND parent_comment_id IN (`), buf...), ')'))
			rows, err = repository.dataBases.Postgres.pool.Query(ctx, query, userId)
			defer rows.Close()
			if err != nil {
				repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't get comments dislikes", err.Error())
				return comments, err
			}

			for rows.Next() {
				var commentId uint
				err = rows.Scan(&commentId)
				if err != nil {
					repository.dataBases.Postgres.logger.Error("GetCommentsByPostId can't scan comments dislikes", err.Error())
					return comments, err
				}
				for i = 0; i < len(comments); i++ {
					if comments[i].Id == commentId {
						comments[i].LikesStatus = -1
					}
				}
			}
		}
	}

	return comments, nil
}

func (repository *PostPostgres) CreateComment(ctx context.Context, userId uint, comment Entities.CommentDTO, date int64) (commentId uint, err error) {
	var row pgx.Row
	if comment.ParentCommentId > 1 {
		row = repository.dataBases.Postgres.pool.QueryRow(ctx, `INSERT INTO comments (parent_post_id, data, date, parent_user_id, files, parent_comment_id)
			VALUES ($1, $2, $6, $3, $4, $5) RETURNING id`,
			comment.ParentPostId, comment.Data, userId, comment.Files, comment.ParentCommentId, date)
	} else {
		row = repository.dataBases.Postgres.pool.QueryRow(ctx, `INSERT INTO comments (parent_post_id, data, date, parent_user_id, files)
			VALUES ($1, $2, $5, $3, $4) RETURNING id`,
			comment.ParentPostId, comment.Data, userId, comment.Files, date)
	}
	err = row.Scan(&commentId)
	if err != nil {
		return 0, err
	}
	return commentId, nil
}

func (repository *PostPostgres) LikeComment(ctx context.Context, userId uint, commentId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("LikePost can't start transaction", err.Error())
		return err
	}

	rows, err := tx.Exec(ctx, `DELETE FROM comments_dislikes WHERE parent_comment_id=$1 AND user_id=$2 RETURNING True`, commentId, userId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			_ = tx.Rollback(ctx)
			repository.dataBases.Postgres.logger.Error("LikeComment can't delete dislike", err.Error())
			return err
		}
	} else {
		if rows.RowsAffected() == 1 {
			_, err = tx.Exec(ctx, `UPDATE comments SET dislikes=dislikes-1 WHERE id=$1`, commentId)
			if err != nil {
				_ = tx.Rollback(ctx)
				repository.dataBases.Postgres.logger.Error("LikePost can't update likes", err.Error())
				return err
			}
		}
	}

	_, err = tx.Exec(ctx, `INSERT INTO comments_likes (parent_comment_id, user_id) VALUES ($1, $2)`, commentId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("LikePost can't insert like", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE comments SET likes = likes+1 WHERE id = $1`, commentId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("LikePost can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) UnlikeComment(ctx context.Context, userId uint, commentId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("UnlikeComment can't start transaction", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM comments_likes WHERE parent_comment_id = $1 AND user_id = $2`, commentId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UnlikeComment can't delete like", err.Error())
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE comments SET likes = likes-1 WHERE id = $1`, commentId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UnlikeComment can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) DislikeComment(ctx context.Context, userId uint, commentId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("DislikePost can't start transaction", err.Error())
		return err
	}

	rows, err := tx.Exec(ctx, `DELETE FROM comments_likes WHERE parent_comment_id=$1 AND user_id=$2 RETURNING True`, commentId, userId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			_ = tx.Rollback(ctx)
			repository.dataBases.Postgres.logger.Error("DislikeComment can't delete dislike", err.Error())
			return err
		}
	} else {
		if rows.RowsAffected() == 1 {
			_, err = tx.Exec(ctx, `UPDATE comments SET likes=likes-1 WHERE id=$1`, commentId)
			if err != nil {
				_ = tx.Rollback(ctx)
				repository.dataBases.Postgres.logger.Error("DislikePost can't update likes", err.Error())
				return err
			}
		}
	}

	_, err = tx.Exec(ctx, `INSERT INTO comments_dislikes (parent_comment_id, user_id) VALUES ($1, $2)`, commentId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("DislikePost can't insert like", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE comments SET dislikes = dislikes+1 WHERE id = $1`, commentId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("DislikePost can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) UndislikeComment(ctx context.Context, userId uint, commentId uint) error {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		repository.dataBases.Postgres.logger.Error("UndislikeComment can't start transaction", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM comments_dislikes WHERE parent_comment_id = $1 AND user_id = $2`, commentId, userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UndislikeComment can't delete like", err.Error())
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE comments SET dislikes = dislikes-1 WHERE id = $1`, commentId)
	if err != nil {
		_ = tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("UndislikeComment can't update likes", err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (repository *PostPostgres) UpdateComment(ctx context.Context, userId uint, commentId uint, updateDTO Entities.CommentUpdateDTO) error {
	_, err := repository.dataBases.Postgres.pool.Exec(ctx, `UPDATE comments 
		SET data = $1, files = $2 
		WHERE id = $3 AND parent_user_id = $4
	`, updateDTO.Data, updateDTO.Files, commentId, userId)
	return err
}

func (repository *PostPostgres) DeleteComment(ctx context.Context, userId uint, commentId uint) error {
	_, err := repository.dataBases.Postgres.pool.Exec(ctx, `DELETE FROM comments 
		WHERE id = $1 AND parent_user_id = $2
	`, commentId, userId)
	return err
}

/*endregion*/

/*region survey */

func (repository *PostPostgres) VoteInSurvey(ctx context.Context, userId uint, surveyId uint, votedFor uint16) error {
	//var votedForForPostgres uint16
	//for _, v := range votedFor {
	//	if v > 9 {
	//		continue
	//	}
	//	// We set a bit of number of votedFor to 1
	//	votedForForPostgres |= 1 << v
	//}

	// "UPDATE surveys SET " is 19 bytes. In worst case we need to use "slnv=slnv+1," 10 times it is 11 * 10 = 110 bytes.
	//"WHERE parent_post_id=$user_id", where user_id is 8 bytes so it is 21 + 8 = 29 bytes. So we need to use 19 + 110 + 29 = 158 bytes.
	query2Buf := make([]byte, 0, 158)
	query2Buf = append(query2Buf, "UPDATE surveys SET "...)
	// TODO check for multivoices
	for i := 0; i < 10; i++ {
		if votedFor&(1<<uint16(i)) > 0 {
			switch i {
			case 0:
				query2Buf = append(query2Buf, fastbytes.S2B("sl0v=sl0v+1,")...)
			case 1:
				query2Buf = append(query2Buf, fastbytes.S2B("sl1v=sl1v+1,")...)
			case 2:
				query2Buf = append(query2Buf, fastbytes.S2B("sl2v=sl2v+1,")...)
			case 3:
				query2Buf = append(query2Buf, fastbytes.S2B("sl3v=sl3v+1,")...)
			case 4:
				query2Buf = append(query2Buf, fastbytes.S2B("sl4v=sl4v+1,")...)
			case 5:
				query2Buf = append(query2Buf, fastbytes.S2B("sl5v=sl5v+1,")...)
			case 6:
				query2Buf = append(query2Buf, fastbytes.S2B("sl6v=sl6v+1,")...)
			case 7:
				query2Buf = append(query2Buf, fastbytes.S2B("sl7v=sl7v+1,")...)
			case 8:
				query2Buf = append(query2Buf, fastbytes.S2B("sl8v=sl8v+1,")...)
			case 9:
				query2Buf = append(query2Buf, fastbytes.S2B("sl9v=sl9v+1,")...)
			}
		}
	}
	query2Buf[len(query2Buf)-1] = ' '

	query2Buf = append(query2Buf, fastbytes.S2B("WHERE parent_post_id=")...)
	query2Buf = append(query2Buf, fastbytes.S2B(strconv.Itoa(int(surveyId)))...)

	tx, err := repository.dataBases.Postgres.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		repository.dataBases.Postgres.logger.Error("vote in survey error (failed to begin transaction): ", err.Error())
		return err
	}
	query := `INSERT INTO surveys_voices (parent_survey_id, user_id, voices) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, query, surveyId, userId, votedFor)
	if err != nil {
		tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("vote in survey error (insert voices): ", err.Error())
		return err
	}

	_, err = tx.Exec(ctx, fastbytes.B2S(query2Buf))
	if err != nil {
		tx.Rollback(ctx)
		repository.dataBases.Postgres.logger.Error("vote in survey error (update surveys): ", err.Error())
		return err
	}

	tx.Commit(ctx)
	return err
}

/*endregion*/
