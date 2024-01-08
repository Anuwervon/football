package main

import (
	"errors"
	"fmt"
	"football/players"
	"football/team"
	"football/user"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type apiResponse struct {
	Code    int
	Status  string
	Payload interface{}
}

func main() {
	r := gin.Default()
	r.GET("/users", AuthMiddleware(), GetUsersHandler)
	r.GET("/teams", AuthMiddleware(), GetTeamsHandler)
	r.GET("/players", AuthMiddleware(), GetPlayersHandler)
	r.PUT("/teams", AuthMiddleware(), UpdateTeamsHandler)
	r.POST("/register", UserRegisterHandler)
	r.POST("/login", UserLoginHandler)
	r.Run(":8080")
}

func GetUsersHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: user.Users})
}

func UpdateTeamsHandler(c *gin.Context) {
	var t team.Team
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(400, apiResponse{Payload: err.Error()})
		return
	}

	p, _ := c.Get("user")
	u := p.(user.User)

	//baroi har yak elementi da druni team.Teams mo 2ta variable i,v mesozem
	//har yak krugda i,v harxelay baroi ki i,v znacheniyayi element da druni slice elemento
	for i, v := range team.Teams {
		if t.ID == v.ID && u.ID == v.UserID {
			team.Teams[i].Name = t.Name
			break
		}
	}
	c.JSON(200, apiResponse{Payload: team.Teams})
}

func GetTeamsHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: team.Teams})
}

func GetPlayersHandler(c *gin.Context) {
	c.JSON(200, apiResponse{Payload: players.Players})
}

func UserRegisterHandler(c *gin.Context) {
	var u user.User
	//daniyora ay body-i zapros mexonem da druni u(u-variable user.User) menavisem
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

func UserLoginHandler(c *gin.Context) {
	var u user.User
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(400, apiResponse{Payload: err.Error()})
		return
	}

	if ok := user.Login(u); !ok {
		c.JSON(http.StatusUnauthorized, apiResponse{Payload: "incorrect login or password"})
		return
	}

	token, err := generateJWTToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Payload: err.Error()})
		return
	}

	c.JSON(200, apiResponse{Payload: token})
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Strip "Bearer " prefix
		// tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method and return the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("jwtSecret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		x := token.Claims.(jwt.MapClaims)["sub"].(map[string]interface{})

		var u = user.User{
			Login:    x["login"].(string),
			Password: x["password"].(string),
		}
		user.FillUserFromLoginAndPassword(&u)

		// Set the user information in the context
		c.Set("user", u)
		c.Next()
	}
}

func generateJWTToken(u user.User) (string, error) {
	// Define the claims for the JWT token
	claims := jwt.MapClaims{
		"sub": u,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		"iat": time.Now().Unix(),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte("jwtSecret"))
}

// 1. register
// 2. create team
// 3. create 20 random players and assign to the team
