package players

import (
	"math/rand"
	"time"

	"github.com/icrowley/fake"
)

var idCnt int = 1

type Player struct {
	ID      int
	Name    string
	Surname string
	Age     int
	Number  int
	TeamID  int
}

var Players []Player

func Generate(teamID int) {
	seed := time.Now().UTC().UnixNano()
	fake.Seed(seed)
	for i := 0; i < 20; i++ {
		p := Player{
			ID:      idCnt,
			Name:    fake.MaleFirstName(),
			Surname: fake.MaleLastName(),
			Age:     rand.Intn(23) + 18,
			Number:  rand.Intn(100),
			TeamID:  teamID,
		}
		idCnt++

		Players = append(Players, p)
	}
}
