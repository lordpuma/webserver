package Types

import (
	"database/sql"
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

type User struct {
	Id        int
	Name      string
	Username  string
	Email     string
	FirstName string
	LastName  string
	ShortName string
	Color     string
	BgColor   string
}

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "Basic User Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).Id, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).Name, nil
			},
		},
		"username": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).Username, nil
			},
		},
		"email": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).Email, nil
			},
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).FirstName, nil
			},
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).LastName, nil
			},
		},
		"shortName": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).ShortName, nil
			},
		},
		"bgColor": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).BgColor, nil
			},
		},
		"color": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(User).Color, nil
			},
		},
		"perms": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var perm string
				var perms []string
				rows, err := database.Db.Query("select perm from perms where user_id = ?", p.Source.(User).Id)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&perm)
					if err != nil {
						log.Fatal(err)
					}
					perms = append(perms, perm)
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}

				return perms, nil
			},
		},
		"workplaces": &graphql.Field{
			Type: graphql.NewList(WorkplaceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var idw int32
				var workplaces []Workplace
				rows, err := database.Db.Query("select id from workplaces where id IN (select workplace_id from users_workplaces where user_id = ?)", p.Source.(User).Id)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&idw)
					if err != nil {
						log.Fatal(err)
					}
					workplaces = append(workplaces, LoadWorkplaceById(int(idw)))
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				return workplaces, nil
			},
		},
	},
})

func LoadUserById(id int) User {
	var (
		username   string
		email      sql.NullString
		first_name string
		last_name  string
		bg_color   string
		color      string
	)
	err := database.Db.QueryRow("SELECT username, email, first_name, last_name, bg_color, color FROM users WHERE id = ?", id).
		Scan(&username, &email, &first_name, &last_name, &bg_color, &color)

	if err != nil {
		panic(err)
	}

	return User{
		Id:        id,
		Name:      first_name + " " + last_name,
		Username:  username,
		ShortName: string([]rune(first_name)[0]) + ". " + last_name,
		Color:     color,
		BgColor:   bg_color,
		LastName:  last_name,
		FirstName: first_name,
		Email:     email.String,
	}
}
