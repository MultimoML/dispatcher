# Dispatcher

Microservice for managing grocery items.

Available endpoints:
- [`/live`](https://multimo.ml/products/live): Liveliness check
- [`/ready`](https://multimo.ml/products/ready): Readiness check
- [`/v1/all`](https://multimo.ml/products/v1/all): Returns a list of all products
- [`/v1/:id`](https://multimo.ml/products/v1/000000000000000000567522): Returns a single product by id

Available query parameters:
- [`limit`](https://multimo.ml/products/v1/all?limit=3): The number of returned products
- [`offset`](https://multimo.ml/products/v1/all?limit=3&offset=5): Which product to start from
- [`sort`](https://multimo.ml/products/v1/all?limit=3&sort=price): Sort by field (none, name, price, category)
- [`order`](https://multimo.ml/products/v1/all?limit=3&sort=name&order=desc): Sort order (asc, desc)
- [`history`](https://multimo.ml/products/v1/all?limit=1&history=full): Limit number of prices to display per item (last, none, all)
- [`category`](https://multimo.ml/products/v1/all?limit=3&category=Olives): Filter by category name

Branches:
- [`main`](https://github.com/MultimoML/dispatcher/tree/main): Contains latest development version
- [`prod`](https://github.com/MultimoML/dispatcher/tree/prod): Contains stable, tagged releases

## Setup/installation

Prerequisites:
- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)

Example usage:
- See all available options: `make help`
- Run microservice in a container: `make run`
- Release a new version: `make release ver=x.y.z`

All work should be done on `main`, `prod` should never be checked out or manually edited.
When releasing, the changes are merged into `prod` and both branches are pushed.
A GitHub Action workflow will then build and publish the image to GHCR, and deploy it to Kubernetes.

## License

Multimo is licensed under the [GNU AGPLv3 license](LICENSE).