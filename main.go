// credit - go-graphql hello world example
package main

import (
	"encoding/json"

	"log"
	"net/http"
	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	ID        int
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"Body": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"tutorials": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	},
)

var tutorialType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tutorial",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: authorType,
			},
			"comments": &graphql.Field{
				Type: graphql.NewList(commentType),
			},
		},
	},
)

var AuthorInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "AuthorInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

func populate() ([]Tutorial, []Author) {
	author1 := Author{ID: 1, Name: "Elliot Forbes", Tutorials: []int{1}}
	author2 := Author{ID: 2, Name: "Max Lopez", Tutorials: []int{2}}
	tutorial1 := Tutorial{
		ID:     1,
		Title:  "Go GraphQL Tutorial",
		Author: author1,
		Comments: []Comment{
			Comment{Body: "First Comment"},
		},
	}
	tutorial2 := Tutorial{
		ID:     2,
		Title:  "Pandit Football Tutorial",
		Author: author2,
		Comments: []Comment{
			Comment{Body: "Second Comment"},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial1)
	tutorials = append(tutorials, tutorial2)

	var authors []Author
	authors = append(authors, author1)
	authors = append(authors, author2)

	return tutorials, authors
}

func main() {
	var tutorials, authors = populate()
	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get Tutorial By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"author": &graphql.Field{
			Type:        authorType,
			Description: "Get Author By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, author := range authors {
						if int(author.ID) == id {
							return author, nil
						}
					}
				}
				return nil, nil
			},
		},
		"tutorialList": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Tutorial List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
		"authorList": &graphql.Field{
			Type:        graphql.NewList(authorType),
			Description: "Get Author List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return authors, nil
			},
		},
	}

	mutationFields := graphql.Fields{
		"createTutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Create a new Tutorial",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"author": &graphql.ArgumentConfig{
					Type: AuthorInput,
				},
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				author := params.Args["author"].(map[string]interface{})
				authorType := Author{
					ID:   author["id"].(int),
					Name: author["name"].(string),
				}

				tutorial := Tutorial{
					ID:     params.Args["id"].(int),
					Title:  params.Args["title"].(string),
					Author: authorType,
				}

				tutorials = append(tutorials, tutorial)
				return tutorial, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)

	// query := `
	// 	mutation {
	// 		createTutorial(title: "Hello World", id: 5, author: {id: 10, name: "Bastard"}) {
	// 			title
	// 			id
	// 			author{
	// 				id
	// 				name
	// 			}
	// 		}
	// 	}
	// `

	// params := graphql.Params{Schema: schema, RequestString: query}
	// r := graphql.Do(params)
	// if len(r.Errors) > 0 {
	// 	log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	// }
	// rJSON, _ := json.Marshal(r)
	// fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}

	// query = `
  //   {
  //       tutorialList {
  //           id
	// 					title
	// 					author{
	// 						id
	// 						name
	// 					}
  //       }
  //   }
	// `
	// params = graphql.Params{Schema: schema, RequestString: query}
	// r = graphql.Do(params)
	// if len(r.Errors) > 0 {
	// 	log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	// }
	// rJSON, _ = json.Marshal(r)
	// fmt.Printf("%s \n", rJSON)
}
