// Package starwars provides a example schema and resolver based on Star Wars characters.
//
// Source: https://github.com/graphql/graphql.github.io/blob/source/site/_core/swapiSchema.js
package resolvers

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/lordpuma/webserver/database"
	"github.com/playlyfe/go-graphql"
	"io"
	"log"
	"strings"
)

var Schema = `
	{
		query: Query
		mutation: Mutation
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		thisUser: User!
		user(Id: Int): User!
		users: [User]
		shifts(Date: String!, Workplace: Int, User: Int): [Shift]
		shift(Id: Int!): Shift!
		workplace(Id: Int!): Workplace!
		workplaces: [Workplace]
		freeUsers(Date: String!, Workplace: Int!): [User]
	}
	# The mutation type, represents all updates we can make to our data
	type Mutation {
		insertShift(Date: String!, Userid: Int!, Workplaceid: Int!,  Note: String): Shift!
		insertUser(Username: String!, FirstName: String!, LastName: String!, BgColor: String!, Color: String!, Email: String, Perms: [String], Workplaces: [Int]): User!
		insertWorkplace(Name: String!, BgColor: String!, Color: String!): Workplace!
		editShift(Id: Int!, Userid: Int, Workplaceid: Int,  Note: String): Shift!
		deleteShift(Id: Int!): Int!
		editWorkplace(Id: Int!, Name: String, BgColor: String, Color: String): Workplace!
		editUser(Id: Int!, Username: String, FirstName: String, LastName: String, BgColor: String, Color: String, Email: String , Password: String, Workplaces: [Int], Perms: [String]): User!
	}

	type User {
		id: Int!
		name: String!
		email: String
		username: String!
		firstName: String!
		lastName: String!
		shortName: String!
		perms: [String]
		bgColor: String!
		color: String!
		workplaces: [Workplace]
	}

	type Shift {
		id: Int!
		user: User!
		workplace: Workplace!
		date: String!
		note: String
	}

	type Workplace {
		id: Int!
		name: String!
		bgColor: String!
		color: String!
	}
`

func User(id int32) map[string]interface{} {
	var (
		username   string
		email      sql.NullString
		first_name string
		last_name  string
		perms      []string
		bg_color   string
		color      string
		workplaces []map[string]interface{}
	)
	err := database.Db.QueryRow("SELECT username, email, first_name, last_name, bg_color, color FROM users WHERE id = ?", id).Scan(&username, &email, &first_name, &last_name, &bg_color, &color)
	if err != nil {
		panic(err)
	}

	var perm string
	rows, err := database.Db.Query("select perm from perms where user_id = ?", id)
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

	if err != nil {
		panic(err)
	}

	var idw int32
	rows, err = database.Db.Query("select id from workplaces where id IN (select workplace_id from users_workplaces where user_id = ?)", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&idw)
		if err != nil {
			log.Fatal(err)
		}
		workplaces = append(workplaces, Workplace(idw))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"id":         id,
		"username":   username,
		"email":      email.String,
		"name":       first_name + " " + last_name,
		"firstName":  first_name,
		"lastName":   last_name,
		"shortName":  string([]rune(first_name)[0]) + ". " + last_name,
		"perms":      perms,
		"bgColor":    bg_color,
		"color":      color,
		"workplaces": workplaces,
	}
}

// SHIFT

func Shift(id int32) map[string]interface{} {
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
	return map[string]interface{}{
		"id":        id,
		"workplace": Workplace(workplace_id),
		"user":      User(user_id),
		"date":      date,
		"note":      note,
	}
}

//WORKPLACE
func Workplace(id int32) map[string]interface{} {
	var (
		name     string
		bg_color string
		color    string
		//users    []map[string]interface{}
	)
	err := database.Db.QueryRow("SELECT name, bg_color, color FROM workplaces WHERE id = ?", id).Scan(&name, &bg_color, &color)
	if err != nil {
		panic(err)
	}

	//var user_id int32
	//rows, err := database.Db.Query("SELECT user_id FROM users_workplaces WHERE workplace_id = ?", id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	err := rows.Scan(&user_id)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	users = append(users, User(user_id))
	//}
	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}

	return map[string]interface{}{
		"id":      id,
		"name":    name,
		"bgColor": bg_color,
		"color":   color,
		//"users":   users,
	}
}

