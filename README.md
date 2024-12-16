# Calendly Clone
## Product Requirements Document

[Product Requirements Document](./PRD.md)

## Setup

[Setup](./Setup.md)

## System Limitations
Due to the time constraints and the demonstrative nature of the system, the following limitations are present in it:
- Authentication and Authorization are not implemented.
- UI is not implemented.
- Robust input validation and error handling is not implemented. The current system focuses on getting the core functionalities working on their happy-paths.
- Robust logging and observability is not implemented.
- Mutliple timezones are not supported. All times are considered to be in UTC.
- Few edge cases are ignored for the sake of brevity and keeping the system simple. E.g. if an event is booked across multiple days, for instance from 11pm to 1am, the system will not handle it correctly.
- CI/CD is not implemented.

## Code Architecture
![code-architecture](https://github.com/user-attachments/assets/03c0b284-44e7-4a6c-9bd8-06b0adb55414)


### File Structure

```
calendly
├── .
├── .
├── app
│   └── app.go
├── cmd
│   └── rest.go
├── db
│   └── models
│       ├── availability-model.go
│       ├── calendar-model.go
│       ├── event-model.go
│       └── user-model.go
├── pkg
│   ├── dto.go
│   ├── entities.go
│   ├── schedule
│   │   ├── availability.go
│   │   ├── calendar.go
│   │   ├── event.go
│   │   ├── service.go
│   │   └── types.go
│   └── user
│       └── service.go
├── repo
│   ├── id.go
│   ├── pg_repo
│   │   ├── availability-repository.go
│   │   ├── calendar-repository.go
│   │   ├── event-repository.go
│   │   └── user-repository.go
│   └── repositories.go
├── rest
│   ├── availability-handler.go
│   ├── calendar-handler.go
│   ├── event-handler.go
│   ├── types.go
│   └── user-handler.go
└── tests
    └── integration
        ├── calendar_test.go
        ├── down.sh
        ├── mockdata
        │   ├── events.json
        │   └── users.json
        ├── overlap_test.go
        ├── setup.sql
        ├── up.sh
        └── user_test.go
```
