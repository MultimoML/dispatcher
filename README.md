# Dispatcher

Microservice for managing grocery items.

Available endpoints:
- `/live`: Liveliness check
- `/ready`: Readiness check
- `/v1/all`: Returns a list of all products
- `/v1/:id`: Returns a single product by id

Available query parameters:
- `limit`: The number of returned products
- `offset`: Which product to start from
- `sort`: Sort by field (none, name, price, category)
- `order`: Sort order (asc, desc)
- `history`: Limit number of prices to display per item (last, none, all)
- `category`: Filter by category name

Branches:
- `main`: Contains latest development version
- `prod`: Contains stable, tagged releases

## Setup/installation

Prerequisites:
- Go
- Docker

To run the microservice run `make run`.
To see all the available options run `make help`.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).
