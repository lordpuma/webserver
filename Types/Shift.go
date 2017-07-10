package Types

import (
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

type Shift struct {
	Id           int
	Date         string
	Note         string
	user_id      int
	workplace_id int
}

var ShiftType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Shift",
	Description: "Basic Shift Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Shift).Id, nil
			},
		},
		"user": &graphql.Field{
			Type: UserType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return LoadUserById(p.Source.(Shift).user_id), nil
			},
		},
		"workplace": &graphql.Field{
			Type: WorkplaceType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return LoadWorkplaceById(p.Source.(Shift).workplace_id), nil
			},
		},
		"date": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Shift).Date, nil
			},
		},
		"note": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Shift).Note, nil
			},
		},
	},
})

func LoadShiftById(id int) Shift {
	var (
		workplace_id int32
		user_id      int32
		date         string
		note         string
	)
	err := database.Db.QueryRow("SELECT workplace_id, user_id, date, note FROM shifts WHERE id = ?", id).Scan(&workplace_id, &user_id, &date, &note)
	if err != nil {
		panic(err)
	}
	return Shift{Id: id, Date: date, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)}
}

func LoadShiftsByWorkplace(date string, workplace int) []Shift {
	var r []Shift
	var (
		id           int
		workplace_id int
		user_id      int
		note         string
	)
	rows, err := database.Db.Query("SELECT id, workplace_id, user_id, note FROM shifts WHERE date = ? AND workplace_id = ?", date, workplace)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &workplace_id, &user_id, &note)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, Shift{Id: id, Date: date, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func LoadShiftsByUser(date string, user int) Shift {
	var (
		id           int
		workplace_id int
		user_id      int
		note         string
	)
	err := database.Db.QueryRow("SELECT id, workplace_id, user_id, note FROM shifts WHERE date = ? AND user_id = ?", date, user).Scan(&id, &workplace_id, &user_id, &note)
	if err != nil {
		return Shift{Id: 0}
	}
	return Shift{Id: id, Date: date, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)}
}
