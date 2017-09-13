package Types

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"io"
	"strings"
)

var RootMutation = graphql.Fields{
	"editUser": &graphql.Field{
		Type: UserType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"Username": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"FirstName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"LastName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"BgColor": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Color": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Perms": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
			"Workplaces": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Username"] != nil {
				_, err := database.Db.Exec("UPDATE users SET username = ? WHERE id = ?", p.Args["Username"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["FirstName"] != nil {
				_, err := database.Db.Exec("UPDATE users SET first_name = ? WHERE id = ?", p.Args["FirstName"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["LastName"] != nil {
				_, err := database.Db.Exec("UPDATE users SET last_name = ? WHERE id = ?", p.Args["LastName"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["BgColor"] != nil {
				_, err := database.Db.Exec("UPDATE users SET bg_color = ? WHERE id = ?", p.Args["BgColor"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Color"] != nil {
				_, err := database.Db.Exec("UPDATE users SET color = ? WHERE id = ?", p.Args["Color"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Email"] != nil {
				_, err := database.Db.Exec("UPDATE users SET email = ? WHERE id = ?", p.Args["Email"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Password"] != nil {
				if p.Args["Password"] == "1" {
					_, err := database.Db.Exec("UPDATE users SET pass = NULL WHERE id = ?", p.Args["Id"].(int))
					if err != nil {
						return nil, err
					}
				} else {
					h := md5.New()
					io.WriteString(h, p.Args["Password"].(string))
					pass := new(bytes.Buffer)
					fmt.Fprintf(pass, "%x", h.Sum(nil)) //cast pass hash to var
					_, err := database.Db.Exec("UPDATE users SET pass = ? WHERE id = ?", pass.String(), p.Args["Id"].(int))
					if err != nil {
						return nil, err
					}
				}
			}
			if p.Args["Perms"] != nil {
				database.Db.Exec("DELETE FROM perms WHERE user_id = ?", p.Args["Id"].(int))
				for _, value := range p.Args["Perms"].([]interface{}) {
					database.Db.Exec("INSERT INTO perms (user_id, perm) VALUES (?, ?)", p.Args["Id"].(int), value)
				}
			}
			if p.Args["Workplaces"] != nil {
				database.Db.Exec("DELETE FROM users_workplaces WHERE user_id = ?", p.Args["Id"].(int))
				for _, value := range p.Args["Workplaces"].([]interface{}) {
					database.Db.Exec("INSERT INTO users_workplaces (user_id, workplace_id) VALUES (?, ?)", p.Args["Id"].(int), value)
				}
			}
			return LoadUserById(p.Args["Id"].(int)), nil
		},
	},
	"insertUser": &graphql.Field{
		Type: UserType,
		Args: graphql.FieldConfigArgument{
			"Username": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"FirstName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"LastName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"BgColor": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Color": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Perms": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
			"Workplaces": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var email string
			var (
				out sql.Result
				err error
			)
			username := strings.ToLower(p.Args["Username"].(string))
			if p.Args["Email"] == nil {
				out, err = database.Db.Exec("INSERT INTO users (username, first_name, last_name, bg_color, color, email) VALUES (?, ?, ?, ?, ?, ?)", username, p.Args["FirstName"].(string), p.Args["LastName"].(string), p.Args["BgColor"].(string), p.Args["Color"].(string), nil)
			} else {
				email = p.Args["Email"].(string)
				out, err = database.Db.Exec("INSERT INTO users (username, first_name, last_name, bg_color, color, email) VALUES (?, ?, ?, ?, ?, ?)", username, p.Args["FirstName"].(string), p.Args["LastName"].(string), p.Args["BgColor"].(string), p.Args["Color"].(string), email)
			}
			if err != nil {
				return nil, err
			}
			id, er := out.LastInsertId()
			if er != nil {
				return nil, err
			}
			if p.Args["Perms"] != nil {
				for _, value := range p.Args["Perms"].([]interface{}) {
					database.Db.Exec("INSERT INTO perms (user_id, perm) VALUES (?, ?)", id, value)
					if err != nil {
						return nil, err
					}
				}
			}
			if p.Args["Workplaces"] != nil {

				for _, value := range p.Args["Workplaces"].([]interface{}) {
					database.Db.Exec("INSERT INTO users_workplaces (user_id, workplace_id) VALUES (?, ?)", id, value)
					if err != nil {
						return nil, err
					}
				}
			}
			return LoadUserById(int(id)), nil
		},
	},
	"editWorkplace": &graphql.Field{
		Type: WorkplaceType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"Name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"BgColor": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Color": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Name"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET name = ? WHERE id = ?", p.Args["Name"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["BgColor"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET bg_color = ? WHERE id = ?", p.Args["BgColor"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Color"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET color = ? WHERE id = ?", p.Args["Color"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			return LoadWorkplaceById(p.Args["Id"].(int)), nil
		},
	},
	"insertWorkplace": &graphql.Field{
		Type: WorkplaceType,
		Args: graphql.FieldConfigArgument{
			"Name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"BgColor": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Color": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			out, err := database.Db.Exec("INSERT INTO workplaces (name, bg_color, color) VALUES (?, ?, ?)", p.Args["Name"].(string), p.Args["BgColor"].(string), p.Args["Color"].(string))
			if err != nil {
				return nil, err
			}
			id, er := out.LastInsertId()
			if er != nil {
				return nil, err
			}
			return LoadWorkplaceById(int(id)), nil
		},
	},
	"editRace": &graphql.Field{
		Type: RaceType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"Name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Active": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Name"] != nil {
				_, err := database.Db.Exec("UPDATE races SET name = ? WHERE id = ?", p.Args["Name"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Active"] != nil {
				_, err := database.Db.Exec("UPDATE races SET active = ? WHERE id = ?", p.Args["Active"].(string), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			return LoadRaceById(p.Args["Id"].(int)), nil
		},
	},
	"insertRace": &graphql.Field{
		Type: RaceType,
		Args: graphql.FieldConfigArgument{
			"Name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Active": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Active"].(bool) {
				_, err := database.Db.Exec("UPDATE races SET active = FALSE")
				if err != nil {
					return nil, err
				}
			}
			out, err := database.Db.Exec("INSERT INTO races (name, active) VALUES (?, ?)", p.Args["Name"].(string), p.Args["Active"].(bool))
			if err != nil {
				return nil, err
			}
			id, er := out.LastInsertId()
			if er != nil {
				return nil, err
			}
			return LoadRaceById(int(id)), nil
		},
	},
	"insertShift": &graphql.Field{
		Type: ShiftType,
		Args: graphql.FieldConfigArgument{
			"Userid": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"Workplaceid": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"Note": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Date": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Note"] == nil {
				p.Args["Note"] = ""
			}
			out, err := database.Db.Exec("INSERT INTO shifts (user_id, workplace_id, date, note) VALUES (?, ?, ?, ?)", p.Args["Userid"].(int), p.Args["Workplaceid"].(int), p.Args["Date"], p.Args["Note"])
			if err != nil {
				return nil, err
			}
			id, _ := out.LastInsertId()
			return LoadShiftById(int(id)), nil
		},
	},
	"editShift": &graphql.Field{
		Type: ShiftType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"Userid": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"Workplaceid": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"Note": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"Date": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if p.Args["Userid"] != nil {
				_, err := database.Db.Exec("UPDATE shifts SET user_id = ? WHERE id = ?", p.Args["Userid"].(int), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Workplaceid"] != nil {
				_, err := database.Db.Exec("UPDATE shifts SET workplace_id = ? WHERE id = ?", p.Args["Workplaceid"].(int), p.Args["Id"].(int))
				if err != nil {
					return nil, err
				}
			}
			if p.Args["Note"] != nil {
				if p.Args["Note"] == "RESET NOTE PLS" {
					_, err := database.Db.Exec("UPDATE shifts SET note = ? WHERE id = ?", "", p.Args["Id"].(int))
					if err != nil {
						return nil, err
					}
				} else {
					_, err := database.Db.Exec("UPDATE shifts SET note = ? WHERE id = ?", p.Args["Note"].(string), p.Args["Id"].(int))
					if err != nil {
						return nil, err
					}
				}
			}
			return LoadShiftById(p.Args["Id"].(int)), nil
		},
	},
	"deleteShift": &graphql.Field{
		Type: graphql.Int,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			out, err := database.Db.Exec("DELETE FROM shifts WHERE id = ?", p.Args["Id"].(int))
			if err != nil {
				return nil, err
			}
			id, _ := out.LastInsertId()
			return id, nil
		},
	},
	"deleteResult": &graphql.Field{
		Type: graphql.Int,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			out, err := database.Db.Exec("DELETE FROM results WHERE id = ?", p.Args["Id"].(int))
			if err != nil {
				return nil, err
			}
			id, _ := out.LastInsertId()
			return id, nil
		},
	},
	"setActiveRace": &graphql.Field{
		Type: RaceType,
		Args: graphql.FieldConfigArgument{
			"Id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"Active": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			_, err := database.Db.Exec("UPDATE races SET active = FALSE")
			if err != nil {
				return nil, err
			}
			_, err = database.Db.Exec("UPDATE races SET active = ? WHERE id = ?", p.Args["Active"].(bool), p.Args["Id"].(int))
			if err != nil {
				return nil, err
			}
			return LoadRaceById(p.Args["Id"].(int)), nil
		},
	},
}
