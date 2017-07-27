package Types

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

var RootQuery = graphql.Fields{
	"thisUser": &graphql.Field{
		Type: UserType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadUserById(p.Context.Value("user_id").(int)), nil
		},
	},
	"user": &graphql.Field{
		Type: UserType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadUserById(p.Args["Id"].(int)), nil
		},
	},
	"users": &graphql.Field{
		Type: graphql.NewList(UserType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadUsersList(), nil
		},
	},
	"workplace": &graphql.Field{
		Type: WorkplaceType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadWorkplaceById(p.Args["Id"].(int)), nil
		},
	},
	"workplaces": &graphql.Field{
		Type: graphql.NewList(WorkplaceType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadWorkplacesList(), nil
		},
	},
	"shift": &graphql.Field{
		Type: ShiftType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadShiftById(p.Args["Id"].(int)), nil
		},
	},

	"shifts": &graphql.Field{
		Type: graphql.NewList(ShiftType),
		Args: graphql.FieldConfigArgument{
			"Date": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"Workplace": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"User": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			workplace, isWOK := p.Args["Workplace"].(int)
			date, isDOK := p.Args["Date"].(string)
			user, isUOK := p.Args["User"].(int)
			if isDOK && isWOK {
				return LoadShiftsByWorkplace(date, workplace), nil
			}
			if isUOK && isWOK {
				return LoadShiftsByUser(date, user), nil
			}
			return nil, nil
		},
	},

	"allShifts": &graphql.Field{
		Type: graphql.NewList(graphql.NewList(AllShiftsType)),
		Args: graphql.FieldConfigArgument{
			"Date": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var (
				id           int
				workplace_id int
				user_id      int
				day          int
			)
			var days = make(map[int]map[int][]Res)

			date, isDOK := p.Args["Date"].(string)

			if isDOK {
				rows, err := database.Db.Query("SELECT id, user_id, workplace_id, DATE_FORMAT(date, '%e') AS day FROM shifts WHERE DATE_FORMAT(date, '%Y-%m') = ? ORDER BY day, workplace_id", date)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&id, &user_id, &workplace_id, &day)
					if err != nil {
						log.Fatal(err)
					}
					//if days[day] == nil {
					//	log.Printf("day %d is nil", day)
					//	days[day] = make(map[int][]Res)
					//}
					//if days[day][workplace_id] == nil {
					//	log.Printf("workplace %d is nil", workplace_id)
					//	days[day][workplace_id] = []Res
					//}
					days[day][workplace_id] = append(days[day][workplace_id], Res{user_id, id})
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(days)
				return days, nil
			}
			return nil, nil
		},
	},
	"freeUsers": &graphql.Field{
		Type: graphql.NewList(UserType),
		Args: graphql.FieldConfigArgument{
			"Date": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"Workplace": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return FreeUsers(p.Args["Date"].(string), p.Args["Workplace"].(int)), nil
		},
	},
}

type Res struct {
	User_id int
	Id      int
}

var AllShiftsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "AllShiftsType",
	Description: "Basic Workplace Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Res).Id, nil
			},
		},
		"user_id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Res).User_id, nil
			},
		},
	},
})
