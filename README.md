# Calendly Clone
## Product Requirements Document

[Product Requirements Document](./PRD.md)

## Installation and Setup

[Installation and Setup](./Setup.md)

## System Limitations
Due to the time constraints and the demonstrative nature of the system, the following limitations are present in it:
- Authentication and Authorization are not implemented.
- UI is not implemented.
- Robust input validation and error handling is not implemented. The current system focuses on getting the core functionalities working on their happy-paths.
- Robust logging and observability is not implemented.
- Mutliple timezones are not supported. All times are considered to be in UTC to avoid complexity.
- Few edge cases are ignored for the sake of brevity and keeping the system simple. E.g. if an event is booked across multiple days, for instance from 11pm to 1am, the system will not handle it correctly.
- CI/CD is not implemented.

## API Documentation
Complete API Documentation is available in [Postman Collection](./Calendly.postman_collection.json)

### Highlighted API Routes

**Get Calendar**
```sh
GET /api/v1/calendars/{calendar_id}
```
- This route is used to get the calendar details.
- supports filtering between time range.
- It returns the
  - availability rules
  - events
  - free intervals
  - scheduling links (e.g. `/book/:calendar_id/slots/30-mins`)


**Set Availability**
```sh
POST /api/v1/availabilities
```

This route is used to set the availability rules for a calendar. Currently, the route supports setting one or more blocks of availability for a day of the week.
e.g:
- Monday 9am to 1pm and 2pm to 5pm
- Tuesday 10am to 12pm

**Get Overlaps Between Calendars**
```sh
GET /api/v1/calendars/availabilities/overlap
```
- accepts list of calendar ids
- accept time range to find overlaps between the calendars
- accounts for events booked in the time range in one or more of the calendars
- returns overlaps intervals - free time for all the participants

**Booking URL**
```sh
GET /book/:calendar_id/slots/30-mins?from=2024-12-14T09:00&to=2024-12-21T18:00
```
- accepts calendar id
- accepts slot duration
- accepts time range to find free slots
- returns free time intervals of duration in the given time range for the given calendar
- accounts for events booked in the time range in the calendars and avoids double booking

## DB Schema
![DB Schema](https://github.com/user-attachments/assets/e8aaafe8-3110-4939-af95-37b301d4934f)

### Models

**UserModel**
```go
type UserModel struct {
	ID    string
	Email string
	Name  string
}
```

**AvailabilityModel**
```go
type AvailabilityModel struct {
	ID         string
	CalendarID string
	Rules      string
}
```
- `Rules` is a JSON encoded string column that stores the availability rules.
- Rules are currently decoded and processed at the application layer as type `pkg.AvailabilityRule` inside [entities file](./pkg/entities.go).
- Rule supports two types of availability rules:
  - Day: e.g. Monday to Friday, 9am to 5pm
  - Date: 2024-01-01, 10am to 12pm (business logic is not implemented for date based rules)

**CalendarModel**
```go
type CalendarModel struct {
	ID             string
	UserID         string
	AvailabilityID *string
}
```
- `AvailabilityID` is a foreign key that references the `AvailabilityModel`.

**EventModel**
```go
type EventModel struct {
	ID          string
	CalendarID  string
	Title       string
	Description *string
	Start       time.Time
	Ending      time.Time
	Invitees    string // comma separated invitee emails
}
```
- `Invitees` is a comma separated string column that stores the invitee emails.
- `Start` and `Ending` are the start and end times of the event.
- `CalendarID` is a foreign key that references the `CalendarModel`.

## Code Architecture
![code-architecture](https://github.com/user-attachments/assets/3f95453c-f9b0-4310-83f0-9482534e5b46)

This codebase provides a monolithic architecture with a single executable binary. The code on a high-level is divided into 3 layers:

- **Presentation Layer**
  - Handles HTTP requests and responses.
  - Exposes the REST API.
  - Contains the handlers for the REST API.
  - Can be used to expose other interfaces like GraphQL, gRPC, etc.
  - Talks to the application layer to implement the business logic.

- **Application Layer**
  - Acts as a unified api interface for the presentation layer.
  - Contains the business logic.
  - Contains the services that implement the business logic.
  - Abstracts the data layer from the presentation layer.
  - Designed to be easily testable.
  - Designed to be easily extendable.
  - Talks to the data layer to get and store data.

- **Data Layer**
  - Contains the data models.
  - Contains the repository interfaces.
  - Contains the repository implementations.
  - Uses [repository pattern](https://www.umlboard.com/design-patterns/repository.html) to abstract the data layer from the application layer.
  - Can be used to implement calendar solutions like Google Calendar, Microsoft Calendar, etc.
  - Can be extended to support async data sources like Kafka and SQS.

### File Structure
```
calendly
├── .
├── .
├── app
│   └── app.go                        # application
├── cmd
│   └── rest.go                       # rest api server
├── db
│   └── models
│       ├── availability-model.go
│       ├── calendar-model.go
│       ├── event-model.go
│       └── user-model.go
├── pkg
│   ├── dto.go                        # data transfer objects
│   ├── entities.go                   # shareable entities between layers
│   ├── schedule
│   │   ├── availability-service.go
│   │   ├── calendar-service.go
│   │   ├── event-service.go
│   │   ├── service.go                # service interface
│   │   └── types.go
│   └── user
│       └── service.go
├── repo
│   ├── id.go                         # id generator
│   ├── pg_repo
│   │   ├── availability-repository.go
│   │   ├── calendar-repository.go
│   │   ├── event-repository.go
│   │   └── user-repository.go
│   └── repositories.go               # repository interfaces
├── rest                              # rest api handlers
│   ├── availability-handler.go
│   ├── calendar-handler.go
│   ├── event-handler.go
│   ├── types.go
│   └── user-handler.go
└── tests
    └── integration
        ├── calendar_test.go            # calendar integration tests
        ├── down.sh                     # down script
        ├── mockdata
        │   ├── events.json
        │   └── users.json
        ├── overlap_test.go             # overlap integration tests
        ├── setup.sql                   # db setup script
        ├── up.sh                       # up script
        └── user_test.go                # user integration tests
```
## Integration Tests
Follow the [Setup and Installation Document](./Setup.md) to run the integration tests.

- Runs test cases for the test suite 6.1 and 6.2 mentioned in the [PRD](./PRD.md).
- Other test cases are not implemented.
- Uses [Testify](https://pkg.go.dev/github.com/stretchr/testify) for assertions and mocks.

## Further Improvements
- Refactor scheduling logic to it's own package to enable versioning and independant testing.
- Implement complete test cases mentioned in the [PRD](./PRD.md).
- Implement authentication and authorization.
- Implement basic UI to demo the system.
- Handle graceful server shutdowns.
- Handle scheduling edge-cases.
- Handle unhappy paths in business logic.
- Introduce robust error handling.
- Implement logging and observability.
- Implement CI/CD.
- Add support for multiple timezones.
- Add support for async notifications.
- Add support for google and microsoft calendar integrations.
- Add support for multiple calendars per user.
- Add CI/CD pipeline.