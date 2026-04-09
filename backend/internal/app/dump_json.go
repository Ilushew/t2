package app

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PlaceDump структура для дампа места
type PlaceDump struct {
	Name              string  `json:"name"`
	Price             string  `json:"price"`
	Time              float64 `json:"time"`
	TypesOfMovement   string  `json:"types_of_movement"`
	Category          string  `json:"category"`
	LatStart          float64 `json:"lat_start"`
	LonStart          float64 `json:"lon_start"`
	LatEnd            float64 `json:"lat_end"`
	LonEnd            float64 `json:"lon_end"`
	IsIndoor          bool    `json:"is_indoor"`
	WithChild         bool    `json:"with_child"`
	WithPets          bool    `json:"with_pets"`
	Description       string  `json:"description"`
}

// importFromJSON импортирует данные из JSON-файла в таблицу places
func importFromJSON(ctx context.Context, pool *pgxpool.Pool, inputPath string) bool {
	// Проверяем существование файла
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return false
	}

	log.Printf("JSON-импорт: читаю файл %s", inputPath)

	// Читаем файл
	jsonData, err := os.ReadFile(inputPath)
	if err != nil {
		log.Printf("JSON-импорт пропущен: ошибка чтения файла: %v", err)
		return false
	}

	// Парсим JSON
	var places []PlaceDump
	if err := json.Unmarshal(jsonData, &places); err != nil {
		log.Printf("JSON-импорт пропущен: ошибка парсинга JSON: %v", err)
		return false
	}

	log.Printf("JSON-импорт: найдено %d записей", len(places))

	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Printf("JSON-импорт пропущен: транзакция: %v", err)
		return false
	}
	defer tx.Rollback(ctx)

	// Очищаем таблицу
	_, err = tx.Exec(ctx, `TRUNCATE TABLE places RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Printf("JSON-импорт: ошибка очистки таблицы: %v", err)
		return false
	}

	// Вставляем данные
	for i, p := range places {
		_, err = tx.Exec(ctx, `
			INSERT INTO places (name, price, time, types_of_movement, category,
			                    lat_start, lon_start, lat_end, lon_end,
			                    is_indoor, with_child, with_pets, description)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, p.Name, p.Price, p.Time, p.TypesOfMovement, p.Category,
			p.LatStart, p.LonStart, p.LatEnd, p.LonEnd,
			p.IsIndoor, p.WithChild, p.WithPets, p.Description)
		if err != nil {
			log.Printf("ошибка вставки строки %d (name=%s): %v", i+1, p.Name, err)
			return false
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("JSON-импорт пропущен: commit: %v", err)
		return false
	}

	log.Printf("JSON-импорт УСПЕШЕН: импортировано %d записей", len(places))
	return true
}
