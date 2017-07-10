package Types

import (
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

type Workplace struct {
	Id      int
	Name    string
	BgColor string
	Color   string
}

var WorkplaceType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Workplace",
	Description: "Basic Workplace Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Workplace).Id, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Workplace).Name, nil
			},
		},
		"bgColor": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Workplace).BgColor, nil
			},
		},
		"color": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Workplace).Color, nil
			},
		},
		//"users": &graphql.Field{
		//	Type: graphql.NewList(UserType),
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		var users []User
		//		var user_id int32
		//		rows, err := database.Db.Query("SELECT user_id FROM users_workplaces WHERE workplace_id = ?", p.Source.(Workplace).Id)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		defer rows.Close()
		//		for rows.Next() {
		//			err := rows.Scan(&user_id)
		//			if err != nil {
		//				log.Fatal(err)
		//			}
		//			users = append(users, LoadUserById(int(user_id)))
		//		}
		//		err = rows.Err()
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		return users, nil
		//	},
		//},
	},
})

func LoadWorkplaceById(id int) Workplace {
	var (
		name     string
		bg_color string
		color    string
	)
	err := database.Db.QueryRow("SELECT name, bg_color, color FROM workplaces WHERE id = ?", id).Scan(&name, &bg_color, &color)
	if err != nil {
		panic(err)
	}

	return Workplace{
		Id:      id,
		Name:    name,
		Color:   color,
		BgColor: bg_color,
	}
}

func LoadWorkplacesList() []Workplace {
	var r []Workplace
	var (
		id       int
		name     string
		bg_color string
		color    string
	)
	rows, err := database.Db.Query("select id, name, bg_color, color from workplaces")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &bg_color, &color)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, Workplace{
			Id:      id,
			Name:    name,
			Color:   color,
			BgColor: bg_color,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return r
}
