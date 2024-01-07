package main

import (
	"errors"
	"football/players"
	"football/team"
	"football/user"
	"strings"

	"github.com/gin-gonic/gin"
)

type apiResponse struct {
	Code    int
	Status  string
	Payload interface{}
}

func main() {
	r := gin.Default()
	r.GET("/users", GetUsersHandler)
	r.GET("/teams", GetTeamsHandler)
	r.GET("/players", GetPlayersHandler)
	r.POST("/register", UserRegisterHandler)
	r.Run(":8080")
}

func GetUsersHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: user.Users})
}

func GetTeamsHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: team.Teams})
}

func GetPlayersHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: players.Players})
}

func UserRegisterHandler(c *gin.Context) {
	var u user.User
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(400, apiResponse{Payload: err.Error()})
		return
	}

	err = validateRegister(&u)
	if err != nil {
		c.JSON(400, apiResponse{Payload: err.Error()})
		return
	}

	err = user.Register(&u)
	if err != nil {
		c.JSON(400, apiResponse{Payload: err.Error()})
		return
	}

	t := team.Team{Name: u.Nickname, UserID: u.ID}
	team.Register(&t)

	players.Generate(t.ID)

	c.JSON(200, apiResponse{Payload: "ok"})
}

func validateRegister(u *user.User) error {
	u.Login = strings.TrimSpace(u.Login)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Password = strings.TrimSpace(u.Password)

	if u.Login == "" {
		return errors.New("no login provided")
	}

	if u.Password == "" {
		return errors.New("no password provided")
	}

	if u.Nickname == "" {
		return errors.New("no nickname provided")
	}

	if len(strings.Split(u.Login, " ")) > 1 {
		return errors.New("login should be one word")
	}

	return nil
}

// 1. register
// 2. create team
// 3. create 20 random players and assign to the team
