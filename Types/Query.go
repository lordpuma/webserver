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
	"race": &graphql.Field{
		Type: RaceType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadRaceById(p.Args["Id"].(int)), nil
		},
	},
	"races": &graphql.Field{
		Type: graphql.NewList(RaceType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadRacesList(), nil
		},
	},
	"result": &graphql.Field{
		Type: ResultType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadResultById(p.Args["Id"].(int)), nil
		},
	},
	"raceResults": &graphql.Field{
		Type: graphql.NewList(FormattedResultType),
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadResultsByRace(p.Args["Id"].(int)), nil
		},
	},
	"results": &graphql.Field{
		Type: graphql.NewList(ResultType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return LoadResultsList(), nil
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
		Type: graphql.NewList(DayType),
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
				d            string
				note         string
			)
			var days []Day
			date, isDOK := p.Args["Date"].(string)
			if isDOK {
				rows, err := database.Db.Query("SELECT id, workplace_id, user_id, date, note, DATE_FORMAT(date, '%e') AS day FROM shifts WHERE DATE_FORMAT(date, '%Y-%m') = ? ORDER BY day, workplace_id", date)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&id, &workplace_id, &user_id, &d, &note, &day)
					if err != nil {
						log.Fatal(err)
					}
					var needle *Day
					for k, v := range days {
						if v.Day == day {
							needle = &days[k]
						}
					}
					if needle == nil {
						days = append(days, Day{day, []W{{workplace_id, []Shift{{Id: id, Date: d, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)}}}}})
					} else {
						fmt.Println(needle)
						var n *W
						for ke, ve := range needle.Workplace {
							if ve.Id == workplace_id {
								n = &needle.Workplace[ke]
							}
						}
						if n == nil {
							*needle = Day{needle.Day,
								append(needle.Workplace,
									W{workplace_id, []Shift{{Id: id, Date: d, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)}}})}
						} else {
							*n = W{n.Id, append(n.Shifts, Shift{Id: id, Date: d, Note: note, user_id: int(user_id), workplace_id: int(workplace_id)})}
						}
					}

				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}

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

type W struct {
	Id     int
	Shifts []Shift
}

type Day struct {
	Day       int
	Workplace []W
}

var WType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "WType",
	Description: "Basic Workplace Object",
	Fields: graphql.Fields{
		"Id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(W).Id, nil
			},
		},
		"Shifts": &graphql.Field{
			Type: graphql.NewList(ShiftType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(W).Shifts, nil
			},
		},
	},
})

var DayType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "DayType",
	Description: "Basic Workplace Object",
	Fields: graphql.Fields{
		"Day": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Day).Day, nil
			},
		},
		"Workplaces": &graphql.Field{
			Type: graphql.NewList(WType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Day).Workplace, nil
			},
		},
	},
})
