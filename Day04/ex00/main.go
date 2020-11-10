package main

import (
	"fmt"
	"log"

	"github.com/dm-alexi/go_piscine/Day04/ex00/swagger/restapi"
	"github.com/dm-alexi/go_piscine/Day04/ex00/swagger/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func main() {
	pricelist := map[string]int64{"CE": 10, "AA": 15, "NT": 17, "DE": 21, "YR": 23}
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewCandyServerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = 3333
	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(
		func(params operations.BuyCandyParams) middleware.Responder {
			if pricelist[*params.Order.CandyType] == 0 || *params.Order.CandyCount <= 0 {
				response := operations.BuyCandyBadRequestBody{Error: "Invalid candy type"}
				if *params.Order.CandyCount <= 0 {
					response.Error = "Invalid candy count"
				}
				return operations.NewBuyCandyBadRequest().WithPayload(&response)
			}
			change := *params.Order.Money - *params.Order.CandyCount*pricelist[*params.Order.CandyType]
			if change < 0 {
				response := operations.BuyCandyPaymentRequiredBody{Error: fmt.Sprintf("You need %d more money!", -change)}
				return operations.NewBuyCandyPaymentRequired().WithPayload(&response)
			}
			response := operations.BuyCandyCreatedBody{Thanks: "Thank you!", Change: &change}
			return operations.NewBuyCandyCreated().WithPayload(&response)
		})
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
