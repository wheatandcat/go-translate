package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/wheatandcat/go-translate/lib"
	"github.com/wheatandcat/go-translate/lib/translate"
	"google.golang.org/appengine"
)

var ctx context.Context

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
					log.Println(lang)
					log.Println(words)

					w := []string{}
					for i := 0; i < words.Len(); i++ {
						w = append(w, words.Index(i).Interface().(string))
					}

					t, err := translate.Translate(ctx, w, lang)
					if err != nil {
						return nil, err
					}

					word := lib.NewWord()
					word.Words = w
					word.Translates = t
					word.Lang = lang

					if err := word.Put(ctx); err != nil {
						return nil, err
					}

					return t, nil
				},
			},
		},
	})

func customHandler(schema *graphql.Schema) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		opts := handler.NewRequestOptions(r)

		ctx = appengine.NewContext(r)

		rootValue := map[string]interface{}{
			"response": rw,
			"request":  r,
		}

		params := graphql.Params{
			Schema:         *schema,
			RequestString:  opts.Query,
			VariableValues: opts.Variables,
			OperationName:  opts.OperationName,
			RootObject:     rootValue,
		}

		result := graphql.Do(params)

		jsonStr, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")

		rw.Write(jsonStr)
	}
}

var schemaConfig = graphql.SchemaConfig{
	Query:    query,
	Mutation: mutation,
}

func init() {
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", customHandler(&schema))
}
