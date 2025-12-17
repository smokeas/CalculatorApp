package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// структура выражения которая будет хранится в бд
type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"` // выражение
	Result     string `json:"result"`
}

// структура запроса на вычисления
type CalculationRequest struct {
	Expression string `json:"expression"`
}

// глобальная переменная - слайс , нужно инициализировать (иначе при попытке сайта получить запрос будет nul и ошибка возращатся )
var Calculations = []Calculation{}

// принимает строку       возращает строку и ошибку
func CalculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression) // создаем выражени (1+1)
	if err != nil {
		return "", err // передали 1++1
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err
}

func getCalculations(c echo.Context) error {
	/*Клиент (браузер)
	    ↓ HTTP GET запрос
	Функция getCalculations(c echo.Context)
	    ↓ c.JSON() формирует JSON ответ
	    ↓ return отправляет ошибку (или nil)
	Клиент получает JSON в теле ответа*/
	return c.JSON(http.StatusOK, Calculations)
}

func main() {
	e := echo.New() // инициализировать (создать) обработчик echo

	e.Use(middleware.CORS())   // CORS - на нашем пк сайт пытается отправить данные на сервер . цепь в нашем запросек
	e.Use(middleware.Logger()) // каждый запрос перед обработкой пройдет через него и попадет в консоль

	//      путь             функция
	e.GET("/calculations", getCalculations)

	e.Start("localhost:8080")
}
