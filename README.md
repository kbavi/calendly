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
[Postman Collection](./Calendly.postman_collection.json)

## DB Schema
![DB Schema](https://github.com/user-attachments/assets/e8aaafe8-3110-4939-af95-37b301d4934f)

## Code Architecture
![code-architecture](https://github.com/user-attachments/assets/03c0b284-44e7-4a6c-9bd8-06b0adb55414)
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
