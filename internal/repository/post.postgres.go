package repository

import (
	"GoServer/Entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
)

type PostPostgres struct {
	database *pgxpool.Pool
}

func NewPostPostgres(db *pgxpool.Pool) *PostPostgres {
	return &PostPostgres{
		database: db,
	}
}

/*region posts*/

func (repository *PostPostgres) CreateAPost(ctx context.Context, id uint, postDTO Entities.CreateAPostDTO, surveyDTO Entities.CreateASurveyDTO, date string) error {
	var (
		postId uint
	)
	row := repository.database.QueryRow(ctx, `INSERT INTO posts (parent_user_id, data, date, files, have_a_survey) VALUES ($1, $2, $5, $3, $4) RETURNING id`,
		id, postDTO.Data, postDTO.Files, postDTO.HaveASurvey, date)
	err := row.Scan(&postId)
	if err != nil {
		return err
	}

	if postDTO.HaveASurvey {
		_, err = repository.database.Exec(ctx, `INSERT INTO surveys (parent_post_id, data, background, is_multivoices) VALUES ($1, $2, $3, $4)`,
			postId, surveyDTO.Data, surveyDTO.Background, surveyDTO.IsMultiVoices)
	}

	return err
}

func (repository *PostPostgres) GetPostsByUserID(ctx context.Context, userID uint, offset uint) ([]Entities.Post, []Entities.Survey, error) {

	rows, err := repository.database.Query(ctx, `SELECT id, likes, liked_by, dislikes, disliked_by, data, date, files, have_a_survey
		FROM posts WHERE parent_user_id = $1 ORDER BY id DESC LIMIT 20 OFFSET $2`, userID, offset)
	if err != nil {
		return []Entities.Post{}, []Entities.Survey{}, err
	}

	var (
		surveys          []Entities.Survey
		necessarySurveys []uint
		posts            []Entities.Post
	)
	for rows.Next() {
		var (
			id          uint
			data        string
			date        string
			likes       uint
			likedBy     []int32
			dislikes    uint
			dislikedBy  []int32
			files       []string
			haveASurvey bool
		)
		err = rows.Scan(&id, &likes, &likedBy, &dislikes, &dislikedBy, &data, &date, &files, &haveASurvey)
		if err != nil {
			log.Println(err)
			continue
		}

		if haveASurvey {
			necessarySurveys = append(necessarySurveys, id)
		}
		posts = append(posts, Entities.Post{
			ID:          id,
			Likes:       likes,
			LikedBy:     likedBy,
			Dislikes:    dislikes,
			DislikedBy:  dislikedBy,
			Data:        data,
			Date:        date,
			Files:       files,
			HaveASurvey: haveASurvey,
		})
	}

	if len(necessarySurveys) > 0 {
		array := "("
		for _, v := range necessarySurveys {
			array = fmt.Sprintf("%s%d,", array, v)
		}
		array = array[:len(array)-1] + ")"
		query := `SELECT parent_post_id, data, sl0v, sl1v, sl2v, sl3v, sl4v, sl5v, sl6v, sl7v, sl8v, sl9v, sl0vby, sl1vby, sl2vby, sl3vby, sl4vby, sl5vby, sl6vby, sl7vby, sl8vby, sl9vby, voted_by, background, is_multivoices
			FROM surveys WHERE parent_post_id IN ` + array
		rows, err = repository.database.Query(ctx, query)
		if err != nil {
			return posts, []Entities.Survey{}, err
		}
		for rows.Next() {
			var (
				ParentPostID  int
				Data          []string
				SL0V          int
				SL0VBY        []int32
				SL1V          int
				SL1VBY        []int32
				SL2V          int
				SL2VBY        []int32
				SL3V          int
				SL3VBY        []int32
				SL4V          int
				SL4VBY        []int32
				SL5V          int
				SL5VBY        []int32
				SL6V          int
				SL6VBY        []int32
				SL7V          int
				SL7VBY        []int32
				SL8V          int
				SL8VBY        []int32
				SL9V          int
				SL9VBY        []int32
				VotedBy       []int32
				Background    string
				IsMultiVoices bool
			)
			err = rows.Scan(&ParentPostID, &Data, &SL0V, &SL1V, &SL2V, &SL3V, &SL4V, &SL5V, &SL6V, &SL7V, &SL8V, &SL9V,
				&SL0VBY, &SL1VBY, &SL2VBY, &SL3VBY, &SL4VBY, &SL5VBY, &SL6VBY, &SL7VBY, &SL8VBY, &SL9VBY,
				&VotedBy, &Background, &IsMultiVoices)
			if err != nil {
				log.Println(err)
				continue
			}

			surveys = append(surveys, Entities.Survey{
				ParentPostID: ParentPostID, Data: Data, SL0V: SL0V, SL1V: SL1V, SL2V: SL2V,
				SL3V: SL3V, SL4V: SL4V, SL5V: SL5V, SL6V: SL6V, IsMultiVoices: IsMultiVoices,
				SL7V: SL7V, SL8V: SL8V, SL9V: SL9V, VotedBy: VotedBy, Background: Background,
				SL0VBY: SL0VBY, SL1VBY: SL1VBY, SL2VBY: SL2VBY, SL3VBY: SL3VBY, SL4VBY: SL4VBY,
				SL5VBY: SL5VBY, SL6VBY: SL6VBY, SL7VBY: SL7VBY, SL8VBY: SL8VBY, SL9VBY: SL9VBY,
			})
		}
	}

	err = rows.Err()
	return posts, surveys, err
}

