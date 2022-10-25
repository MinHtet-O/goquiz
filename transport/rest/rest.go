package rest

import (
	_ "github.com/lib/pq"
	"goquiz/service"
	"sync"
)

const version = "1.0.0"

type Application struct {
	model *service.Model
	wg    sync.WaitGroup
}

func NewRESTServer(model *service.Model) service.Transport {
	app := &Application{
		model: model,
	}
	return app
}

// TODO: figure out why scrap results are different each time the program run
// TODO: change standard logging to json format, replace all log stdout stderr with json logging
// TODO: why file write is not working
// TODO: move URL field from questions to categories
// TODO: make category table composite primary key (id,name). So, there would be no duplicate category name
