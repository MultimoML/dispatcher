# Dispatcher

Microservice for managing grocery items.

Available endpoints:
- `/live`: Liveliness check
- `/ready`: Readiness check
- `/v1/all`: Returns a list of all products

Branches:
- `main`: Contains stable, tagged releases
- `dev`: Contains latest development version

## Setup/installation

To run the microservice using Docker Compose run `make compose`.

To see other available options run `make help`.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).
