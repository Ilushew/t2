package services

import "time"

// Place описывает точку маршрута
// Эти поля должны совпадать с тем, что ты ожидаешь от реальной ML-модели
type Place struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    Rating      float64 `json:"rating"`
    Category    string  `json:"category"`
}

// MockMLService имитирует работу нейросети
type MockMLService struct{}

// NewMockMLService создает новый экземпляр заглушки
func NewMockMLService() *MockMLService {
    return &MockMLService{}
}

// GetRecommendations возвращает хардкодный список мест
// Сигнатура функции должна совпадать с будущим реальным сервисом
func (s *MockMLService) GetRecommendations(pref string) ([]Place, time.Duration, error) {
    // База мест Удмуртии для тестов
    places := []Place{
        {
            ID:          1,
            Name:        "Музейно-выставочный комплекс им. М.Т. Калашникова",
            Description: "Современный музей, посвященный истории оружия и биографии конструктора.",
            Latitude:    56.8498,
            Longitude:   53.2001,
            Rating:      4.9,
            Category:    "culture",
        },
        {
            ID:          2,
            Name:        "Архитектурно-этнографический музей Лудорвай",
            Description: "Музей деревянного зодчества под открытым небом. Ветряные мельницы и старинные избы.",
            Latitude:    56.7269,
            Longitude:   53.3550,
            Rating:      4.8,
            Category:    "culture",
        },
        {
            ID:          3,
            Name:        "Ижевский зоопарк",
            Description: "Один из крупнейших зоопарков России. Отлично подходит для семейного отдыха.",
            Latitude:    56.8389,
            Longitude:   53.1927,
            Rating:      4.7,
            Category:    "family",
        },
        {
            ID:          4,
            Name:        "Водопад Шумиловский (Шаркан)",
            Description: "Живописный водопад в лесу. Популярное место для фотосессий и прогулок.",
            Latitude:    57.0569,
            Longitude:   53.9069,
            Rating:      4.6,
            Category:    "nature",
        },
        {
            ID:          5,
            Name:        "Гора Лысова (Ижевск)",
            Description: "Панорамная видовая точка на берегу Ижевского пруда. Лучшее место для встречи заката.",
            Latitude:    56.8598,
            Longitude:   53.2289,
            Rating:      4.8,
            Category:    "nature",
        },
    }

    // Простая логика фильтрации по предпочтениям
    // Если пользователь выбрал конкретный тип, показываем соответствующие места
    // Если "all" или не совпало — показываем всё
    if pref != "" && pref != "all" {
        var filtered []Place
        for _, p := range places {
            if p.Category == pref {
                filtered = append(filtered, p)
            }
        }
        // Если по фильтру ничего нет, возвращаем всё, чтобы не было пусто
        if len(filtered) > 0 {
            return filtered, 500 * time.Millisecond, nil
        }
    }

    return places, 500 * time.Millisecond, nil
}