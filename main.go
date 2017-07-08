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

	"github.com/lordpuma/webserver/Types"
	"github.com/rs/cors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//var executor *graphql.Executor

func init() {
	var err error
	if err != nil {
		panic(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int32
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

func test() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint(r.Context().Value("user_id"))))
	})
}

//func queryHand() http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		var v map[string]interface{}
//		ctx := map[string]interface{}{"user_id": r.Context().Value("user_id")}
//
//		body, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			fmt.Fprintf(w, "%s", err)
//		}
//		var m interface{}
//		errr := json.Unmarshal(body, &m)
//		if errr != nil {
//			fmt.Fprintf(w, "%s", errr)
//		}
//
//		q := m.(map[string]interface{})["query"].(string)
//		if m.(map[string]interface{})["variables"] != nil {
//			v = m.(map[string]interface{})["variables"].(map[string]interface{})
//		}
//
//		result, err := executor.Execute(ctx, q, v, "")
//		resp, _ := json.Marshal(result)
//		fmt.Fprintf(w, "%s", resp)
//
//	})
//}
func main1() {
	//r := resolvers.GetResolvers()
	var err error
	//executor, err = graphql.NewExecutor(resolvers.Schema, "Query", "Mutation", r)

	//db, err := sql.Open("mysql", "root:password@tcp(db:3306)/database")
	db, err := sql.Open("mysql", "root:pass@/database")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	database.Connect(db)

	//CORS STARTS HERE - DEV ONL	Y
	mux := http.NewServeMux()

	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	mux.Handle("/test", authMiddleware(test()))
	//mux.Handle("/query", authMiddleware(queryHand()))
	//mux.Handle("/query", authMiddleware(queryHand()))
	//mux.Handle("/query", queryHand())

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

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.AllowAll().Handler(mux)

	//CORS ENDS HERE - DEV ONLY

	//PRODUCTION START

	//http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write(page)
	//}))
	//
	//http.Handle("/query", &relay.Handler{Schema: schema})

	//PRODUCTION END

	//SEED DB
	time.Sleep(5000 * time.Millisecond)
	var c int32
	err = database.Db.QueryRow("select count(id) as c from users").Scan(&c)
	if err != nil {
		log.Fatal(err)
	}
	if c == 0 {
		_, err := database.Db.Exec("INSERT INTO users (username, first_name, last_name, bg_color, color, email) VALUES (?, ?, ?, ?, ?, ?)", "lordpuma", "Tomáš", "Korený", "#000000", "#FFFFFF", "")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(http.ListenAndServe(":8080", handler))

}

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
	db, err := sql.Open("mysql", "root:pass@/database")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	database.Connect(db)

	fields := graphql.Fields{
		"user": &graphql.Field{
			Type: Types.UserType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadUserById(p.Args["id"].(int)), nil
			},
		},
		"workplace": &graphql.Field{
			Type: Types.WorkplaceType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Types.LoadWorkplaceById(p.Args["id"].(int)), nil
			},
		},

	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
		{
			user(id: 1) {workplaces{id}}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}

}
