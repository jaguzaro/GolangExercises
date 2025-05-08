package main

import (
	ctr "UsersAPI/controllers"
	stc "UsersAPI/structs"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/get-users", GetUsersHandler)
	r.Run(":8080")

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
