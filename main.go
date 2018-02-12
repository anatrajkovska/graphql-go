package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/graphql-go/handler"

	"github.com/graphql-go/graphql"
)

type Options struct {
	ID     int `url:"id,omitempty"`
	UserID int `url:"userId,omitempty"`
}

type Time struct {
	Time string `json:"time"`
}

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"userId": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"body": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var timeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Time",
	Fields: graphql.Fields{
		"time": &graphql.Field{
			Type: graphql.String,
		},
	},
})

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
		"getCurrentTime": &graphql.Field{
			Type: timeType,
			Args: graphql.FieldConfigArgument{
				"city": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"continent": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				city := p.Args["city"].(string)
				continent := p.Args["continent"].(string)
				loc, _ := time.LoadLocation(continent + "/" + city)
				now := time.Now().In(loc)
				timeNow := Time{now.String()}

				return timeNow, nil
			},
		},
		"getPosts": &graphql.Field{
			Type: &graphql.List{
				OfType: postType,
			},
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"userID": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				opts := Options{}
				url := "https://jsonplaceholder.typicode.com/posts"

				id, okID := p.Args["id"]
				userID, okUserID := p.Args["userID"]

				if okID == true {
					opts.ID = id.(int)
				}

				if okUserID == true {
					opts.UserID = userID.(int)
				}

				v, _ := query.Values(opts)

				if v.Encode() != "" {
					url = url + "?" + v.Encode()
				}

				var posts = []Post{}

				response, err := http.Get(url)

				if err != nil {
					fmt.Printf("%s", err)
				} else {
					defer response.Body.Close()

					err := json.NewDecoder(response.Body).Decode(&posts)

					if err != nil {
						fmt.Printf("%s", err)
					}
				}

				return posts, nil
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
