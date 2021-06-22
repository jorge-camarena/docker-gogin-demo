package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type Task struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Text       string `json:"text"`
	DatePosted string `json:"datePosted"`
	DueDate    string `json:"dueDate"`
	Reminder   bool   `json:"reminder"`
}

//Route handlers for ToDo List App
func createUser(c *gin.Context) {
	body, _ := c.GetRawData()
	var user User
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO users (id, name, email, password, age) VALUES ($1, $2, $3, $4, $5)", user.Id, user.Name, user.Email, user.Password, user.Age)

	if err != nil {
		log.Fatal(err)
	}
}

func loginUser(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	rows, _ := db.Query("SELECT name FROM users WHERE email = ($1) AND password = ($2)", email, password)
	var name string
	for rows.Next() {
		rows.Scan(&name)
	}
	if len(name) == 0 {
		c.JSON(404, gin.H{
			"message":        "Could not login because email and user does not exist",
			"loggedInStatus": false,
		},
		)
	} else {
		c.JSON(200, gin.H{
			"message":        "Successfully logged in",
			"loggedInStatus": true,
		})
	}
}

func postTask(c *gin.Context) {
	body, _ := c.GetRawData()
	var task Task
	err := json.Unmarshal(body, &task)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO tasks (id, name, email, text, dateposted, duedate, reminder) VALUES ($1, $2, $3, $4, $5, $6, $7)", task.Id, task.Name, task.Email, task.Text, task.DatePosted, task.DueDate, task.Reminder)
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(c *gin.Context) {
	email := c.Query("email")
	rows, _ := db.Query("SELECT text FROM tasks WHERE email = ($1)", email)
	var tasks []string
	for rows.Next() {
		var text string
		rows.Scan(&text)
		tasks = append(tasks, text)
	}
	c.JSON(
		200,
		tasks,
	)
}

func deleteTask(c *gin.Context) {
	task_id := c.Query("id")
	loggedInStatus := c.Query("loggedInStatus")
	b, _ := strconv.ParseBool(loggedInStatus)

	if b == true {
		_, err := db.Exec("DELETE FROM tasks WHERE id = ($1)", task_id)
		fmt.Println(err)
		c.JSON(200, gin.H{
			"message": "Successfully deleted task",
		})
	} else {
		c.JSON(404, gin.H{
			"message": "You must be logged in and authenticated in order to delete a task",
		})
	}
}

func deleteUser(c *gin.Context) {
	user_id := c.Query("id")
	loggedInStatus := c.Query("loggedInStatus")
	b, _ := strconv.ParseBool(loggedInStatus)

	if b == true {
		db.Exec("DELETE FROM users WHERE id = ($1)", user_id)
		c.JSON(200, gin.H{
			"message": "Successfully deleted user",
		})
	} else {
		c.JSON(404, gin.H{
			"message": "You must be logged in and authenticated in order to delete your account",
		})
	}
}

func main() {
	godotenv.Load()

	var (
		host, _     = os.LookupEnv("HOST")
		port        = 5432
		user, _     = os.LookupEnv("USER")
		password, _ = os.LookupEnv("PASSWORD")
		dbname, _   = os.LookupEnv("DBNAME")
	)

	fmt.Println(host)
	//print statement for debugging
	fmt.Println("Hello Word")
	router := gin.Default()

	//Serving static assets
	router.Static("/home", "./client")

	//Connecting to PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	//Initialize all necessary routes
	router.GET("/api/login", loginUser)
	router.GET("/api/get-tasks", getTasks)
	router.POST("/api/create-user", createUser)
	router.POST("/api/post-task", postTask)
	router.DELETE("/api/delete-task", deleteTask)
	router.DELETE("/api/delete-user", deleteUser)

	//Run server on port :8080
	router.Run()
}
