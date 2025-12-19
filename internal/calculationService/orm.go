package calculationService

// структура выражения которая будет хранится в бд
type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"` // выражение
	Result     string `json:"result"`
}

// структура запроса на вычисления
type CalculationRequest struct {
	Expression string `json:"expression"`
}
