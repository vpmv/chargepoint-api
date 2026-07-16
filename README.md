ChargePoint API
---

Demo application


# Preface
I started out doing a little research on the requirements.
My initial instinct was to keep the project as light and portable as possible, sticking with a single binary and a SQLite database, since the requirements are minimal.
Since one of the key features is calculating point-to-point distance based on geographical data, I checked if my assumptions would meet the criteria.
I quickly found out there's a Postgres extension specialized in calculating geographical/geometrical data: PostGIS.
I spent a fair amount of time researching the best way to implement this, before and during development.

The rest of the application can remain fairly minimal. I'm used to separating the API and backend, using DTOs (data transfer objects), making for a clean presentation layer whilst leaving implementation details within the data layer(s).
Since only need one actual resource is required, we won't be splitting up the API further than necessary.

Although the requirements don't speak of Authentication/Authorization, I'm implementing an some boilerplate code into the API, which should be adapted to production standards.

## Task deviations
I've deviated from the task slightly by implementing Vendor ID instead of a generic ID. I figured this would be more alike a real world scenario, where vendors reference their assets by a unique identifier. For the purpose of [seeding](#database-seeding) I used the following format:
`<VENDOR_SHORTCODE><REGION>-<ASSETID>` e.g. `FN11-321`.

I chose to use Fuego as my HTTP microframework. I'm using it because it's lightweight and easy to use; it's my go-to framework. It's also a good fit for REST APIs.

The project layout may seem a bit excessive, but I wanted to keep the codebase organized and easy to navigate. I also wanted to make it easy to add new resources in the future. 

## AI disclosure

I refrained from using any generative AI services (e.g. Claude, Junie, Aider) in the development of this application. The main architecture is taken from earlier applications. I used AI to aid in parts of the research and used it alongside as a sparring partner. ChatGPT was a great help in integrating PostGIS, with which I was wholly unfamiliar. 


# Dev setup and runtime
The application is written in Go and uses PostgreSQL as a database. It comes with a configured Docker Compose file to run the application and the database. No other setup is required.


## Configuration
The default application environment variables are located in /config/.env.development. You can override them by creating your own file (.env.development.local or .env.local). You can configure the Docker Compose bind mount to use a custom directory, in conjunction with the basedir flag (--basedir /path/to/dir/).

## Run the application
```shell
docker compose up
```

### Authorization
As stated earlier, Authorization tokens have been implemented. A token provider has been omitted, since this is wildly out of scope, and replaced by a hard-coded "SimpleAuthenticator". The default token is `secret`. Please supply the token in the Authorization header:
```shell
curl --request GET \
  --url http://localhost:8989/api/v1/chargepoints \
  --header 'Accept: application/json, application/xml' \
  --header 'Authorization: Bearer secret'
```

Or in the OpenAPI UI: `Authorization: [ Bearer secret ]`

## Database seeding
The database can be seeded with mock records by adding the --seed flag to the application entrypoint.
This allows for easy functionality testing, mainly of the endpoint `/api/v1/chargepoints/location`.

```yaml
services:
  go:
    # [...]
    command: ["/server", "--seed"]
```

The seeder adds mock records to the database. It is intended for development purposes only and should not be used in production.

Seeded data simulates charge points by a set number of vendors, grouped by a logical number of regions. You can alter the amount of records per vendor region, with the environment variable SEED_COUNT. The default value is 100, times 10 regions = 1000 records per vendor. There are five pre-configured vendors:
- Shell (SH)
- Esso (ES)
- Gulf (GU)
- OK
- FastNed (FN)

$R = 5 * 10 * SC$

## Documentation / OpenAPI
The application automatically generates OpenAPI documentation. When the application is running, visit: http://localhost:8989/openapi to view the live documentation test suite (provided by Stoplight). The API can be extended to add more descriptions and clarifications to the interface. 

The generated documentation is automatically written to `/doc/openapi.json` on the file system. Please refer to the volume bindings in `docker-compose.yml`.

## Tests
At the time of writing, tests have not been implemented, but functionality has been tested manually.