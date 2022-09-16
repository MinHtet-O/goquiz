# Goquiz

## About the project

Goquiz provides the scrapers to get the engineering multiple-choice questions across the sites, structure the data, populate the database, and Serve API endpoints to retrieve the questions. 
You can also start the services with db-less mode. 

Currently, Goquiz can scrap the MCQ questions from the following sites. 
- https://www.javatpoint.com/

Scrapers for more web sites will be provided as the project progress. There are over **4000+ MCQ questions** from 74 different categories and all of them are credited to the respective original web source. 
This project is educational purpose only.


## Technical Description ##

The project structure was inspired by **hexagonal architecture**. Actors such as transport and repository are separated from service, which allows service to be a technology-agnostic component with only business logic inside.

The actors are loaded during runtime based on the config, i.e. , postgres data mode and in-memory data model can be swapped easily just by config value. It allows mocking and significantly improves testability.

The scraper spawns go routine to scrap each web page, consolidate and categorized the data before populating to the database.

Goquiz can be served as either REST or gRPC API ( still in development ). PostgresSQL is used as persistent storage to store the categorized data. For easier deployment, db-less mode with in-memory data model can also be used.

## Functions

#### Web Scraping
* Get the mcq questions from the web sites
* Group the questions into separate categories
* Write the questions for each categories to the Postgres database

#### API to retrieve questions
* Get categories by category ID endpoint to retrieve all available categories
* Get questions by categories endpoint to retrieve all questions from one category
* Key based API authentication
* API rate limiting

## Setup instructions

#### Prerequisite
Download and install [go](https://go.dev/doc/install) on your machine and clone the goquiz project. For deployment, you can deploy without database or setup and populate the database first before serving the API. With later option, you don't have to scrap the web everytime you start/ restart the service.

#### DB-less mode
1. Build the Project<br>```go build ./cmd/api```
2. Start API service to retrieve questions 
    - without apikey authentication <br>```./api```
    - with apikey authentication <br>```./api -apikey=<your_api_key>```
3. For more startup parameters <br>```./api --help```
    
#### DB mode
1. [Setup migrate cli](https://github.com/golang-migrate/migrate)
2. [Setup Postgres and create a database](https://www.prisma.io/dataguide/postgresql/setting-up-a-local-postgresql-database)
3. Expose database service name <br> ```export GOQUIZ_DB=postgres://<username>:@localhost/<db_name>?sslmode=disable```
4. Migrate the database, create necessary tables<br>```migrate -path=./resources/migrations -database=$GOQUIZ_DB up```
5. Build the Project<br>```go build ./cmd/api```
6. Scrap the questions and populate the database<br>```./api -populate-db -db-dsn=$GOQUIZ_DB```
7. Start API service to retrieve questions 
    - without apikey authentication <br>```./api -db-dsn=$GOQUIZ_DB```
    - with apikey authentication <br>```./api -db-dsn=$GOQUIZ_DB -apikey=<your_api_key>```
8. For more startup parameters <br>```./api --help```

## API Routes

1. Get all categories<br> ```curl --request GET \
  --url http://localhost:4000/v1/categories \
  --header 'Authorization: Key 1234'```
2. Get questions by category ID<br>```curl --request GET \
  --url 'http://localhost:4000/v1/questions?category_id=1'```
3. Create new category<br>```curl --request POST \
  --url http://localhost:4000/v1/categories \
  --header 'Content-Type: application/json' \
  --data '{
	"name":"Go Programming"
}'```
4. Create new question by category<br>```curl --request POST \
  --url http://localhost:4000/v1/questions \
  --header 'Authorization: Key 1234' \
  --header 'Content-Type: application/json' \
  --data '{
	"categ_id": 75,
	"text": "Which company created go programming language?",
	"answers": [
		"Apple",
		"Google",
		"Amazon",
		"Facebook"
	],
	"correct_answer": 2,
	"codeblock": "fmt.Println(\"Hello! Example codeblock\")",
	"explanation": "Go was originally designed at Google in 2007"
}'```

## Process Diagram
![alt text](https://github.com/MinHtet-O/goquiz/blob/main/resources/diagrams/goquiz_communication.png)

[//]: # (## Layout)

[//]: # ()
[//]: # (```tree)

[//]: # ()
[//]: # (├── .gitignore)

[//]: # (├── README.md)

[//]: # (├── cmd)

[//]: # (│   ├── api)

[//]: # (│   │   └── categories.go)

[//]: # (│   │   └── errors.go)

[//]: # (│   │   └── healthcheck.go)

[//]: # (│   │   └── main.go)

[//]: # (│   │   └── middleware.go)

[//]: # (│   │   └── questions.go)

[//]: # (│   │   └── routes.go)

[//]: # (│   │   └── server.go)

[//]: # (├── go.mod)

[//]: # (├── go.sum)

[//]: # (├── migrations)

[//]: # (│   └── 000001_create_db_category.down.sql)

[//]: # (│   └── 000001_create_db_category.up.sql)

[//]: # (│   └── 000002_create_db_questions.down.sql)

[//]: # (│   └── 000002_create_db_questions.up.sql)

[//]: # (│   └── 000003_view_questions_by_categories.down.sql)

[//]: # (│   └── 000003_view_questions_by_categories.up.sql)

[//]: # (├── pkg)

[//]: # (│   ├── model)

[//]: # (│   │   └── errors.go)

[//]: # (│   │   └── format.go)

[//]: # (│   │   └── io.go)

[//]: # (│   │   └── model.go)

[//]: # (│   │   └── postgres)

[//]: # (│   │   │   └── categories.go)

[//]: # (│   │   │   └── model.go)

[//]: # (│   │   │   └── questions.go)

[//]: # (│   ├── scraper)

[//]: # (│   │   └── network.go)

[//]: # (│   │   └── nodes.go)

[//]: # (│   │   └── parser.go)

[//]: # (│   │   └── scrapper.go)

[//]: # (│   └── validator)

[//]: # (│   │   └── category.go)

[//]: # (│   │   └── validator.go)

[//]: # (├── files)

[//]: # (│   └── mcq.txt)

[//]: # (```)

[//]: # (A brief description of the directory layout:)

[//]: # (* `cmd` contains main packages, each subdirectory of `cmd` can be built into executable.)

[//]: # (* * `api` that provides endpoints to retrieve the questions)

[//]: # (* `migrations` contains migration files to create the necessary tables)

[//]: # (* `pkg` contans most of the business logic)

[//]: # (* * `scraper` includes logic that fetch mcq questions from the web and write to the database)

[//]: # (* * `model` includes data models and it's helper functions like formatter and io.)

[//]: # (* * `validator` includes validation logic for data model)

[//]: # (* `files` contans static file assets. Currently, there is mcq.txt that contains the mcq urls from the sites.)

### ToDo
- [ ] Add GRPC Transport
- [ ] Add Structure logger
- [ ] Dockerize Deployment
- [ ] Deploy the API to Digital Ocean VPS
- [ ] Unit Testings for Model and APIs
- [ ] Add github CI to automate unit tests and deployment
- [ ] Update/ Delete endpoints for questions and categories 
- [ ] Add web scraper to fetch MCQs from www.sanfoundry.com

