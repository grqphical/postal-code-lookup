# Canadian Postal Code API

[![Go Tests](https://github.com/grqphical/postal-code-lookup/actions/workflows/go.yml/badge.svg)](https://github.com/grqphical/postal-code-lookup/actions/workflows/go.yml)

This is an API that allows you to submit a Canadian postal code and recieve information about such as which province in belongs too, which city its in, etc.

## Installation

Clone the repo, make sure you have Golang installed, then run:

```bash
$ make run
```

## Docs

Swagger API Docs are available at `/docs`

## Endpoints

### `/`

A simple status endpoint that returns `OK` if the API is up

### `/v1/postal-code/:postalCode`

This is the main endpoint of the API. When given a valid Canadian postal code, it will return all relevant information about it.

If the postal code is invalid it will return an HTTP Status Code of `400`

All data returned from this endpoint will be in JSON format

## License

This project is licensed under the GPL V3 License
