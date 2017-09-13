package Types

import (
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

type Race struct {
	Id     int
	Name   string
	Active bool
}

var RaceType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Race",
	Description: "Basic Race Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Race).Id, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Race).Name, nil
			},
		},
		"active": &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Race).Active, nil
			},
		},
	},
})

func LoadRaceById(id int) Race {
	var (
		name   string
		active bool
	)
	err := database.Db.QueryRow("SELECT name, active FROM races WHERE id = ?", id).Scan(&name, &active)
	if err != nil {
		panic(err)
	}

	return Race{
		Id:     id,
		Name:   name,
		Active: active,
	}
}

func LoadRacesList() []Race {
	var r []Race
	var (
		id     int
		name   string
		active bool
	)
	rows, err := database.Db.Query("select id, name, active from races")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &active)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, Race{
			Id:     id,
			Name:   name,
			Active: active,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return r
}
