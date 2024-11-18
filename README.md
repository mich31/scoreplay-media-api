# ScorePlay Media API
## Description
This project is an implementation of a REST API providing endpoints to manage media assets &amp; tags. It provides the following functionalities:
- Create a tag
- List all tags
- Search tags by name
- Delete tag
- Create a media
- Search medias by tag

## Architecture
This application has been implemented with [Go](https://go.dev/doc/install) and [Fiber](https://docs.gofiber.io/) which is a famous framework to easily build REST APIs in [Go](https://go.dev/doc/install). 

It uses a [PostgreSQL](https://www.postgresql.org/) database. **PostgreSQL** is easy to use as a SQL database and handles well the logic of this application. Any other SQL database like [MySQL](https://www.mysql.com/) or NoSQL like [MongoDB](https://www.mongodb.com/) could have been use in this case. This database includes 3 tables: media (media entities), tags (tag entities), media_tags(manage many-to-many association between medias and tags).

[GORM](https://gorm.io/) manages interactions between the application and the database. This ORM library is easy to use and provides a straightforward [documentation](https://gorm.io/docs/).

[MinIO](https://min.io/) is used here as a storage service to manage &amp; store media files created. It seems more relevant to use a dedicated storage service for media management than use a database for scalability, security and cost-effectiveness concerns. [MinIO](https://min.io/) provides a pretty simple Go SDK, similar functionalities than [Amazon S3](https://aws.amazon.com/s3/) or any other famous cloud storage service (GCP, Azure Blob Storage), a WEBUI (available at `http://127.0.0.1:9001` if you run it via a Docker) and an API (available at `http://127.0.0.1:9000` if you run it via a Docker). You can find the credentials (`MINIO_ROOT_USER` & `MINIO_ROOT_PASSWORD`) in `.env.example` file.

For simplicity and effectiveness, both the [PostgreSQL](https://www.postgresql.org/) database and [MinIO](https://min.io/) will be run as Docker containers.

**P.S:** You may need to create an access key via **MinIO** WebUI (available at http://127.0.0.1:9000 if you run it via a Docker) if you get an error (`The Access Key Id you provided does not exist in our records`) when running the application. This case is handled through the instructions set in the `docker-compose.yaml` file.

## How to run
### Prerequisites
1. Install [Docker](https://docs.docker.com/get-docker/) and Docker Compose
2. Install [Go 1.23+](https://go.dev/doc/install)
3. Clone this repository
### Environment Setup
1. Open a terminal and go to the root folder of this repository
```
cd scoreplay-media-api
```
2. Create a copy of `.env.example` file and name it `.env`:
```
cp .env.example .env
```
### Run the service
1. Run the database and storage service using Docker Compose
```
docker compose up -d
```
2. Install the dependencies
```
go mod tidy
```
3. Build &amp; generate the artifact
```
go build
```
4. Run the service
```
./scoreplay-media-api
```
The service is running on `http://localhost:3000`.

API documentation is available in `docs/swagger.yaml` or `http://localhost:3000/swagger`.

## Testing
Run tests
```
go test ./...
```

## Technologies
- [Go 1.23+](https://go.dev/doc/install)
- [Fiber v2](https://docs.gofiber.io/)
- [GORM](https://gorm.io/)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [GoDotEnv](github.com/joho/godotenv)
- [MinIO](https://min.io/)
- [Swaggo](https://github.com/swaggo/swag)

## Improvements
- **Logging**: More structured logs need to be added with different log levels (DEBUG, INFO, ERROR) for better monitoring and debugging.
- **Testing**: The current tests suite covers the controllers logic and an integration the  services. Test coverage needs to be improved with more unit and integration tests. An in-memory could even be used for to end-to-end integrations and less rely on mocks.
- **Pagination & Filtering**: Pagination needs to be implemented for scalability and better performance. Additional filters can be added to provide limiting &amp; sorting capabilities.
- **Storage**: [MinIO](https://min.io/) is a nice solution for prototyping. The integration with production-ready service like [Amazon S3](https://aws.amazon.com/s3/), [Google Cloud Storage](https://cloud.google.com/storage) or [Azure Blob Storage](https://azure.microsoft.com/en-us/products/storage/blobs) can be implemented.
- **API documentation**: [Swaggo](https://github.com/swaggo/swag) helps to generate swagger documentation with annotations but there is room for improvement on the result. In my opinion, it is interestin to use this library to get a 1st draft version and then improve it.
- **File management**:
    1. A limit can be set for the input file size on `POST /api/medias` endpoint. It depends on the product requirements but it could help to control resource consumption and service availability.
    2. File processing can be improved by delegating file upload to a messaging service
    3. It might be useful to implement file compression and thumbnail generation. This will help to manage costs especially for large files if the storage is managed by a cloud service.
    4. File type checks should be implemented for security concerns.
- **Caching**: Caching can be implemented for the most used tags &amp; medias using a technology like [Redis](https://redis.io). It could help to maintain a good performance on a system which may have to handle a large amount of medias &amp; tags.
