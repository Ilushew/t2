// +build ignore

// Скрипт для создания дампа БД в JSON
// Использование: go run scripts/dump_db.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/udmurtia_trip?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Ошибка пинга БД: %v", err)
	}

	outputPath := "static/places_dump.json"

	rows, err := pool.Query(ctx, `
		SELECT name, price, time, types_of_movement, category,
		       lat_start, lon_start, lat_end, lon_end,
		       is_indoor, with_child, with_pets, description
		FROM places ORDER BY id
	`)
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}
	defer rows.Close()

	var places []PlaceDump
	for rows.Next() {
		var p PlaceDump
		if err := rows.Scan(&p.Name, &p.Price, &p.Time, &p.TypesOfMovement, &p.Category,
			&p.LatStart, &p.LonStart, &p.LatEnd, &p.LonEnd,
			&p.IsIndoor, &p.WithChild, &p.WithPets, &p.Description); err != nil {
			log.Fatalf("Ошибка сканирования: %v", err)
		}
		places = append(places, p)
	}

	jsonData, err := json.MarshalIndent(places, "", "  ")
	if err != nil {
		log.Fatalf("Ошибка JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Ошибка записи файла: %v", err)
	}

	log.Printf("✓ Дамп успешно создан: %s (%d записей)", outputPath, len(places))
}
