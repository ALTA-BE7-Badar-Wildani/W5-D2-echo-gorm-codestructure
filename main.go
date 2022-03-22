package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var (
	DB *gorm.DB
)


// Entity config
type Config struct {
	DB_Username string
	DB_Password string
	DB_Port string
	DB_Host string
	DB_Name string
}

// Init database
func InitDB()  {
	config := Config{
		DB_Username: "root",
		DB_Password: "root",
		DB_Port: "3306",
		DB_Host: "localhost",
		DB_Name: "alta_echo_gorm",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	fmt.Println(dsn)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't connect to database")
	}
}

// Entity User
type User struct {
	gorm.Model
	Name string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

//  Init Migration
func InitMigration() {
	DB.AutoMigrate(&User{})
}



// GET	/users
// Get all User
func GetUsersController(c echo.Context) error {
	// Get Users
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all user",
		"users": users,
	})
}


// GET /users/:id
// Get user controller
func GetUserController(c echo.Context) error {
	// Get User
	id := c.Param("id")
	user := User{}
	tx := DB.Find(&user, id)
	if tx.RowsAffected <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{} {
			"message": "User with ID " + id + " is not found",
		})
	}
	// Return Response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success getting a user",
		"users": user,
	})
}

// POST /users
// Create new user
func CreateUserController(c echo.Context) error {
	// Retrieve request body
	user := User{}
	c.Bind(&user)
	// Create user
	err := DB.Save(&user)
	if err.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	// return response
	return c.JSON(http.StatusOK, map[string]interface{} {
		"message": "success create new user",
		"user": user,
	})
}


// PUT	/users/:id
// Update a spesific User
func UpdateUserController(c echo.Context) error {
	// Get User id
	id := c.Param("id")
	user := User{}
	tx := DB.Find(&user, id)
	if tx.RowsAffected <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{} {
			"message": "User with ID " + id + " is not found",
		})
	}

	// Update
	c.Bind(&user)
	tx = DB.Save(&user)
	if tx.RowsAffected <= 0 {
		return c.JSON(http.StatusInternalServerError, map[string]interface{} {
			"message": "User with ID " + id + " can't be updated",
		})
	}
	// return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success updating user",
		"user": user,
	})
}


// DELETE /users/:id
// Delete spesific user
func DeleteUserController(c echo.Context) error {
	// Get Id param
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{} {
			"message": "ID param is not provided",
		})
	}
	// Delete user
	tx := DB.Delete(&User{}, id)
	if tx.RowsAffected <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{} {
			"message": "Can't delete user with provided ID",
		})
	}
	// return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success deleting a user",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}



func init() {
	InitDB()
	InitMigration()
}


func main() {
	e := echo.New()
	e.GET("/users", GetUsersController)
	e.POST("/users", CreateUserController)
	e.GET("/users/:id", GetUserController)
	e.PUT("/users/:id", UpdateUserController)
	e.DELETE("/users/:id", DeleteUserController)

	e.Logger.Fatal(e.Start(":8000"))
}