package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/db"
	"github.com/Eco-Sort/eco_sort_backend/repository/mariadb"
	"github.com/joho/godotenv"
	"github.com/thatisuday/commando"
)

func Init() {
	if godotenv.Load(".env") != nil {
		fmt.Println("Unable to load .env file, using global varibale")
	}
	db.InitMariadb()
	db.InitGcStorage()
}

func main() {
	commando.
		Register("seed:generate").
		SetDescription("generating seed for category, sorting, and 1 admin user").
		SetShortDescription("generating seed for category, sorting, and 1 admin user").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			Init()
			db.Mariadb.AutoMigrate(
				&domain.User{},
				&domain.Sorting{},
				&domain.Category{},
				&domain.ImageObject{},
				&domain.Garbage{},
			)
			masterCategoryRepo := mariadb.NewMariadbCategoryRepository(db.Mariadb)
			masterSortingRepo := mariadb.NewMariadbSortingRepository(db.Mariadb)

			categoryContent, err := os.ReadFile("./seed/category.json")
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}

			var category []domain.CategoryRequest
			err = json.Unmarshal(categoryContent, &category)
			if err != nil {
				log.Fatal("Error when unmarshaling json category: ", err)
			}

			sortingContent, err := os.ReadFile("./seed/sorting.json")
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}

			var sorting []domain.SortingRequest
			err = json.Unmarshal(sortingContent, &sorting)
			if err != nil {
				log.Fatal("Error when unmarshaling json sorting: ", err)
			}

			for _, s := range sorting {
				res, err := masterSortingRepo.Create(s)
				if err != nil {
					log.Fatal("Error when creating sorting: ", err)
				}
				fmt.Println(res)
			}

			for _, s := range category {
				res, err := masterCategoryRepo.Create(s)
				if err != nil {
					log.Fatal("Error when creating category: ", err)
				}
				fmt.Println(res)
			}
		})
	commando.Parse(nil)
}
