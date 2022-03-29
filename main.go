package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

// open json file & return data
func retriveJobsFromFile() func() []Job {
	return func() []Job {
		jsonf, err := os.Open("data.json")

		if err != nil {
			fmt.Printf("Error opening json file: %s", err)
		}
		data, _ := ioutil.ReadAll(jsonf)
		defer jsonf.Close()

		var jobsData []Job

		err = json.Unmarshal(data, &jobsData)

		if err != nil {
			fmt.Printf("Error unmarshalling json file: %s", err)
		}
		return jobsData

	}
}

func gqlSchema(queryJobs func() []Job) graphql.Schema {
	fields := graphql.Fields{
		"jobs": &graphql.Field{
			Type:        graphql.NewList(jobType),
			Description: "List of jobs",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return queryJobs(), nil
			},
		},
		"job": &graphql.Field{
			Type:        jobType,
			Description: "Get a job by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, job := range queryJobs() {
						if int(job.Id) == id {
							return job, nil
						}
					}
				}
				return nil, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic("failed to create new schema, error: " + err.Error())
	}
	return schema
}

func processQuery(query string) (result string) {
	retrieveJobs := retriveJobsFromFile()

	params := graphql.Params{Schema: gqlSchema(retrieveJobs), RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	return fmt.Sprintf("%s", rJSON)
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
		fmt.Fprintf(w, "%s", processQuery(rBody.Query))
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
