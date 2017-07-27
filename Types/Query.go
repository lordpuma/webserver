package Types

import "github.com/graphql-go/graphql"

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
