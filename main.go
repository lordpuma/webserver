package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lordpuma/webserver/database"
	//"github.com/lordpuma/webserver/resolvers"
	//"github.com/playlyfe/go-graphql"
	"github.com/graphql-go/graphql"
	_ "github.com/joho/godotenv/autoload"

	"github.com/lordpuma/webserver/Types"
	"github.com/rs/cors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//var executor *graphql.Executor

var schema graphql.Schema

func init() {
	var err error
	if err != nil {
		panic(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.Header.Get("token") != "" {
			err := database.Db.QueryRow("SELECT user_id FROM logins WHERE token = ?", r.Header.Get("token")).Scan(&id)
			if err != nil {
				//return shift(0), nil
			}
			if id == 0 {
				resp, _ := json.Marshal(map[string]interface{}{"error": "bad_token"})
				w.Write(resp)
			} else {
				ctx := context.WithValue(context.Background(), "user_id", id)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			resp, _ := json.Marshal(map[string]interface{}{"error": "no_header"})
			w.Write(resp)
		}
	})
}

//
//func test() http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte(fmt.Sprint(r.Context().Value("user_id"))))
//	})
//}

func queryHand() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v map[string]interface{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
		}
		var m interface{}
		errr := json.Unmarshal(body, &m)
		if errr != nil {
			fmt.Fprintf(w, "%s", errr)
		}

		q := m.(map[string]interface{})["query"].(string)

		if m.(map[string]interface{})["variables"] != nil {
			v = m.(map[string]interface{})["variables"].(map[string]interface{})
		}

		params := graphql.Params{Schema: schema, RequestString: q, Context: r.Context(), VariableValues: v}
		ret := graphql.Do(params)
		if len(ret.Errors) > 0 {
			log.Fatalf("failed to execute graphql operation, errors: %+v", ret.Errors)
		}

		resp, _ := json.Marshal(ret)
		fmt.Fprintf(w, "%s", resp)

	})
}

