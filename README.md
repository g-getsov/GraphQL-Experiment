# GraphQL-Experiment
A GraphQL API backed by a Neo4J DB built in Go running in a docker container with available code hotswap and Delve Go debugger for remote debugging

## How to build & use

1. Clone the project in the src folder of your $GOPATH

2. Open a terminal in the project folder and use `docker-compose build` or `docker-compose build -e hotswap=1 debug=0` to build the docker images (replace the 0 with 1 depending on what features you want the container to have enabled). Hotswap and debug are turned off by default.

3. Once the images are finished building run `docker-compose up` to spin up the containers

4. Once the docker containers are running graphql is ready to be used.

5. To navigate to the graphiql GUI open a browser and go to `http://localhost:9000/graphiql`

6. To navigate to the neo4j GUI open a browsr and go to `http://localhost:7474/browser/`

7. To programaticaly access the graphql API make http requests to `http://localhost:9000/graphql`

##### NOTES
When running the project in debug mode in the container you can connect to the debugger remotely at the following address `localhost:2345`

When running the project with the hotswap app listening for changes the project will get rebuild everytime a go file is changed.

To ssh into a running container use `docker exec -ti <put container id here> sh` 

## Stack & credits:

1. Neo4J - The Internet-Scale, Native Graph Database (https://neo4j.com/)

2. graphql-go/graphql - An implementation of GraphQL for Go / Golang (https://github.com/graphql-go/graphql)

3. Golang - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. (https://golang.org/)

4. Delve - Delve is a debugger for the Go programming language. The goal of the project is to provide a simple, full featured debugging tool for Go. (https://github.com/derekparker/delve)

5. Fresh - Fresh is a command line tool that builds and (re)starts your web application everytime you save a Go or template file. (https://github.com/pilu/fresh)

6. Docker - Docker is an open platform for developers and sysadmins to build, ship, and run distributed applications, whether on laptops, data center VMs, or the cloud (https://www.docker.com/)

## TO DO:
1. Figure out a way to run hotswap and the remote debugger at the same time (dlv attach to the new PID after refresh ?)

2. Figure out why running the project in debug mode automaticaly trigers a breakpoint and halts the program at the first line of code in the main function of the main.go file
