# Project Scope
There is an abundance of self-hosted CalDav servers out there, yet basically nothing to actually display the data with!

Luna aims to be a *self-hosted CalDav calendar client web application*. That is, Luna will let you connect to CalDav servers, iCal links and more. It will manage all the connections in the background and serve you a view with all your calendars and their events, as you have come to expect from a commercial calendar application.

Ultimately, Luna should become a fully offline-capable *Progressive Web Application (PWA)*, so you can still manage your calendars without an internet connection. Upon reconnecting to the internet, all your events will synchronize with the backend.

This approach ensures out-of-the-box platform cross-compatibility: You will be able to use Luna inside your favourite web browser, as well as install it as a PWA.

Currently, this project is still very much **work-in-progress** and is nowhere near being production-ready. Nevertheless, I have high hopes for one day replacing the propriatary calendar apps I use with Luna. 

Feel free to get in contact through GitHub discussions if you are interested in contributing.

You may also follow the progress in the [development roadmap](https://todo.opisek.net/share/dvEazOyRLEYThqxohVosnqKskYLyoZ4nS8rQ63G1/auth?view=280).

# Deployment
Since Luna is not ready to be used yet, this only serves as instructions on how to get Luna up and running for development purposes!

## Docker
Currently, no pre-made images are available. Instead, you can generate and run the images simply by typing `make` in the root directory of this repository.
Make sure you have **make** and **docker** installed.

## Baremetal
For baremetal deployment, you must ensure your system has:
- **make**
- **node.js** (v20.11.1)
- **go** (go1.23)
- a running **postgres** (version 16) database

For the backend, create an `.env` file in the `backend` directory inside the repository and fill it out accordingly to `.env.example`. To start the backend, run `make` inside the `backend` directory.

Proceed in the same way for the frontend inside the `frontend` directory.

# API
## Current State
The current frontend does not implement all functionality provided by the backend yet. For testing and development purposes, tools like Postman can be used to interact with the API directly. Note that everything is still very much under development and is subject to change.

## General
- All bodies are to be passed as `multipart/form-data`.
- All endpoints except [unauthenticated ones](#unauthenticated) require an access token received from the [Login](#login) endpoint. It is to be passed in the request header as a *bearer token* or as the cookie *token*.
- Parameters passed via the URL are indicated with angular brackets, e.g. `<ID>`

## Design Decisions
### Separation of Concerns
An early draft was to address calendars in explicit relation to their sources and events to their calendars. For example, editing an event would have worked through `/api/sources/<SourceID>/calendars/<CalendarID>/events/<EventID>`. This was later scrapped in favour of the much simpler endpoints `/api/calendars` and `/api/events`.

The reasoning for this change is that Luna is supposed to be a calendar **aggregator** as one of its main principles. After adding one's sources, the user should no longer need to care about them when viewing or manipulating calendars and events. While this simplification of the API makes the implementation of the backend slightly more challenging, it is worth the effort in my opinion.

### UUIDs
Luna uses UUIDs for all its IDs. This has a few reasons:
- Avoiding conflicts in distributed scenarios (future-proofing)
- Better security due to unpredictable IDs; in particular, a potential attacker can neither guess IDs of any resources, nor can they deduct information from IDs, like amount of registered users (since IDs are not consecutive).
- While UUIDv4 is used as a base for some IDs to ensure uniquity and unpredictability, other IDs are built on top of these pseudo-random identifiers using the deterministic UUIDv5. This determinism built on top of random "base" IDs provides design-level collision resistance while maintaining deterministic ways to derive the IDs.

### Local IDs
Luna uses its own (UU)IDs for every resource accessed through it. Therefore, the ID, over which you access a calendar or an event over Luna is different from the underlying IDs used by the upstream source. This has a few reasons:
- Different sources might use different ID types. Luna instead uses the same ID scheme for everything.
- Better security due to hiding the nature of the upstream sources from potential eavesdroppers.
- If an event with the same ID is present in two different calendars (this can have legitimate operational reasons), Luna will still be able to distinguish between them due to different internal IDs.

## Endpoints
### Unauthenticated
#### Login

- **Path**: ``/api/login``
- **Method**: ``POST``
- **Body**: `username`, `password`
- **Purpose**: Returns an authorization token

#### Register
- **Path**: ``/api/register``
- **Method**: ``POST``
- **Body**: `username`, `password`, `email`
- **Purpose**: Creates a new user

#### Version
- **Path**: ``/api/version``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Returns the current backend version. This will be used by the frontend to verify compatibility based on the major version.

### Sources
#### Get Sources
- **Path**: ``/api/sources``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Returns a list of the user's calendar sources.

#### Get Source
- **Path**: ``/api/sources/<ID>``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Returns details for a user's specific source, including authentication data.

#### Put Source
- **Path**: ``/api/sources``
- **Method**: ``PUT``
- **Body**: `name`, `type`, `auth_type`
- **Purpose**: Puts a new calendar source in the database. The authentication information is encrypted by PostegreSQL.

Depending on the `type` field, additional information may need to be passed:
- `caldav`: `url`
- `ical`: Not yet implemented

Depending on the `auth_type` field, additional information may need to be passed:
- `none`: No additional information
- `basic`: `username`, `password`
- `bearer`: `token`
- `oauth`: Not yet implemented

#### Patch Source
- **Path**: ``/api/sources/<ID>``
- **Method**: ``PATCH``
- **Body**: `name`, `type`, `auth_type`, depending on which values should be updated. If `type` and `auth_type` are set, additional information must be provided, as described in the [Put Source](#put-source) endpoint
- **Purpose**: Edit an existing source

#### Delete Source
- **Path**: ``/api/sources/<ID>``
- **Method**: ``DELETE``
- **Body**: Empty
- **Purpose**: Deletes a source from the database.

### Calendars
#### Get Calendars
- **Path**: ``/api/sources/<ID>/calendars``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Fetches calendars from the specified source.

#### Get Calendar
- **Path**: ``/api/calendars/<ID>``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Fetches a specific calendar from its appropriate source.

#### Put Calendar
- **Path**: ``/api/sources/<ID>/calendars``
- **Method**: ``PUT``
- **Body**: `name`, `color`
- **Purpose**: Add a new calendar to the specified source in the upstream, as well as the local database.

#### Patch Calendar
- **Path**: ``/api/calendars/<ID>``
- **Method**: ``PATCH``
- **Body**: `name`, `color`, depending on which values should be updated.
- **Purpose**: Updates specific fields of a calendar in the local database and the upstream source.
- **Note**: This endpoint strives to not erase any values set by other applications that are not supported by Luna.

#### Delete Calendar
- **Path**: ``/api/calendars/<ID>``
- **Method**: ``DELETE``
- **Body**: Empty
- **Purpose**: Deletes the source from the local database and the upstream source.

### Events
#### Get Events
- **Path**: ``/api/calendars/<ID>/events``
- **Method**: ``GET``
- **Search Parameters**: `start`, `end` (both in RFC-3339 format and at most one year apart)
- **Purpose**: Fetches events from the specified calendar.

#### Get Event
- **Path**: ``/api/events/<ID>``
- **Method**: ``GET``
- **Body**: Empty
- **Purpose**: Fetches a specific event from its appropriate calendar.

#### Put Event
- **Path**: ``/api/calendars/<ID>/events``
- **Method**: ``PUT``
- **Body**: `name`, `desc`, `color`, `date_start`, `date_end`, `date_duration`
- **Purpose**: Add a new event to the specified calendar in the upstream, as well as the local database.

The description field is optional. Either the end date or the event duration is to be specified, not both and not neither.

#### Patch Event
- **Path**: ``/api/events/<ID>``
- **Method**: ``PATCH``
- **Body**: `name`, `desc`, `color`, `date_start`, `date_end`, `date_duration`, depending on which values should be updated.
- **Purpose**: Updates specific fields of an event in the local database and the upstream source.
- **Note**: If `desc` should not change, it must be set to its previous values, since leaving it empty implies deleting the description. This endpoint strives to not erase any values set by other applications that are not supported by Luna.

The description field is optional. Either the end date or the event duration is to be specified, not both and not neither.

#### Delete Event
- **Path**: ``/api/events/<ID>``
- **Method**: ``DELETE``
- **Body**: Empty
- **Purpose**: Deletes the event from the local database and the upstream source calendar source.