func GetResolvers() map[string]interface{} {
	var resolvers = map[string]interface{}{
		"Query/user": func(params *graphql.ResolveParams) (interface{}, error) {
			var id int32 = params.Args["Id"].(int32)
			return User(id), nil
		},
		"Query/thisUser": func(params *graphql.ResolveParams) (interface{}, error) {
			var id int32 = params.Context.(map[string]interface{})["user_id"].(int32)
			return User(id), nil
		},
		"Query/users": func(params *graphql.ResolveParams) (interface{}, error) {
			var r []interface{}
			var id int32
			rows, err := database.Db.Query("select id from users")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					log.Fatal(err)
				}
				r = append(r, User(id))
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			return r, nil
		},
		"Query/freeUsers": func(params *graphql.ResolveParams) (interface{}, error) {
			var r []interface{}
			var id int32
			rows, err := database.Db.Query("SELECT id from users WHERE id NOT IN (SELECT user_id FROM vacations WHERE date = ?) AND id IN (SELECT user_id FROM users_workplaces WHERE workplace_id = ?) AND id NOT IN (SELECT user_id FROM shifts WHERE date = ? AND workplace_id = ?)", params.Args["Date"].(string), params.Args["Workplace"].(int32), params.Args["Date"].(string), params.Args["Workplace"].(int32))
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					log.Fatal(err)
				}
				r = append(r, User(id))
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			return r, nil
		},
		"Mutation/editUser": func(params *graphql.ResolveParams) (interface{}, error) {
			if params.Args["Username"] != nil {
				_, err := database.Db.Exec("UPDATE users SET username = ? WHERE id = ?", params.Args["Username"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["FirstName"] != nil {
				_, err := database.Db.Exec("UPDATE users SET first_name = ? WHERE id = ?", params.Args["FirstName"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["LastName"] != nil {
				_, err := database.Db.Exec("UPDATE users SET last_name = ? WHERE id = ?", params.Args["LastName"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["BgColor"] != nil {
				_, err := database.Db.Exec("UPDATE users SET bg_color = ? WHERE id = ?", params.Args["BgColor"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Color"] != nil {
				_, err := database.Db.Exec("UPDATE users SET color = ? WHERE id = ?", params.Args["Color"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Email"] != nil {
				_, err := database.Db.Exec("UPDATE users SET email = ? WHERE id = ?", params.Args["Email"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Password"] != nil {
				if params.Args["Password"] == "1" {
					_, err := database.Db.Exec("UPDATE users SET pass = NULL WHERE id = ?", params.Args["Id"].(int32))
					if err != nil {
						return nil, err
					}
				} else {
					h := md5.New()
					io.WriteString(h, params.Args["Password"].(string))
					pass := new(bytes.Buffer)
					fmt.Fprintf(pass, "%x", h.Sum(nil)) //cast pass hash to var
					_, err := database.Db.Exec("UPDATE users SET pass = ? WHERE id = ?", pass.String(), params.Args["Id"].(int32))
					if err != nil {
						return nil, err
					}
				}
			}
			if params.Args["Perms"] != nil {
				database.Db.Exec("DELETE FROM perms WHERE user_id = ?", params.Args["Id"].(int32))
				for key, value := range params.Args["Perms"].([]interface{}) {
					fmt.Println(key, value)
					database.Db.Exec("INSERT INTO perms (user_id, perm) VALUES (?, ?)", params.Args["Id"].(int32), value)
				}
			}
			if params.Args["Workplaces"] != nil {
				database.Db.Exec("DELETE FROM users_workplaces WHERE user_id = ?", params.Args["Id"].(int32))
				for key, value := range params.Args["Workplaces"].([]interface{}) {
					fmt.Println(key, value)
					database.Db.Exec("INSERT INTO users_workplaces (user_id, workplace_id) VALUES (?, ?)", params.Args["Id"].(int32), value)
				}
			}
			fmt.Println(params.Args["Perms"])
			return User(params.Args["Id"].(int32)), nil
		},
		"Mutation/insertUser": func(params *graphql.ResolveParams) (interface{}, error) {
			//h := md5.New()
			//io.WriteString(h, params.Args["Password"].(string))
			//pass := new(bytes.Buffer)
			//fmt.Fprintf(pass, "%x", h.Sum(nil)) //cast pass hash to var
			var email string
			if params.Args["Email"] == nil {
				email = ""
			} else {
				email = params.Args["Email"].(string)
			}
			username := strings.ToLower(params.Args["Username"].(string))
			out, err := database.Db.Exec("INSERT INTO users (username, first_name, last_name, bg_color, color, email) VALUES (?, ?, ?, ?, ?, ?)", username, params.Args["FirstName"].(string), params.Args["LastName"].(string), params.Args["BgColor"].(string), params.Args["Color"].(string), email)
			if err != nil {
				return nil, err
			}
			id, er := out.LastInsertId()
			if er != nil {
				return nil, err
			}
			if params.Args["Perms"] != nil {
				for key, value := range params.Args["Perms"].([]interface{}) {
					fmt.Println(key, value)
					database.Db.Exec("INSERT INTO perms (user_id, perm) VALUES (?, ?)", id, value)
					if err != nil {
						return nil, err
					}
				}
			}
			fmt.Println(params.Args["Workplaces"])
			if params.Args["Workplaces"] != nil {

				for key, value := range params.Args["Workplaces"].([]interface{}) {
					fmt.Println(key, value)
					database.Db.Exec("INSERT INTO users_workplaces (user_id, workplace_id) VALUES (?, ?)", id, value)
					if err != nil {
						return nil, err
					}
				}
			}
			return User(int32(id)), nil
		},

		"Query/workplace": func(params *graphql.ResolveParams) (interface{}, error) {
			var id int32 = params.Args["Id"].(int32)
			return Workplace(id), nil
		},
		"Query/workplaces": func(params *graphql.ResolveParams) (interface{}, error) {
			var r []interface{}
			var id int32
			rows, err := database.Db.Query("select id from workplaces")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					log.Fatal(err)
				}
				r = append(r, Workplace(id))
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			return r, nil
		},
		"Mutation/editWorkplace": func(params *graphql.ResolveParams) (interface{}, error) {
			if params.Args["Name"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET name = ? WHERE id = ?", params.Args["Name"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["BgColor"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET bg_color = ? WHERE id = ?", params.Args["BgColor"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Color"] != nil {
				_, err := database.Db.Exec("UPDATE workplaces SET color = ? WHERE id = ?", params.Args["Color"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			return Workplace(params.Args["Id"].(int32)), nil
		},
		"Mutation/insertWorkplace": func(params *graphql.ResolveParams) (interface{}, error) {
			out, err := database.Db.Exec("INSERT INTO workplaces (name, bg_color, color) VALUES (?, ?, ?)", params.Args["Name"].(string), params.Args["BgColor"].(string), params.Args["Color"].(string))
			if err != nil {
				return nil, err
			}
			id, er := out.LastInsertId()
			if er != nil {
				return nil, err
			}
			return Workplace(int32(id)), nil
		},

		"Query/shifts": func(params *graphql.ResolveParams) (interface{}, error) {
			if params.Args["Workplace"] != nil {
				var r []interface{}
				var id int32
				rows, err := database.Db.Query("SELECT id FROM shifts WHERE date = ? AND workplace_id = ?", params.Args["Date"].(string), params.Args["Workplace"].(int32))
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&id)
					if err != nil {
						log.Fatal(err)
					}
					r = append(r, Shift(id))
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(r)
				return r, nil
			}
			if params.Args["User"] != nil {
				var id int32
				err := database.Db.QueryRow("SELECT id FROM shifts WHERE date = ? AND user_id = ?", params.Args["Date"].(string), params.Args["User"].(int32)).Scan(&id)
				if err != nil {
					return Shift(0), nil
				}
				return Shift(id), nil
			}
			return Shift(0), nil
		},
		"Query/shift": func(params *graphql.ResolveParams) (interface{}, error) {
			var id int32 = params.Args["Id"].(int32)
			return Shift(id), nil
		},
		"Mutation/insertShift": func(params *graphql.ResolveParams) (interface{}, error) {
			if params.Args["Note"] == nil {
				params.Args["Note"] = ' '
			}
			out, err := database.Db.Exec("INSERT INTO shifts (user_id, workplace_id, date, note) VALUES (?, ?, ?, ?)", params.Args["Userid"].(int32), params.Args["Workplaceid"].(int32), params.Args["Date"], params.Args["Note"])
			if err != nil {
				return nil, err
			}
			id, _ := out.LastInsertId()
			return Shift(int32(id)), nil
		},
		"Mutation/editShift": func(params *graphql.ResolveParams) (interface{}, error) {
			if params.Args["Userid"] != nil {
				_, err := database.Db.Exec("UPDATE shifts SET user_id = ? WHERE id = ?", params.Args["Userid"].(int32), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Workplaceid"] != nil {
				_, err := database.Db.Exec("UPDATE shifts SET workplace_id = ? WHERE id = ?", params.Args["Workplaceid"].(int32), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			if params.Args["Note"] != nil {
				_, err := database.Db.Exec("UPDATE shifts SET note = ? WHERE id = ?", params.Args["Note"].(string), params.Args["Id"].(int32))
				if err != nil {
					return nil, err
				}
			}
			return Shift(params.Args["Id"].(int32)), nil
		},
		"Mutation/deleteShift": func(params *graphql.ResolveParams) (interface{}, error) {
			out, err := database.Db.Exec("DELETE FROM shifts WHERE id = ?", params.Args["Id"].(int32))
			if err != nil {
				return nil, err
			}
			id, _ := out.LastInsertId()
			return id, nil
		},
	}
	return resolvers
}
