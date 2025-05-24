# Roketin Backend Engineer Test Case: Challenge 2

## Case

- API to create and upload movies. Required information related with a movies are at least  title, description, duration, artists, genres
- API to update movie
- API to list all movies with pagination
- API to search movie by title/description/artists/genres

## How To Run

1. Install dependencies
   ```bash
   go mod tidy
   ```

2. Run the application
   ```bash
   go run .
   ```

3. Open Swagger UI API documentation at `http://localhost:8080/swagger/index.html`

## API Endpoints

### Create Movie
- **POST** `/movies`
- Body: 
  ```json
  {
    "title": "Movie Title",
    "description": "Movie description",
    "duration": 120,
    "artists": ["Artist 1", "Artist 2"],
    "genres": ["Action", "Drama"]
  }
  ```

### Update Movie
- **PUT** `/movies/{id}`
- Body: Same as create movie

### List Movies (with pagination)
- **GET** `/movies`
- Query params:
  - page: current pagination page, default: 1
  - limit: max movies per page, default: 10

### Search Movies
- **GET** `/movies/search`
- Query params:
  - page: current pagination page, default: 1
  - limit: max movies per page, default: 10
  - title, optional
  - description, optional
  - artists, optional
  - genres, optional

## Example

```bash
# Create a movie
curl -X POST http://localhost:8080/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "description": "A computer programmer discovers reality is a simulation",
    "duration": 136,
    "artists": ["Keanu Reeves", "Laurence Fishburne"],
    "genres": ["Sci-Fi", "Action"]
  }'

# Get all movies with pagination
curl "http://localhost:8080/movies?page=1&limit=5"

# Search movies
curl "http://localhost:8080/movies/search?title=matrix"

# Update a movie
curl -X PUT http://localhost:8080/movies/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix Reloaded",
    "description": "Updated description",
    "duration": 138,
    "artists": ["Keanu Reeves", "Carrie-Anne Moss"],
    "genres": ["Sci-Fi", "Action"]
  }'
```