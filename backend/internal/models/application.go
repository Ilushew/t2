package models

// Application данные заявки на маршрут
type Application struct {
	// Email для обратной связи
	Email string `json:"email"`
	// Текст заявки (комментарий/пожелания)
	Comment string `json:"comment"`
	// RouteName название маршрута
	RouteName string `json:"route_name"`
}
