package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/microsoft/go-mssqldb"
)

type UserDB struct {
	ID       int
	Name     string
	Gender   string
	Location string
	Email    string
	Phone    string
}

type RandomUserResponse struct {
	Results []User `json:"results"`
	Info    Info   `json:"info"`
}

type User struct {
	Gender   string   `json:"gender"`
	Name     Name     `json:"name"`
	Location Location `json:"location"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
}

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Location struct {
	Street  Street `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Street struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

type Info struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
	Version string `json:"version"`
}

const (
	workers    = 4
	totalUsers = 10000
)

var dbsql *sql.DB

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go setupRouter()
	go connection()
	<-c
	dbsql.Close()
}

func setupRouter() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/users", UsersHandler)
	r.Run()
}

func connection() {
	q := url.Values{}
	q.Add("database", "RandomUsers")

	url := url.URL{
		Scheme:   "sqlserver",
		Host:     "localhost:1433",
		User:     url.UserPassword("sa", "Admin123"),
		RawQuery: q.Encode(),
	}

	fmt.Println(url.String())

	db, err := sql.Open("sqlserver", url.String())
	dbsql = db
	if err != nil {
		fmt.Println(err)
	}
}

func GetUsersFromAPI() ([]User, error) {
	url := "https://randomuser.me/api/?results=500"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error al crear el request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la peticion")
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error al leer respuesta")
	}

	var users RandomUserResponse
	err = json.Unmarshal(resBytes, &users)
	if err != nil {
		fmt.Println("Error al parsear la respuesta")
	}

	if len(users.Results) > 0 {
		return users.Results, nil
	}

	return nil, fmt.Errorf("%v", err)
}

func WorkerAPI(ctx context.Context, cancel context.CancelFunc, jobs <-chan int, results chan<- []User, ops *atomic.Uint64) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-jobs:
			if !ok {
				return
			}

			if ctx.Err() != nil {
				return
			}

			users, err := GetUsersFromAPI()
			if err != nil {
				continue
			}
			current := ops.Load()
			if current == totalUsers {
				cancel()
				return
			}
			ops.Add(500)
			InsertUsers(users)

			select {
			case <-ctx.Done():
				return
			case results <- users:
			}
		}
	}
}

func ManageWorkerPool(ctx context.Context) []User {
	var wg sync.WaitGroup
	var ops atomic.Uint64

	jobs := make(chan int)
	results := make(chan []User)

	jobsctx, jobscancel := context.WithCancel(ctx)
	defer jobscancel()

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			WorkerAPI(jobsctx, jobscancel, jobs, results, &ops)
		}()
	}

	go func() {
		for {
			select {
			case <-jobsctx.Done():
				fmt.Println("Jobs cerrado")
				close(jobs)
				return
			default:
				current := ops.Load()
				if current < totalUsers {
					jobs <- 1
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var users []User
	for result := range results {
		users = append(users, result...)
	}

	fmt.Println("Total Usuarios")
	fmt.Println(ops.Load())

	return users
}

func InsertUsers(users []User) {
	fmt.Println(len(users), "Longitud de usuarios")
	for _, user := range users {
		userdb := UserDB{
			Name:     user.Name.First + user.Name.Last,
			Gender:   user.Gender,
			Location: user.Location.City + "|" + user.Location.State,
			Email:    user.Email,
			Phone:    user.Phone,
		}
		_, err := dbsql.Exec(`
			INSERT INTO USERS(name, gender, location, email, phone)
				VALUES(@p1, @p2, @p3, @p4, @p5)`, userdb.Name, userdb.Gender, userdb.Location, userdb.Email, userdb.Phone)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SelectUsers() {
	var user UserDB
	err := dbsql.QueryRow("SELECT * FROM USERS WHERE id = @p1", 1).Scan(&user.ID, &user.Name, &user.Gender, &user.Location, &user.Email, &user.Phone)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
}

func UsersHandler(ginctx *gin.Context) {
	ctx, cancel := context.WithTimeout(ginctx, 40*time.Second)
	defer cancel()

	users := ManageWorkerPool(ctx)

	ginctx.JSON(200, gin.H{
		"results": users,
	})
}
