package repository

import (
	"context"
	"fmt"

	"github.com/ilushew/udmurtia-trip/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	sq "github.com/Masterminds/squirrel"
)

type PlaceRepository struct {
	pool     *pgxpool.Pool
	psq      sq.StatementBuilderType
	imageRepo *PlaceImageRepository
}

func NewPlaceRepository(pool *pgxpool.Pool, imageRepo *PlaceImageRepository) *PlaceRepository {
	return &PlaceRepository{
		pool:     pool,
		psq:      sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		imageRepo: imageRepo,
	}
}

// GetByIDs возвращает места по списку ID в том же порядке + подгружает картинки
func (r *PlaceRepository) GetByIDs(ctx context.Context, ids []int) ([]models.Place, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `
		SELECT id, name, name_label, price, time, types_of_movement, category,
		       lat_start, lon_start, lat_end, lon_end,
		       is_indoor, with_child, with_pets, description, description_label
		FROM places
		WHERE id = ANY($1)
		ORDER BY array_position($1, id)
	`

	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query places: %w", err)
	}
	defer rows.Close()

	places, err := scanPlaces(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan places: %w", err)
	}

	// Подгружаем картинки если есть репозиторий
	if r.imageRepo != nil && len(places) > 0 {
		imagesMap, err := r.imageRepo.GetByPlaceIDs(ctx, ids)
		if err == nil {
			for i := range places {
				places[i].Images = imagesMap[places[i].ID]
			}
		}
	}

	return places, nil
}

// GetAll возвращает все места из БД (для админки)
func (r *PlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
	query, args, err := r.psq.
		Select("id", "name", "name_label", "price", "time", "types_of_movement", "category",
			"lat_start", "lon_start", "lat_end", "lon_end",
			"is_indoor", "with_child", "with_pets", "description", "description_label").
		From("places").
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query places: %w", err)
	}
	defer rows.Close()

	places, err := scanPlaces(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan places: %w", err)
	}

	// Подгружаем картинки
	if r.imageRepo != nil && len(places) > 0 {
		ids := make([]int, len(places))
		for i, p := range places {
			ids[i] = p.ID
		}
		imagesMap, err := r.imageRepo.GetByPlaceIDs(ctx, ids)
		if err == nil {
			for i := range places {
				places[i].Images = imagesMap[places[i].ID]
			}
		}
	}

	return places, nil
}

func scanPlaces(rows pgx.Rows) ([]models.Place, error) {
	var places []models.Place

	for rows.Next() {
		var p models.Place
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.NameLabel,
			&p.Price,
			&p.Time,
			&p.TypesOfMovement,
			&p.Category,
			&p.LatStart,
			&p.LonStart,
			&p.LatEnd,
			&p.LonEnd,
			&p.IsIndoor,
			&p.WithChild,
			&p.WithPets,
			&p.Description,
			&p.DescriptionLabel,
		); err != nil {
			return nil, fmt.Errorf("failed to scan place: %w", err)
		}
		places = append(places, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return places, nil
}

// UpdatePlace обновляет поля места по ID
func (r *PlaceRepository) UpdatePlace(ctx context.Context, id int, updates map[string]any) error {
	query, args, err := r.psq.
		Update("places").
		SetMap(updates).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

// InsertPlace добавляет новое место в БД
func (r *PlaceRepository) InsertPlace(ctx context.Context, place *models.Place) error {
	query, args, err := r.psq.
		Insert("places").
		Columns(
			"name", "name_label", "price", "time", "types_of_movement", "category",
			"lat_start", "lon_start", "lat_end", "lon_end",
			"is_indoor", "with_child", "with_pets", "description", "description_label",
		).
		Values(
			place.Name, place.NameLabel, place.Price, place.Time,
			place.TypesOfMovement, place.Category,
			place.LatStart, place.LonStart, place.LatEnd, place.LonEnd,
			place.IsIndoor, place.WithChild, place.WithPets, place.Description, place.DescriptionLabel,
		).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

// DeletePlace удаляет место по ID
func (r *PlaceRepository) DeletePlace(ctx context.Context, id int) error {
	query, args, err := r.psq.
		Delete("places").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