func (repository *PostPostgres) LikePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE posts 
		SET likes=likes+1, liked_by=array_append(liked_by, $1)
		WHERE id=$2 AND NOT $1=ANY(liked_by) AND NOT $1=ANY(disliked_by)
	`, userId, postId)

	return err
}

func (repository *PostPostgres) UnlikePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE posts 
		SET likes=likes-1, liked_by=array_remove(liked_by, $1)
		WHERE id=$2 AND $1=ANY(liked_by)
	`, userId, postId)

	return err
}

func (repository *PostPostgres) DislikePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE posts 
		SET dislikes=dislikes+1, disliked_by=array_append(disliked_by, $1)
		WHERE id=$2 AND NOT $1=ANY(liked_by) AND NOT $1=ANY(disliked_by)
	`, userId, postId)

	return err
}

func (repository *PostPostgres) UndislikePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE posts 
		SET dislikes=dislikes-1, disliked_by=array_remove(disliked_by, $1)
		WHERE id=$2 AND $1=ANY(disliked_by)
	`, userId, postId)

	return err
}

func (repository *PostPostgres) DeletePost(ctx context.Context, userId, postId uint) error {
	_, err := repository.database.Exec(ctx, "DELETE FROM posts WHERE id = $1 AND parent_user_id=$2", postId, userId)
	return err
}

/*endregion*/

/*region comments*/

