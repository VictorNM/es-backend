# E-SHARING

![Go](https://github.com/VictorNM/es-backend/workflows/Go/badge.svg?branch=develop)
[![Build Status](https://travis-ci.com/VictorNM/es-backend.svg?branch=develop)](https://travis-ci.com/VictorNM/es-backend)
[![codecov](https://codecov.io/gh/VictorNM/es-backend/branch/develop/graph/badge.svg)](https://codecov.io/gh/VictorNM/es-backend)

[//]: <> (## Getting Started)

[//]: <> (### Prerequisites)

[//]: <> (### Installing)

## Running the tests

```bash
go test ./... -cover
```

[//]: <> (### Break down into end to end tests)

[//]: <> (### And coding style tests)

## Local development

- Following step assuming use are at the PROJECT_ROOT

### Setup Environment file

```bash
# copy the .env file
cp .env.example .env
```

### Running

#### Using Docker

Build

```bash
docker build --tag=es:latest .
```

Run

```bash
# Local inside VM
docker run --network host -d -p 127.0.0.1:8080:8080/tcp --env-file ./.env es:latest

# Local
docker run -d -p 127.0.0.1:8080:8080/tcp --env-file ./.env es:latest
```

#### Using Go

```bash
# Get godotenv binary (require once)
go get github.com/joho/godotenv/cmd/godotenv

# Run in Unix
godotenv go run main.go

# Run in Windows
godotenv.exe go run main.go
```

### API References

Open browser at: `{host}/swagger/index.html`

[//]: <> (## Built With)

[//]: <> (## Contributing)

[//]: <> (## Versioning)

## Authors

* **VictorNM** - *Initial work* - [VictorNM](https://github.com/VictorNM)

* **bobgel12** - [bobgel12](https://github.com/bobgel12)

[//]: <> (## License)

[//]: <> (## Acknowledgments)
