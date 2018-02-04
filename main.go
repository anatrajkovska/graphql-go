package main

import (
	"net/http"

	"github.com/graphql-go/handler"

	"github.com/graphql-go/graphql"
)

type City struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func createDummyData() {
	cities := []string{
		"Skopje",
		"Bitola",
		"Veles",
		"Shtip",
		"Prilep",
		"Gevgelija",
		"Kumanovo",
	}
	temps := []int{
		18,
		19,
		17,
		16,
		15,
		20,
		21,
	}

	for i := 0; i < len(cities); i++ {
		var city = City{
			Name:  cities[i],
			Value: temps[i],
		}
		data = append(data, city)
	}
}

var data = []City{}

var tempType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Temp",
	Fields: graphql.Fields{
		"value": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"getCurrentTemperature": &graphql.Field{
			Type: tempType,
			Args: graphql.FieldConfigArgument{
				"city": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				cityValue := p.Args["city"].(string)

				for _, value := range data {
					if cityValue == value.Name {
						return value, nil
					}
				}
				return nil, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func main() {
	createDummyData()
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
