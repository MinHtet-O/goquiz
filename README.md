# Goquiz

## About the project

Goquiz provides the scrapers to get the engineering multiple-choice questions across the sites, structure the data, populate the database, and serve API endpoints to retrieve the questions.

Currently, Goquiz can scrap the MCQ questions from the following sites. 
- https://www.javatpoint.com/

Scrapers for more web sites will be provided as the project progress. There are over **4000+ MCQ questions** from 74 different categories and all of them are credited to the respective original web source. This project is educational purpose only.

### Functions

#### Web Scraping
* Get the mcq questions from the web sites
* Group the questions into separate categories
* Write the questions for each categories to the Postgres database

#### API to retrieve questions
* Get categories by category ID endpoint to retrieve all available categories
* Get questions by categories endpoint to retrieve all questions from one category
* Key based API authentication
* API rate limiting

### Setup instruction

1. Clone the goquiz project
2. [Setup migrate cli](https://github.com/golang-migrate/migrate)
3. [Setup Postgres and create a database](https://www.prisma.io/dataguide/postgresql/setting-up-a-local-postgresql-database)
4. Expose database service name <br> ```export GOQUIZ_DB=postgres://<username>:@localhost/<db_name>?sslmode=disable```
5. Migrate the database, create necessary tables<br>```migrate -path=./migrations -database=$GOQUIZ_DB up```
6. Build the Project<br>```go build ./cmd/api```
7. Scrap the questions and populate the database<br>```./api -scrap -db-dsn=$GOQUIZ_DB```
8. Start API service to retrieve questions 
    - without apikey authentication <br>```./api -db-dsn=$GOQUIZ_DB```
    - with apikey authentication <br>```./api -db-dsn=$GOQUIZ_DB -apikey=<your_api_key>```
9. For more startup parameters <br>```./api --help```

### Process Diagram
![alt text](https://github.com/MinHtet-O/goquiz/blob/main/goquiz_communication.png)

### Layout

```tree

├── .gitignore
├── README.md
├── cmd
│   ├── api
│   │   └── categories.go
│   │   └── errors.go
│   │   └── healthcheck.go
│   │   └── main.go
│   │   └── middleware.go
│   │   └── questions.go
│   │   └── routes.go
│   │   └── server.go
├── go.mod
├── go.sum
├── migrations
│   └── 000001_create_db_category.down.sql
│   └── 000001_create_db_category.up.sql
│   └── 000002_create_db_questions.down.sql
│   └── 000002_create_db_questions.up.sql
│   └── 000003_view_questions_by_categories.down.sql
│   └── 000003_view_questions_by_categories.up.sql
├── pkg
│   ├── model
│   │   └── errors.go
│   │   └── format.go
│   │   └── io.go
│   │   └── model.go
│   │   └── postgres
│   │   │   └── categories.go
│   │   │   └── model.go
│   │   │   └── questions.go
│   ├── scraper
│   │   └── network.go
│   │   └── nodes.go
│   │   └── parser.go
│   │   └── scrapper.go
│   └── validator
│   │   └── category.go
│   │   └── validator.go
├── files
│   └── mcq.txt
```
A brief description of the directory layout:
* `cmd` contains main packages, each subdirectory of `cmd` can be built into executable.
* * `api` that provides endpoints to retrieve the questions
* `migrations` contains migration files to create the necessary tables
* `pkg` contans most of the business logic
* * `scraper` includes logic that fetch mcq questions from the web and write to the database
* * `model` includes data models and it's helper functions like formatter and io.
* * `validator` includes validation logic for data model
* `files` contans static file assets. Currently, there is mcq.txt that contains the mcq urls from the sites.
