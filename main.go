package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() { // dsn - источник данных
	dsn := "host=localhost user=postgres password=secret123 dbname=postgres port=5432 sslmode=disable"

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	//AutoMigrate отвечает за то чтобы в БД создалась модель Calculation
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}

// get - список операций ; post - передача выражения для расчета

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

// глобальная переменная - слайс , нужно инициализировать (иначе при попытке сайта получить запрос будет nul и ошибка возращатся )
//var calculations = []Calculation{}

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

// Основные методы ORM  -  Create , Read , Update , Delete  (crud)

func getCalculations(c echo.Context) error {
	/*Клиент (браузер)
	    ↓ HTTP GET запрос
	Функция getCalculations(c echo.Context)
	    ↓ c.JSON() формирует JSON ответ
	    ↓ return отправляет ошибку (или nil)
	Клиент получает JSON в теле ответа*/
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
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
	}

	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calculations"})
	}

	return c.JSON(http.StatusCreated, calc)

	/*                     куда          что
	calculations = append(calculations, calc) // Первый параметр функции - срез, в который надо добавить, а второй параметр - значение, которое нужно добавить
	return c.JSON(http.StatusCreated, calc)*/
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
	/*
		//перебираем слайс(историю) вычислений
		for i, calculation := range calculations {
			if calculation.ID == id { // и если ID равен тому который мы передали (id)
				calculations[i].Expression = req.Expression // обновляем рузультат
				calculations[i].Result = result
				return c.JSON(http.StatusOK, calculations[i])
			}
		}
	*/
	var Calc Calculation // выражение которое хотим заменить
	if err := db.First(&Calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Could not find expression"}) //выражение не найдено (ошибка со сторон клиента)
	}

	Calc.Expression = req.Expression
	Calc.Result = result

	if err := db.Save(&Calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculations"})
	}

	return c.JSON(http.StatusOK, Calc)
}

func deleteCalculations(c echo.Context) error {
	id := c.Param("id")

	/*
		for i, calculation := range calculations {
			if calculation.ID == id {
				calculations = append(calculations[:i], calculations[i+1:]...) // удаление элемента из слайса
				return c.NoContent(http.StatusNoContent)                       //  NoContent - возращает ответ без тела , и статус код
			}
		}
	*/
	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculations"})
	}
	return c.NoContent(http.StatusNoContent)
	//return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcculation not found"})
}

func main() {
	initDB()
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
