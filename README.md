# User-Posts examples

## Directory structure
- `business_logic` as its name says interacts with http layer to present the data 
  and make requests to user and post clients.
- `datasource` contains the logic for user and posts clients.
- `datasource_models` are the models which requests responses parses the results.
- `service` contains the logic for http serving.
- `service_models` contains the structures that are presented as response bodies.

## To Test
run `go test ./...`

## To Build
run `go build`