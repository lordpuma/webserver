package Types

import (
	"github.com/graphql-go/graphql"
	"github.com/lordpuma/webserver/database"
	"log"
)

type Result struct {
	Id      int
	Name    string
	Time    float32
	race_id int
}

var ResultType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Result",
	Description: "Basic Race Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Result).Id, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Result).Name, nil
			},
		},
		"time": &graphql.Field{
			Type: graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(Result).Time, nil
			},
		},
		"race": &graphql.Field{
			Type: RaceType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return LoadRaceById(p.Source.(Result).race_id), nil
			},
		},
	},
})

var FormattedResultType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "FormattedResult",
	Description: "Formatted Result Object",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(FResult).name, nil
			},
		},
		"average": &graphql.Field{
			Type: graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(FResult).average, nil
			},
		},
		"min": &graphql.Field{
			Type: graphql.Float,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(FResult).min, nil
			},
		},
		"times": &graphql.Field{
			Type: graphql.NewList(graphql.Float),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(FResult).times, nil
			},
		},
		"race": &graphql.Field{
			Type: RaceType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return LoadRaceById(p.Source.(FResult).race_id), nil
			},
		},
	},
})

func LoadResultById(id int) Result {
	var (
		name    string
		time    float32
		race_id int
	)
	err := database.Db.QueryRow("SELECT name, time, race_id FROM results WHERE id = ?", id).Scan(&name, &time, &race_id)
	if err != nil {
		panic(err)
	}

	return Result{
		Id:      id,
		Name:    name,
		race_id: race_id,
	}
}

func LoadResultsList() []Result {
	var r []Result
	var (
		id      int
		name    string
		time    float32
		race_id int
	)
	rows, err := database.Db.Query("select id, name, time, race_id FROM results")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &time, &race_id)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, Result{
			Id:      id,
			Name:    name,
			race_id: race_id,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func LoadResultsByRace(r_id int) []FResult {
	var (
		id      int
		name    string
		time    float32
		race_id int
	)
	var ret []FResult
	var ret2 []FResult
	rows, err := database.Db.Query("select id, name, time, race_id FROM results WHERE race_id = ?", r_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &time, &race_id)
		if err != nil {
			log.Fatal(err)
		}

		var needle *FResult
		for k, v := range ret {
			if v.name == name {
				needle = &ret[k]
			}
		}
		if needle == nil {
			var a []float32
			ret = append(ret, FResult{name, append(a, time), time, time, race_id})
		} else {
			needle.times = append(needle.times, time)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range ret {
		var total float32
		for _, value := range e.times {
			total += value
			if value < e.min {
				e.min = value // found another smaller value, replace previous value in min
			}
		}
		e.average = total / float32(len(e.times))
		ret2 = append(ret2, e)
	}
	return ret2

}

type FResult struct {
	name    string
	times   []float32
	average float32
	min     float32
	race_id int
}
