package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// get - список операций ; post - передача выражения для расчета

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
var calculations = []Calculation{}

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
	return c.JSON(http.StatusOK, calculations)
	// возращает JSON со статус кодом OK И слайс предыдуших выражений
}

func postCalculations(c echo.Context) error {
	//тело запроса начало
	var req CalculationRequest
	if err := c.Bind(&req); err != nil { //если произошла ошибка с декодировкой то что мы передали то ошибка
		// не тот запрос
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := CalculateExpression(req.Expression)
	//не то выражение
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	//конец

	calc := Calculation{ // создаем новую структуру из того что есть
		ID:         uuid.NewString(),
		Expression: req.Expression, // берем то выражение которое пришло с сайта
		Result:     result,
	} //                     куда          что
	calculations = append(calculations, calc) // Первый параметр функции - срез, в который надо добавить, а второй параметр - значение, которое нужно добавить
	return c.JSON(http.StatusCreated, calc)
}

func patchCalculations(c echo.Context) error {
	id := c.Param("id") // достает из заголовка / ссылкаи id которое хотим обновить
	//тело запроса начало
	var req CalculationRequest //передаем запрос
	//декодируем
	if err := c.Bind(&req); err != nil { //если произошла ошибка с декодировкой то что мы передали то ошибка
		// не тот запрос
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := CalculateExpression(req.Expression) // считаем новый результат
	//не то выражение
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	//конец

	//перебираем слайс(историю) вычислений
	for i, calculation := range calculations {
		if calculation.ID == id { // и если ID равен тому который мы передали (id)
			calculations[i].Expression = req.Expression // обновляем рузультат
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcculation not found"}) // вычисление не найдено
}

func deleteCalculations(c echo.Context) error {
	id := c.Param("id")

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations = append(calculations[:i], calculations[i+1:]...) // удаление элемента из слайса
			return c.NoContent(http.StatusNoContent)                       //  NoContent - возращает ответ без тела , и статус код
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcculation not found"})
}

func main() {
	e := echo.New() // инициализировать (создать) обработчик echo

	e.Use(middleware.CORS())   // CORS - на нашем пк сайт пытается отправить данные на сервер . цепь в нашем запросек
	e.Use(middleware.Logger()) // каждый запрос перед обработкой пройдет через него и попадет в консоль

	//все get/post запросы     путь             функция обработчик
	e.GET("/calculations", getCalculations)
	e.POST("/calculations", postCalculations)
	e.PATCH("/calculations/:id", patchCalculations)
	e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8080")
}
