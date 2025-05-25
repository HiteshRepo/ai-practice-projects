# Movie Recommendation Assistant

A Go-based movie recommendation system powered by OpenAI's Assistant API with file search capabilities. This application creates an intelligent assistant that can provide personalized movie recommendations based on a curated dataset of movies.

## Features

- **AI-Powered Recommendations**: Uses OpenAI's Assistant API to provide intelligent movie suggestions
- **Vector Search**: Implements vector store for semantic search through movie data
- **File Upload Management**: Automatically handles file uploads and vector store creation
- **Interactive Query Interface**: Command-line interface for asking movie-related questions
- **Persistent Assistant**: Reuses existing assistants and vector stores to optimize performance

## Prerequisites

- Go 1.24.1 or higher
- OpenAI API key with access to Assistant API
- Internet connection for OpenAI API calls

## Installation

1. Clone or download the project
2. Navigate to the assistant directory:
   ```bash
   cd assistant
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Set your OpenAI API key as an environment variable:

```bash
export OPEN_API_KEY="your-openai-api-key-here"
```

## Usage

Run the application with a movie-related query:

```bash
go run main.go -query "Can you recommend a good action movie?"
```

### Example Queries

```bash
# Get action movie recommendations
go run main.go -query "What are some good action movies from 2023?"

# Ask about specific movies
go run main.go -query "Tell me about Oppenheimer"

# Get recommendations by genre
go run main.go -query "I want to watch a comedy movie, what do you suggest?"

# Ask about ratings
go run main.go -query "What are the highest rated movies in your database?"
```

## Project Structure

```
assistant/
├── main.go                 # Main application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── constants/
│   └── constants.go        # Application constants and configuration
├── data/
│   └── movies.txt          # Movie dataset (12 movies with details)
└── openai/
    ├── client.go           # OpenAI client initialization
    ├── assistant.go        # Assistant creation and management
    ├── file.go             # File upload operations
    ├── thread.go           # Thread management for conversations
    └── vectorstore.go      # Vector store operations
```

## How It Works

1. **Initialization**: The application checks for existing uploaded files and vector stores
2. **File Upload**: If not already uploaded, it uploads the movie dataset to OpenAI
3. **Vector Store**: Creates or reuses a vector store for semantic search
4. **Assistant Creation**: Creates or reuses a movie recommendation assistant
5. **Query Processing**: 
   - Creates a conversation thread
   - Adds user query to the thread
   - Runs the assistant with file search capabilities
   - Returns the assistant's response
6. **Cleanup**: Automatically deletes the conversation thread

## Movie Dataset

The application includes a curated dataset of 12 popular movies from 2022-2023:

- **Oppenheimer** (2023) - Biographical Drama - 8.6 rating
- **Top Gun: Maverick** (2022) - Action Drama - 8.3 rating
- **Barbie** (2023) - Adventure Comedy - 7.2 rating
- **The Menu** (2022) - Comedy Thriller - 7.2 rating
- **Elemental** (2023) - Animated Adventure - 7.0 rating
- **Glass Onion** (2022) - Comedy Mystery - 7.1 rating
- **The Super Mario Bros. Movie** (2023) - Animated Adventure - 7.1 rating
- **A Haunting in Venice** (2023) - Crime Drama - 6.8 rating
- **Blue Beetle** (2023) - Action Adventure - 6.7 rating
- **Asteroid City** (2023) - Comedy Drama - 6.6 rating
- **M3GAN** (2022) - Horror/Sci-Fi - 6.4 rating
- **Expend4bles** (2023) - Action - 5.0 rating

Each movie entry includes:
- Title and release year
- MPAA rating
- Runtime
- IMDb rating
- Detailed plot summary
- Director and main cast

## Dependencies

- **github.com/openai/openai-go**: Official OpenAI Go client
- **github.com/caarlos0/env**: Environment variable parsing
- **github.com/pkg/errors**: Enhanced error handling

## API Components

### Assistant Management
- Creates and manages OpenAI assistants
- Configures file search capabilities
- Handles assistant instructions and behavior

### File Operations
- Uploads movie dataset to OpenAI
- Manages file lifecycle and reuse
- Implements custom file reader for uploads

### Vector Store
- Creates vector stores for semantic search
- Associates uploaded files with vector stores
- Optimizes search performance

### Thread Management
- Creates conversation threads
- Manages message flow
- Handles thread cleanup

## Error Handling

The application includes comprehensive error handling for:
- Missing environment variables
- API communication failures
- File upload issues
- Assistant creation problems
- Thread management errors

## Performance Optimizations

- **Resource Reuse**: Checks for existing assistants, files, and vector stores
- **Efficient Cleanup**: Automatically deletes temporary threads
- **Lazy Loading**: Only creates resources when needed
- **Error Recovery**: Graceful handling of API failures

## Limitations

- Movie dataset is limited to 12 curated movies
- Requires active internet connection
- Dependent on OpenAI API availability and quotas
- Non-movie queries will receive a standard "Sorry, I don't know" response

## Contributing

To extend the movie dataset:
1. Add new movie entries to `data/movies.txt`
2. Follow the existing format: `Title: Year | Rating | Duration | Score rating`
3. Include detailed plot summary, director, and cast information

## License

This project is part of a collection of AI practice projects and is intended for educational purposes.
