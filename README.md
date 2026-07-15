ChargePoint API
---

Demo application


# Preface
I started out doing a little research on the requirements.
My initial instinct was to keep the project as light and portable as possible, sticking with a single binary and a SQLite database, since the requirements are minimal.
Since one of the key features is calculating point-to-point distance based on geographical data, I checked if my assumptions would meet the criteria.
I quickly found out there's a Postgres extension specialized in calculating geographical/geometrical data: PostGIS.
I spent a fair amount of time researching the best way to implement this, before and during development.

The rest of the application can remain fairly minimal. I'm used to separating the API and backend, using DTOs (data transfer objects), making for a clean presentation layer whilst leaving implantation details within the data layer(s).
Since only need one actual resource is required, we won't be splitting up the API further than necessary.

Although the requirements don't speak of Authentication/Authorization, I'm implementing an some boilerplate code into the API, which should be adapted to production standards.


## AI disclosure

I refrained from using any generative AI services (e.g. Claude, Junie, Aider) in the development of this application. The main architecture is taken from earlier applications. I used AI to aid in parts of the research and used it alongside as a sparring partner. ChatGPT was a great help in integrating PostGIS, with which I was wholly unfamiliar. 


# Dev setup & runtime
The application is written in Go and uses PostgreSQL as a database. It comes with a configured Docker Compose file to run the application and the database. No other setup is required.


## Configuration
The default application environment variables are located in /config/.env.development. You can override them by creating your own file (.env.development.local or .env.local). You can configure the Docker Compose bind mount to use a custom directory, in conjunction with the basedir flag (--basedir /path/to/dir/).

## Run the application
```shell
docker compose up
```

## Database seeding
The database can be seeded with mock records by adding the --seed flag to the application entrypoint.
This allows for easy functionality testing of the endpoint `/api/v1/chargepoints/location`.

```yaml
services:
  go:
    # [...]
    command: ["/server", "--seed"]
```

The seeder adds mock records to the database. It is intended for development purposes only and should not be used in production.

Seeded data simulates charge points by a set number of vendors, grouped by a logical number of regions. You can alter the amount of records per vendor region, with the environment variable SEED_COUNT. The default value is 100, times 10 regions = 1000 records.


## Testing
At the time of writing, tests have not been implemented yet, but functionality has been tested manually.