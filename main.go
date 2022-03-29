package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
)

type Job struct {
	Id             int      `json:"id"`
	Position       string   `json:"position"`
	Company        string   `json:"company"`
	Location       string   `json:"location"`
	Description    string   `json:"description"`
	EmploymentType string   `json:"employmentType"`
	SkillsRequired []string `json:"skillsRequired"`
}

var jobType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Job",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"position": &graphql.Field{
				Type: graphql.String,
			},
			"company": &graphql.Field{
				Type: graphql.String,
			},
			"location": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"employmentType": &graphql.Field{
				Type: graphql.String,
			},
			"skillsRequired": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

type reqBody struct {
	Query string `json:"query"`
}

func gqlHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No query data", 400)
			return
		}

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error parsing JSON request body", 400)
		}

		// fmt.Fprintf(w, "%s", processQuery(rBody.Query))

	})
}

func main() {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/api/graphql")
	if err != nil {
		panic(err)
	}

	// http.Handle("/api/graphql", gqlHandler())
	http.Handle("/api/graphiql", graphiqlHandler)
	http.ListenAndServe(":4000", nil)
	fmt.Println("GraphQL server running on http://localhost:4000/api/graphql. GraphiQL listening on http://localhost:4000/api/graphiql 🚀")
}
