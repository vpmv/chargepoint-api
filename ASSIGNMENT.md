Go Coding Assignment: EV Charge Point API
---

# Goal

The goal of this assignment is to build a simple RESTful API server in Go that manages and
queries a database of Electric Vehicle (EV) charge points. The entire assignment should take
approximately 2 hours of coding. We value clean, simple, and idiomatic Go. The solution will
serve as the starting point for a technical discussion during our interview. Usage of LLMs is
allowed, but be wary that you yourself need to be able to explain decisions that you (or the LLM)
made during the interview.

# Requirements
Your project must be a self-contained Go application. When run, it should start an HTTP server
and be ready to accept requests.

## Data Model
Define a ChargePoint struct. It should contain at least the following information:
- ID: A unique identifier.
- Name: A human-readable name.
- Location: Coordinates of the charge point.
- Status: The current status of the charge point (e.g. AVAILABLE, OCCUPIED, OFFLINE).

## Database
- Use a database solution of your choice to store the charge points.

## API Server
Your server must expose the following JSON API endpoints:

- POST `/chargepoints`
  - **Action**: Creates a new charge point.
  - **Request** Body: A JSON object representing a ChargePoint.
  - **Response**: The newly created ChargePoint object (including its ID) as JSON.
  - **Validation**: Perform basic validation on the input data.
- GET `/chargepoints/{id}`
  - **Action**: Retrieves a single charge point by its ID.
  - **Response**: The corresponding ChargePoint object as JSON.
  - **Error**: Return a 404 Not Found if the ID does not exist.
- GET `/chargepoints/nearby`
  - **Action**: Finds all charge points within a given radius of a location.
  - **Query** Parameters:
     - lat: The latitude of the search center.
     - lon: The longitude of the search center.
     - radius: The search radius in kilometers (km).
  - **Response**: A JSON array of the ChargePoint objects that are within the radius.