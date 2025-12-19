package calculationService

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type CalcService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationByID(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo Repository
}

func NewCalcService(r Repository) CalcService {
	return &calcService{repo: r}
}

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return Calculation{}, fmt.Errorf("invalid expression: %w", err)
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return Calculation{}, fmt.Errorf("evaluation error: %w", err)
	}

	calc := Calculation{
		ID:         "",
		Expression: expression,
		Result:     fmt.Sprintf("%v", result),
	}
	created, err := s.repo.Create(calc)
	if err != nil {
		return Calculation{}, err
	}
	return created, nil
}

func (s *calcService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAll()
}

func (s *calcService) GetCalculationByID(id string) (Calculation, error) {
	return s.repo.GetByID(id)
}

func (s *calcService) UpdateCalculation(id, expression string) (Calculation, error) {

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return Calculation{}, err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         id,
		Expression: expression,
		Result:     fmt.Sprintf("%v", result),
	}
	return s.repo.Update(calc)
}

func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.Delete(id)
}
