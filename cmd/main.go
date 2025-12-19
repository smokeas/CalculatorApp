package main

import (
	"log"

	"CalculatorApp/internal/calculationService"
	"CalculatorApp/internal/db"
	"CalculatorApp/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// миграция модели
	if err := database.AutoMigrate(&calculationService.Calculation{}); err != nil {
		log.Fatalf("failed to automigrate: %v", err)
	}

	repo := calculationService.NewGormRepository(database)
	service := calculationService.NewCalcService(repo)
	h := handlers.NewHandlers(service)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/calculations", h.GetCalculations)
	e.POST("/calculations", h.PostCalculation)
	e.PATCH("/calculations/:id", h.PatchCalculation)
	e.DELETE("/calculations/:id", h.DeleteCalculation)

	log.Fatal(e.Start(":8080"))
}