func (repository *PostPostgres) GetCommentsByPostId(ctx context.Context, postId uint, offset uint) ([]Entities.Comment, error) {
	rows, err := repository.database.Query(ctx, `
        SELECT id, parent_post_id, data, date, parent_user_id, likes, likes_by, dislikes, dislikes_by, files, parent_comment_id
        FROM comments
        WHERE parent_post_id = $1
        ORDER BY id DESC 
        LIMIT 20 OFFSET $2
    `, postId, offset)
	if err != nil {
		return nil, err
	}

	var comments = []Entities.Comment{}
	for rows.Next() {
		var comment Entities.Comment
		err = rows.Scan(&comment.ID, &comment.ParentPostID, &comment.Data, &comment.Date, &comment.ParentUserId, &comment.Likes, &comment.LikedBy, &comment.Dislikes, &comment.DislikedBy, &comment.Files, &comment.ParentCommentId)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repository *PostPostgres) CreateComment(ctx context.Context, userId uint, comment Entities.CommentDTO, date string) (commentId uint, err error) {
	var row pgx.Row
	if comment.ParentCommentId > 1 {
		row = repository.database.QueryRow(ctx, `INSERT INTO comments (parent_post_id, data, date, parent_user_id, files, parent_comment_id)
			VALUES ($1, $2, $6, $3, $4, $5) RETURNING id`,
			comment.ParentPostID, comment.Data, userId, comment.Files, comment.ParentCommentId, date)
	} else {
		row = repository.database.QueryRow(ctx, `INSERT INTO comments (parent_post_id, data, date, parent_user_id, files)
			VALUES ($1, $2, $5, $3, $4) RETURNING id`,
			comment.ParentPostID, comment.Data, userId, comment.Files, date)
	}
	err = row.Scan(&commentId)
	if err != nil {
		return 0, err
	}
	return commentId, nil
}

func (repository *PostPostgres) LikeComment(ctx context.Context, userID uint, commentID uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE comments 
		SET likes = likes + 1, likes_by = array_append(likes_by, $1)
		WHERE id = $2 AND NOT $1 = ANY(likes_by) AND NOT $1 = ANY(dislikes_by)
	`, userID, commentID)
	return err
}

func (repository *PostPostgres) UnlikeComment(ctx context.Context, userID uint, commentID uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE comments 
		SET likes = likes - 1, likes_by = array_remove(likes_by, $1)
		WHERE id = $2 AND $1 = ANY(likes_by)
	`, userID, commentID)
	return err
}

func (repository *PostPostgres) DislikeComment(ctx context.Context, userID uint, commentID uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE comments 
		SET dislikes = dislikes + 1, dislikes_by = array_append(dislikes_by, $1)
		WHERE id = $2 AND NOT $1 = ANY(likes_by) AND NOT $1 = ANY(dislikes_by)
	`, userID, commentID)
	return err
}

func (repository *PostPostgres) UndislikeComment(ctx context.Context, userID uint, commentID uint) error {
	_, err := repository.database.Exec(ctx, `UPDATE comments 
		SET dislikes = dislikes - 1, dislikes_by = array_remove(dislikes_by, $1)
		WHERE id = $2 AND $1 = ANY(dislikes_by)
	`, userID, commentID)
	return err
}

func (repository *PostPostgres) UpdateComment(ctx context.Context, userID uint, commentID uint, updateDTO Entities.CommentUpdateDTO) error {
	_, err := repository.database.Exec(ctx, `UPDATE comments 
		SET data = $1, files = $2 
		WHERE id = $3 AND parent_user_id = $4
	`, updateDTO.Data, updateDTO.Files, commentID, userID)
	return err
}

func (repository *PostPostgres) DeleteComment(ctx context.Context, userID uint, commentID uint) error {
	_, err := repository.database.Exec(ctx, `DELETE FROM comments 
		WHERE id = $1 AND parent_user_id = $2
	`, commentID, userID)
	return err
}

/*endregion*/

/*region Survey */

func (repository *PostPostgres) VoteInSurvey(ctx context.Context, userId uint, surveyId uint, votedFor []uint8) error {
	var votedForForPostgres string
	for _, v := range votedFor {
		votedForForPostgres = votedForForPostgres + "sl" + strconv.Itoa(int(v)) + "v=sl" + strconv.Itoa(int(v)) + "v+1," + "sl" + strconv.Itoa(int(v)) + "vby=array_append(sl" + strconv.Itoa(int(v)) + "vby, $1),"
	}
	query := `UPDATE surveys SET ` + votedForForPostgres + ` voted_by = array_append(voted_by, $1) WHERE parent_post_id = $2 AND NOT $1 = ANY(voted_by)`
	_, err := repository.database.Exec(ctx, query, userId, surveyId)
	return err
}

/*endregion*/
