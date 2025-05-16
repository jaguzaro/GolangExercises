package controllers

import (
	"UsersAPI/db"
	dtb "UsersAPI/db"
	"UsersAPI/db/models"
	stc "UsersAPI/structs"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"
)

func GetRandomUser(ctx context.Context) ([]stc.User, error) {
	/* newCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel() */
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://randomuser.me/api/?results=500", nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear peticion")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la peticion: %w", err)
	}
	defer resp.Body.Close()

	var randomUser stc.RandomUserResponse

	if err := json.NewDecoder(resp.Body).Decode(&randomUser); err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta")
	}

	if len(randomUser.Results) > 0 {

		return randomUser.Results, nil
	}
	return nil, fmt.Errorf("no se encontraron datos")
}

func RequestProcessor(ctx context.Context, jobs <-chan int, results chan<- []stc.User, jobctx context.Context) {
	for {
		select {
		case <-jobctx.Done():
			return
		case <-jobs:
			user, err := GetRandomUser(ctx)
			if err != nil {
				fmt.Println("error al obtener usuario: " + err.Error())
				continue
			}
			results <- user
		}
	}
}

func JobProducer(jobctx context.Context, jobs chan<- int) {
	for {
		select {
		case <-jobctx.Done():
			fmt.Println("productor detenido")
			close(jobs)
			return
		default:
			jobs <- 1
			time.Sleep(1 * time.Second)
		}
	}
}

func ManageProcessor(ctx context.Context, done chan<- []stc.User) {
	/* numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU) */
	numWorkers := 500
	totalUsers := 12000

	jobs := make(chan int)
	results := make(chan []stc.User)

	jobctx, jobcancel := context.WithCancel(ctx)
	defer jobcancel()
	var users []stc.User
	var wg sync.WaitGroup
	//var mu sync.Mutex

	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			RequestProcessor(ctx, jobs, results, jobctx)
		}()
	}

	go JobProducer(jobctx, jobs)

	for res := range results {
		users = append(users, res...)
		fmt.Println(len(users))
		if len(users) >= totalUsers {
			fmt.Println("Maximo de usuarios")
			jobcancel()
			break
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	InsertIntoDB(ctx, users)

	done <- users
}

func InsertIntoDB(ctx context.Context, users []stc.User) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	db := dtb.GetDBController().WithContext(ctx)

	jobs := make(chan stc.User, len(users))

	workers := 500
	wg.Add(500)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				for k := 0; k < 3; k++ {
					mutex.Lock()
					result := db.Create(
						&models.Users{
							Name:     j.Name.First + " " + j.Name.Last,
							Gender:   j.Gender,
							Location: j.Location.City + " " + j.Location.State,
							Email:    j.Email,
							Phone:    j.Phone,
						})
					mutex.Unlock()
					if result.Error != nil {
						fmt.Println("Result: ", result)
						continue
					}
					break
				}
			}
		}()
	}

	for _, user := range users {
		jobs <- user
	}

	close(jobs)
	wg.Wait()
}

//Worker Function
//Manager

func workerSelectUser(db *gorm.DB, jobs <-chan int, result chan<- models.Users) {
	var mutex sync.Mutex
	for j := range jobs {
		var user models.Users
		mutex.Lock()
		resultDB := db.First(&user, j)
		mutex.Unlock()
		if resultDB.Error != nil {
			fmt.Println("error to get user from db", resultDB.Error)
			continue
		}
		result <- user
	}
}

func ManageQueries(ctx context.Context, done chan<- []models.Users) {
	workers := 100
	totalUsers := 10000

	var wg sync.WaitGroup
	jobs := make(chan int)
	results := make(chan models.Users, totalUsers)

	wg.Add(workers)
	db := db.GetDBController().WithContext(ctx)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			workerSelectUser(db, jobs, results)
		}()
	}

	for i := 1; i <= totalUsers; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var users []models.Users
	for result := range results {
		users = append(users, result)
	}

	done <- users
}
