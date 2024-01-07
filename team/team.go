package team

var Teams []Team
var idCnt int = 1

type Team struct {
	ID     int
	UserID int
	Name   string
}

func Register(t *Team) {
	t.ID = idCnt
	idCnt++
	Teams = append(Teams, *t)
}
