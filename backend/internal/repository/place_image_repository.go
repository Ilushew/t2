package repository

import (
	"context"
	"fmt"

	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	sq "github.com/Masterminds/squirrel"
)

type PlaceImageRepository struct {
	pool *pgxpool.Pool
	psq  sq.StatementBuilderType
}

func NewPlaceImageRepository(pool *pgxpool.Pool) *PlaceImageRepository {
	return &PlaceImageRepository{
		pool: pool,
		psq:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// GetByPlaceIDs возвращает картинки для указанных ID мест, сгруппированные по place_id
func (r *PlaceImageRepository) GetByPlaceIDs(ctx context.Context, placeIDs []int) (map[int][]models.PlaceImage, error) {
	if len(placeIDs) == 0 {
		return make(map[int][]models.PlaceImage), nil
	}

	query, args, err := r.psq.
		Select("place_id", "id", "filename", "sort_order").
		From("place_images").
		Where(sq.Eq{"place_id": placeIDs}).
		OrderBy("place_id", "sort_order").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	result := make(map[int][]models.PlaceImage)
	for rows.Next() {
		var img models.PlaceImage
		if err := rows.Scan(&img.PlaceID, &img.ID, &img.Filename, &img.SortOrder); err != nil {
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}
		result[img.PlaceID] = append(result[img.PlaceID], img)
	}

	return result, nil
}

// Create добавляет картинку
func (r *PlaceImageRepository) Create(ctx context.Context, img *models.PlaceImage) error {
	query, args, err := r.psq.
		Insert("place_images").
		Columns("place_id", "filename", "sort_order").
		Values(img.PlaceID, img.Filename, img.SortOrder).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(&img.ID)
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	return nil
}

// Delete удаляет картинку по ID
func (r *PlaceImageRepository) Delete(ctx context.Context, id int) error {
	query, args, err := r.psq.
		Delete("place_images").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
