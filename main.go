package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/YONG-LIN-LIANG/go-fiber-postgres/models"
	"github.com/YONG-LIN-LIANG/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct{
	Author		string			`json:"author"`
	Title			string			`json:"title"`
	Publisher	string			`json:"publisher"`
}

type Repository struct{
	DB *gorm.DB
}

func(r *Repository) GetBookByID(context *fiber.Ctx) error{
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("the ID is", id)
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"could not get the book",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"book id fetched successfully",
		"data":bookModel,
	})
	return nil
}

func(r *Repository) DeleteBook(context *fiber.Ctx) error{
	// 使用package model的Books type
	bookModel := models.Books{}
	id := context.Params("id")
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"could not delete book",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"books delete successfully",
	})
	return nil
}

func(r *Repository) GetBooks(context *fiber.Ctx) error{
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"books fetched successfully",
		"data":	bookModels,
	})
	return nil
}

// 最後那個error是一定要return error的意思
func(r *Repository) CreateBook(context *fiber.Ctx) error{
	book := Book{}
	// 根據Book struct type decode json
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create book"})
		return err
	}
	
	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"book has been added"})
	return nil
}

// 這樣寫是因為要定義 SetupRoutes 是 r 的方法
func(r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	// 使用storage package中的Config type，賦值後丟給config
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		DBName: os.Getenv("DB_NAME"),
	}
	// 建立資料庫連線
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}
	// 如果連線無錯誤，呼叫 models package 的 MigrateBooks func 來遷移(Migrate)我們的表到資料庫中
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}