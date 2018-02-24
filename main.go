package main

import (
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/wheatandcat/go-translate/utils/translate"
)

var query = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		},
	})

var mutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"translate": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Args: graphql.FieldConfigArgument{
					"words": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
					"lang": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					words := reflect.ValueOf(p.Args["words"])
					lang, _ := p.Args["lang"].(string)

					w := []string{}
					for i := 0; i < words.Len(); i++ {
						w = append(w, words.Index(i).Interface().(string))
					}

					r, err := translate.Translate(w, lang)
					if err != nil {
						return nil, err
					}

					return r, nil
				},
			},
		},
	})

var schemaConfig = graphql.SchemaConfig{
	Query:    query,
	Mutation: mutation,
}

func main() {

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
