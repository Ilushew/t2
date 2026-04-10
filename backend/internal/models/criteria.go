package models

// TripCriteria — критерии для подбора маршрута
type TripCriteria struct {
	Duration   string   `json:"duration"`
	Company    string   `json:"company"`
	HasCar     bool     `json:"has_car"`
	Budget     int      `json:"budget"`
	Interests  []string `json:"interests,omitempty"`
	WithPets   bool     `json:"with_pets"`
	Query      string   `json:"query"`
}

// BudgetLabels — текстовые подписи для значений бюджета
var BudgetLabels = []string{"Экономный", "Средний", "Не важен"}

// BudgetLabel возвращает текстовую подпись для значения бюджета
func (t TripCriteria) BudgetLabel() string {
	if t.Budget < 0 || t.Budget >= len(BudgetLabels) {
		return BudgetLabels[1] // по умолчанию "Средний"
	}
	return BudgetLabels[t.Budget]
}
