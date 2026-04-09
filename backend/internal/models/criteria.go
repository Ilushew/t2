package models

// TripCriteria — критерии для подбора маршрута
type TripCriteria struct {
	Duration   string   `json:"duration"`            // длительность поездки
	Company    string   `json:"company"`             // компания: один, семья, друзья
	HasCar     bool     `json:"has_car"`             // есть ли автомобиль
	Budget     string   `json:"budget"`              // бюджет: экономный, средний, не важен
	Interests  []string `json:"interests,omitempty"` // интересы (множественный выбор)
	WithPets   bool     `json:"with_pets"`           // путешествие с животными
}
