//go:debug x509negativeserial=1
package main

import (
	ctr "UsersAPI/controllers"
	dtb "UsersAPI/db"
	"UsersAPI/db/models"
	stc "UsersAPI/structs"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	mainctx, maincancel := context.WithCancel(context.Background())
	defer maincancel()

	appSignal := make(chan os.Signal, 1)
	signal.Notify(appSignal, os.Interrupt)
	go func() {
		<-appSignal
		maincancel()
	}()

	go func() {
		db, err := dtb.CreateConnection(mainctx)
		if err != nil {
			fmt.Println(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println(err)
		}

		sqlDB.SetMaxOpenConns(500)
		sqlDB.SetMaxIdleConns(100)
		sqlDB.SetConnMaxLifetime(1 * time.Minute)
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: createServer(),
	}

	go func() {
		fmt.Println("Servidor iniciado en el puerto :8080")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Error al iniciar el servidor")
		}
	}()
	<-mainctx.Done()

}

func GetUsersHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	users, err := CallGetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"users": users})

}

func CallGetUsers(ctx context.Context) ([]stc.User, error) {
	done := make(chan []stc.User)
	go ctr.ManageProcessor(ctx, done)
	select {
	case users := <-done:
		fmt.Println("Usuarios recolectados", len(users))
		return users, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("Timeout")
	}
}

func CallQueryDB(ctx context.Context) ([]models.Users, error) {
	done := make(chan []models.Users)
	go ctr.ManageQueries(ctx, done)
	select {
	case users := <-done:
		return users, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout")
	}
}

func QueryDBHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	users, err := CallQueryDB(ctx)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func createServer() *gin.Engine {
	r := gin.Default()
	r.GET("/get-users", GetUsersHandler)
	r.GET("/query-db", QueryDBHandler)
	return r
}
