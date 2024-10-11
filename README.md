# music-library
EM test task

This repository contains the implementation of an online song library API, offering functionalities such as adding, retrieving, updating and deleting songs.

## Features:

### RESTful API:
- GET /songs: Retrieves a list of songs, allowing filtering by any field and pagination.
- GET /songs/{songId}/lyrics: Retrieves the lyrics of a specific song, paginated by verses.
- DELETE /songs/{songId}: Deletes a song by its ID.
- PUT /songs/{songId}: Updates the information of a song by its ID.
- POST /songs: Adds a new song to the library. The request body should be in the following JSON format:
```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}
```
#### External API Integration:
 When adding a new song, the API makes a request to an external API to enrich the song information.  

#### PostgreSQL Database:
 The enriched song information is stored in a PostgreSQL database, with the database schema defined using migrations during service startup. To start db conatainer:  
run `make db` command  

#### Getting Started:

#### Clone the repository:

`git clone https://github.com/existanz/music-library`

#### Install dependencies:
`go mod download`

#### Configure the environment:

Copy the `.env.example` file to `.env` and fill in the necessary configuration parameters, including database credentials and external API endpoint.  
#### Run the API:

`make run`

### Access the API documentation:

Open
http://localhost:4001/docs/index.html
in your browser.