//func main1() {
//	//r := resolvers.GetResolvers()
//	var err error
//	//executor, err = graphql.NewExecutor(resolvers.Schema, "Query", "Mutation", r)
//
//	//db, err := sql.Open("mysql", "root:password@tcp(db:3306)/database")
//	db, err := sql.Open("mysql", "root:pass@/database")
//	if err != nil {
//		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
//	}
//	defer db.Close()
//	database.Connect(db)
//
//	//CORS STARTS HERE - DEV ONL	Y
//	mux := http.NewServeMux()
//
//	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write(page)
//	}))
//
//	mux.Handle("/test", authMiddleware(test()))
//	//mux.Handle("/query", authMiddleware(queryHand()))
//	//mux.Handle("/query", authMiddleware(queryHand()))
//	//mux.Handle("/query", queryHand())
//
//	mux.HandleFunc("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		body, err := ioutil.ReadAll(r.Body)
//		var m interface{}
//		err = json.Unmarshal(body, &m)
//
//		if m != nil {
//			if m.(map[string]interface{})["pass"] != nil && m.(map[string]interface{})["user"] != nil {
//				h := md5.New()
//				io.WriteString(h, m.(map[string]interface{})["pass"].(string))
//				var username string
//				var pass sql.NullString
//				var id []uint8
//				rows, err := database.Db.Query("select id, username, pass from users")
//				if err != nil {
//					log.Fatal(err)
//				}
//				defer rows.Close()
//				for rows.Next() {
//					err := rows.Scan(&id, &username, &pass)
//					if err != nil {
//						log.Fatal(err)
//					}
//					if strings.ToLower(m.(map[string]interface{})["user"].(string)) == strings.ToLower(username) {
//						token := randToken()
//						if pass.Valid {
//							passhash := new(bytes.Buffer)
//							fmt.Fprintf(passhash, "%x", h.Sum(nil)) //cast pass hash to var
//							if passhash.String() == pass.String {
//								_, err := database.Db.Exec("INSERT INTO logins (user_id, token) VALUES (?, ?)", id, token) //Save Token to db
//								if err != nil {
//									panic(err)
//								}
//								resp, err := json.Marshal(map[string]interface{}{"token": token})
//								w.Write(resp) //USER IS LOGIN, send him token
//								return
//							}
//						} else {
//							_, err := database.Db.Exec("INSERT INTO logins (user_id, token) VALUES (?, ?)", id, token) //Save Token to db
//							if err != nil {
//								panic(err)
//							}
//							resp, err := json.Marshal(map[string]interface{}{"token": token, "first": true})
//							w.Write(resp) //USER IS LOGIN, send him token
//							return
//						}
//
//					}
//				}
//				err = rows.Err()
//				resp, err := json.Marshal(map[string]interface{}{"error": "unknown-user-or-pass"})
//				w.Write(resp)
//				return
//			}
//		}
//		resp, err := json.Marshal(map[string]interface{}{"error": err})
//		w.Write(resp)
//	}))
//
//	mux.HandleFunc("/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.Header.Get("token") != "" {
//			_, err := database.Db.Exec("DELETE FROM logins WHERE token = ?", r.Header.Get("token")) //delete token
//			if err != nil {
//				panic(err)
//			}
//			resp, _ := json.Marshal(map[string]interface{}{"success": true})
//			w.Write(resp)
//		} else {
//			resp, _ := json.Marshal(map[string]interface{}{"error": "no_header"})
//			w.Write(resp)
//		}
//
//	}))
//
//	// cors.Default() setup the middleware with default options being
//	// all origins accepted with simple methods (GET, POST). See
//	// documentation below for more options.
//	handler := cors.AllowAll().Handler(mux)
//
//	//CORS ENDS HERE - DEV ONLY
//
//	//PRODUCTION START
//
//	//http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	//	w.Write(page)
//	//}))
//	//
//	//http.Handle("/query", &relay.Handler{Schema: schema})
//
//	//PRODUCTION END
//
//	//SEED DB
//	time.Sleep(5000 * time.Millisecond)
//	var c int
//	err = database.Db.QueryRow("select count(id) as c from users").Scan(&c)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if c == 0 {
//		_, err := database.Db.Exec("INSERT INTO users (username, first_name, last_name, bg_color, color, email) VALUES (?, ?, ?, ?, ?, ?)", "lordpuma", "Tomáš", "Korený", "#000000", "#FFFFFF", "")
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	log.Fatal(http.ListenAndServe(":8080", handler))
//
//}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.0.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				graphQLParams.variables = graphQLParams.variables ? JSON.parse(graphQLParams.variables) : null;
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}

			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func main() {
	db, err := sql.Open("mysql", os.Getenv("DB"))
	//db, err := sql.Open("mysql", "root:pass@/database")
	//db, err := sql.Open("mysql", "root:password@tcp(db:3306)/database")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	database.Connect(db)

	fields := graphql.Fields{
		"thisUser": &graphql.Field{
			Type: Types.UserType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadUserById(p.Context.Value("user_id").(int)), nil
			},
		},
		"user": &graphql.Field{
			Type: Types.UserType,
			Args: graphql.FieldConfigArgument{
				"Id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadUserById(p.Args["Id"].(int)), nil
			},
		},
		"users": &graphql.Field{
			Type: graphql.NewList(Types.UserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadUsersList(), nil
			},
		},
		"workplace": &graphql.Field{
			Type: Types.WorkplaceType,
			Args: graphql.FieldConfigArgument{
				"Id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadWorkplaceById(p.Args["Id"].(int)), nil
			},
		},
		"workplaces": &graphql.Field{
			Type: graphql.NewList(Types.WorkplaceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadWorkplacesList(), nil
			},
		},
		"shift": &graphql.Field{
			Type: Types.ShiftType,
			Args: graphql.FieldConfigArgument{
				"Id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadShiftById(p.Args["Id"].(int)), nil
			},
		},
		"shifts": &graphql.Field{
			Type: graphql.NewList(Types.ShiftType),
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
					return Types.LoadShiftsByWorkplace(date, workplace), nil
				}
				if isUOK && isWOK {
					return Types.LoadShiftsByUser(date, user), nil
				}
				return nil, nil
			},
		},
		"freeUsers": &graphql.Field{
			Type: graphql.NewList(Types.UserType),
			Args: graphql.FieldConfigArgument{
				"Date": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"Workplace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.FreeUsers(p.Args["Date"].(string), p.Args["Workplace"].(int)), nil
			},
		},
	}

	mutations := graphql.Fields{
		"editUser": &graphql.Field{
			Type: Types.UserType,
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
				return Types.LoadUserById(p.Args["Id"].(int)), nil
			},
		},
		"insertUser": &graphql.Field{
			Type: Types.UserType,
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
					Type: graphql.NewList(Types.WorkplaceType),
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
				return Types.LoadUserById(int(id)), nil
			},
		},
		"editWorkplace": &graphql.Field{
			Type: Types.WorkplaceType,
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
				return Types.LoadWorkplaceById(p.Args["Id"].(int)), nil
			},
		},
		"insertWorkplace": &graphql.Field{
			Type: Types.WorkplaceType,
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
				return Types.LoadWorkplaceById(int(id)), nil
			},
		},
		"insertShift": &graphql.Field{
			Type: Types.ShiftType,
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
					p.Args["Note"] = ' '
				}
				out, err := database.Db.Exec("INSERT INTO shifts (user_id, workplace_id, date, note) VALUES (?, ?, ?, ?)", p.Args["Userid"].(int), p.Args["Workplaceid"].(int), p.Args["Date"], p.Args["Note"])
				if err != nil {
					return nil, err
				}
				id, _ := out.LastInsertId()
				return Types.LoadShiftById(int(id)), nil
			},
		},
		"editShift": &graphql.Field{
			Type: Types.ShiftType,
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
					_, err := database.Db.Exec("UPDATE shifts SET note = ? WHERE id = ?", p.Args["Note"].(string), p.Args["Id"].(int))
					if err != nil {
						return nil, err
					}
				}
				return Types.LoadShiftById(p.Args["Id"].(int)), nil
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
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	mux.Handle("/query", authMiddleware(queryHand()))

	mux.HandleFunc("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		var m interface{}
		err = json.Unmarshal(body, &m)

		if m != nil {
			if m.(map[string]interface{})["pass"] != nil && m.(map[string]interface{})["user"] != nil {
				h := md5.New()
				io.WriteString(h, m.(map[string]interface{})["pass"].(string))
				var username string
				var pass sql.NullString
				var id []uint8
				rows, err := database.Db.Query("select id, username, pass from users")
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&id, &username, &pass)
					if err != nil {
						log.Fatal(err)
					}
					if strings.ToLower(m.(map[string]interface{})["user"].(string)) == strings.ToLower(username) {
						token := randToken()
						if pass.Valid {
							passhash := new(bytes.Buffer)
							fmt.Fprintf(passhash, "%x", h.Sum(nil)) //cast pass hash to var
							if passhash.String() == pass.String {
								_, err := database.Db.Exec("INSERT INTO logins (user_id, token) VALUES (?, ?)", id, token) //Save Token to db
								if err != nil {
									panic(err)
								}
								resp, err := json.Marshal(map[string]interface{}{"token": token})
								w.Write(resp) //USER IS LOGIN, send him token
								return
							}
						} else {
							_, err := database.Db.Exec("INSERT INTO logins (user_id, token) VALUES (?, ?)", id, token) //Save Token to db
							if err != nil {
								panic(err)
							}
							resp, err := json.Marshal(map[string]interface{}{"token": token, "first": true})
							w.Write(resp) //USER IS LOGIN, send him token
							return
						}

					}
				}
				err = rows.Err()
				resp, err := json.Marshal(map[string]interface{}{"error": "unknown-user-or-pass"})
				w.Write(resp)
				return
			}
		}
		resp, err := json.Marshal(map[string]interface{}{"error": err})
		w.Write(resp)
	}))

	mux.HandleFunc("/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "" {
			_, err := database.Db.Exec("DELETE FROM logins WHERE token = ?", r.Header.Get("token")) //delete token
			if err != nil {
				panic(err)
			}
			resp, _ := json.Marshal(map[string]interface{}{"success": true})
			w.Write(resp)
		} else {
			resp, _ := json.Marshal(map[string]interface{}{"error": "no_header"})
			w.Write(resp)
		}

	}))

	handler := cors.AllowAll().Handler(mux)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
	//log.Fatal(http.ListenAndServe(":80", handler))

}
