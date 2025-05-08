package controllers

import (
	stc "UsersAPI/structs"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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
		}
	}
}

func ManageProcessor(ctx context.Context, done chan<- []stc.User) {
	/* numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU) */
	numWorkers := 500
	totalUsers := 10000

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

	fmt.Println(users)
	done <- users
}
