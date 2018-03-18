package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/graphql-go/handler"

	"github.com/graphql-go/graphql"
)

type PhotoOptions struct {
	ID      int `url:"id,omitempty"`
	AlbumID int `url:"albumId,omitempty"`
}

type Photo struct {
	ID           int    `json:"id"`
	AlbumID      int    `json:"albumId"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type PostOptions struct {
	ID     int `url:"id,omitempty"`
	UserID int `url:"userId,omitempty"`
}

type Time struct {
	ID   int    `json:"id"`
	Time string `json:"time"`
}

type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type City struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

var photoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Photo",
	Fields: graphql.Fields{
		"albumId": &graphql.Field{
			Type: graphql.Int,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"thumbnailUrl": &graphql.Field{
			Type: graphql.String,
		},
	},
})

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
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"time": &graphql.Field{
			Type: graphql.String,
		},
	},
})

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
		"id": &graphql.Field{
			Type: graphql.ID,
		},
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
				timeNow := Time{rand.Intn(1000), now.String()}

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
				opts := PostOptions{}
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
		"getPhotos": &graphql.Field{
			Type: &graphql.List{
				OfType: photoType,
			},
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"albumID": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				opts := PhotoOptions{}
				url := "https://jsonplaceholder.typicode.com/photos"

				id, okID := p.Args["id"]
				albumID, okAlbumID := p.Args["albumID"]

				if okID == true {
					opts.ID = id.(int)
				}

				if okAlbumID == true {
					opts.AlbumID = albumID.(int)
				}

				v, _ := query.Values(opts)

				if v.Encode() != "" {
					url = url + "?" + v.Encode()
				}

				var photos = []Photo{}

				response, err := http.Get(url)

				if err != nil {
					fmt.Printf("%s", err)
				} else {
					defer response.Body.Close()

					err := json.NewDecoder(response.Body).Decode(&photos)

					if err != nil {
						fmt.Printf("%s", err)
					}
				}

				return photos, nil
			},
		},
	},
})

type HttpHandlerWrapper struct {
	Handler http.Handler
}

func (wrapper *HttpHandlerWrapper) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	wrapper.Handler.ServeHTTP(rw, req)
}

func main() {
	createDummyData()

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	wrapper := &HttpHandlerWrapper{
		Handler: h,
	}
	http.Handle("/graphql", wrapper)
	http.ListenAndServe(":8080", nil)
}
