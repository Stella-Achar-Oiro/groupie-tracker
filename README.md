Groupie Tracker is a web application that displays information about music artists and their concert dates. It fetches data from an external API and presents it in an interactive user interface.

## Project Structure
groupie-tracker/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── handlers.go
│   ├── models/
│   │   └── models.go
│   ├── services/
│   │   └── api.go
│   └── routes/
│       └── routes.go
├── static/
│   ├── index.html
│   ├── styles.css
│   └── script.js
├── tests/
│   └── api_test.go
├── go.mod
└── README.md

## Prerequisites

- Go (version 1.16 or later)
- Git

## Setup

1. Clone the repository:
git clone https://github.com/Stella-Achar-Oiro/groupie-tracker.git
cd groupie-tracker

2. Install dependencies:
go mod tidy

## Running the Application

1. Start the server:
go run cmd/server/main.go
Copy
2. Open a web browser and navigate to `http://localhost:8080`

## Features

- Display a grid of music artists
- Search functionality to filter artists
- Detailed view for each artist including:
- Band members
- Creation date
- First album release date
- Concert locations and dates
- Interactive map showing concert locations

## API Endpoints

- `GET /api/data`: Fetches all data (artists, locations, dates, relations)
- `GET /api/search?q=<query>`: Searches for artists based on the provided query

## Running Tests

Execute the following command from the project root:
go test ./tests

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- Data provided by the Groupie Trackers API
- Maps powered by Leaflet.js and OpenStreetMap