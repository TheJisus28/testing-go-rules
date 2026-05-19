package models

// APIErrorResponse representa el cuerpo de error estándar de la API (Swagger).
type APIErrorResponse struct {
	Status       string `json:"status" example:"validation_error"`
	StatusReason string `json:"status_reason" example:"invalid request body"`
	Data         any    `json:"data"`
}
