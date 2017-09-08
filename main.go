package main
import (
	"log"
	"net/http"
	"github.com/rs/cors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/mnmtanish/go-graphiql"
	sqlDriver "github.com/lib/pq"
	neoDriver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"database/sql/driver"
)

type DbConfig struct {
	Username string
	Password string
	Host string
	Port string
}

// Neo4j config
var neoConf = DbConfig{
	"neo4j",
	"over1234",
	"neo4j",
	"7687",
}

var postgresConf = DbConfig{
	"admin",
	"admin",
	"postgres",
	"5432",
}

type Person struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	From string `json:"from"`
	Friends []Person `json:"friends"`
}

type Hobby struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
}

func getPostgresConnection() (driver.Conn, error) {
	db, err := sqlDriver.Open("postgres://" +postgresConf.Username+ ":" +postgresConf.Password+ "@" +postgresConf.Host+ ":" +postgresConf.Port)

	if err != nil { log.Println("error connecting to neo4j", err) }

	return db, err
}

func getNeoConnection() (neoDriver.Conn, error) {
	db, err := neoDriver.NewDriver().OpenNeo("bolt://"+ neoConf.Username+":"+ neoConf.Password+"@"+ neoConf.Host+":"+ neoConf.Port)

	if err != nil { log.Println("error connecting to neo4j", err) }

	return db, err
}

func getPeople() []Person {

	db, _ := getNeoConnection()
	defer db.Close()

	cypher := `MATCH (n:Person) RETURN ID(n) as id, n.name as name, n.from LIMIT {limit}`

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface {}{ "limit": 25})

	if err != nil {
		log.Println("error querying person:", err)
		// w.WriteHeader(500)
		// w.Write([]byte("an error occured querying the DB"))
		// return
	} else if len(data) == 0 {
		// w.WriteHeader(404)
		// return
	}

	results := mapPeople(data)

	return results
}

func getPerson(name string) Person {

	db, _ := getNeoConnection()
	defer db.Close()

	cypher := "MATCH (p:Person) WHERE p.name = {name} RETURN ID(p) as id, p.name, p.from"

	data, err := db.QueryNeo(cypher, map[string] interface{}{"name": name})

	if err != nil { log.Println("error looking up person")
	} else if data == nil { log.Println("cant find person") }

	fields, _, err := data.NextNeo()

	result := mapPerson(fields)

	return result
}

func mapPeople(rows [][]interface{}) ([]Person) {

	people := make([]Person, len(rows))

	for idx, row := range rows {
		people[idx] =
			mapPerson(row)
	}

	return people
}

func mapPerson(row []interface{}) (Person) {

	person := Person{
		ID:    row[0].(int64),
		Name:  row[1].(string),
		From:  row[2].(string),
	}

	return person
}

func getFriends(id int) []Person {

	db, _ := getNeoConnection()
	defer  db.Close()

	cypher := "MATCH (p:Person)-[r :FRIEND]->(friend) WHERE ID(p) = {id} RETURN ID(friend) as id, friend.name, friend.from"

	data, _, _, err := db.QueryNeoAll(cypher, map[string] interface{}{"id": id})

	if err != nil { log.Println("error looking up person")
	} else if data == nil { log.Println("cant find person") }

	friends := make([]Person, len(data))

	for idx, row := range data {
		friends[idx] =
			mapPerson(row)
	}

	return friends
}

func getHobby(name string) Hobby {

	db, _ := getNeoConnection()
	defer db.Close()

	cypher := "MATCH (h:Hobby) WHERE h.name = {name} RETURN ID(h) as id, h.name"

	data, err := db.QueryNeo(cypher, map[string] interface{}{"name": name})

	if err != nil { log.Println("error looking up hobby")
	} else if data == nil { log.Println("cant find hobby") }

	fields, _, err := data.NextNeo()

	result := Hobby {
		ID:    fields[0].(int64),
		Name:  fields[1].(string),
	}

	return result
}

var hobbyType = graphql.NewObject(

	graphql.ObjectConfig{
		Name: "Hobby",
		Fields: graphql.Fields {
			"id" : &graphql.Field{
				Type: graphql.Int,
			},
			"name" : &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)


func compilePersonType() (*graphql.Object) {

	personType := graphql.NewObject(

		graphql.ObjectConfig{
			Name: "Person",
			Fields: graphql.Fields {
				"id" : &graphql.Field {
					Type: graphql.Int,
				},
				"name": &graphql.Field {
					Type: graphql.String,
				},
				"from": &graphql.Field {
					Type: graphql.String,
				},
			},
		},
	)

	//TODO export the field to make it reusable
	personType.AddFieldConfig(
		"friends",
		&graphql.Field {
			Type: graphql.NewList(personType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getFriends(p.Args["id"].(int)), nil
			},
		},
	)

	return personType
}

var personType = compilePersonType()

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"People": &graphql.Field{
			Type: graphql.NewList(personType),
			Description: "List of people",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getPeople(), nil
			},
		},
		"Person": &graphql.Field{
			Type: personType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Description: "A person",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getPerson(p.Args["name"].(string)), nil
			},
		},
		"Friends": &graphql.Field{
			Type: graphql.NewList(personType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Description: "List of the person's friends",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return getFriends(p.Args["id"].(int)), nil
			},
		},
		"Hobby": &graphql.Field{
			Type: hobbyType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:graphql.String,
				},
			},
			Description: "A hobby",
			Resolve: func(h graphql.ResolveParams) (interface{}, error) {
				return getHobby(h.Args["name"].(string)), nil
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	// serve HTTP
	serveMux := http.NewServeMux()

	// serveMux.HandleFunc("/neo", neo4jHandler)
	serveMux.Handle("/graphql", c.Handler(h))
	serveMux.HandleFunc("/graphiql", graphiql.ServeGraphiQL)
	http.ListenAndServe(":8080", serveMux)
}