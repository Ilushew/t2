package repository

import (
	"context"
	"fmt"

	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	sq "github.com/Masterminds/squirrel"
)

type PlaceCommentRepository struct {
	pool *pgxpool.Pool
	psq  sq.StatementBuilderType
}

func NewPlaceCommentRepository(pool *pgxpool.Pool) *PlaceCommentRepository {
	return &PlaceCommentRepository{
		pool: pool,
		psq:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// GetByPlaceID возвращает все комментарии для маршрута, отсортированные по дате
func (r *PlaceCommentRepository) GetByPlaceID(ctx context.Context, placeID int) ([]models.PlaceComment, error) {
	query, args, err := r.psq.
		Select("id", "place_id", "author", "text", "created_at").
		From("place_comments").
		Where(sq.Eq{"place_id": placeID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	return scanComments(rows)
}

// Create добавляет новый комментарий
func (r *PlaceCommentRepository) Create(ctx context.Context, comment *models.PlaceComment) error {
	query, args, err := r.psq.
		Insert("place_comments").
		Columns("place_id", "author", "text").
		Values(comment.PlaceID, comment.Author, comment.Text).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

// Count возвращает количество комментариев для маршрута
func (r *PlaceCommentRepository) Count(ctx context.Context, placeID int) (int, error) {
	query, args, err := r.psq.
		Select("COUNT(*)").
		From("place_comments").
		Where(sq.Eq{"place_id": placeID}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}

func scanComments(rows interface{ Next() bool; Scan(...any) error; Err() error }) ([]models.PlaceComment, error) {
	var comments []models.PlaceComment

	for rows.Next() {
		var c models.PlaceComment
		if err := rows.Scan(&c.ID, &c.PlaceID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return comments, nil
}
