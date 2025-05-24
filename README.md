# AI Practice Projects

This repository contains a collection of AI-powered applications and tools, including both Go and Node.js projects, that leverage APIs such as OpenAI and Hugging Face for various practical use cases.

## Projects

### [Hugging Face Inference](./hugging-face-inference/README.md)

A Node.js project demonstrating the use of the Hugging Face Inference API for chat completion, sentiment classification, and translation tasks.

**Key Features:**
- Chat completion using large language models
- Sentiment and emotion classification
- Multilingual translation
- Simple CLI interface with task-based flags
- Easy setup with Hugging Face API token

### [Motivational Speaker](./motivational-speaker/README.md)

A command-line tool that generates positive and encouraging responses to user inputs using a fine-tuned OpenAI model.

**Key Features:**
- Uploads training data to OpenAI
- Creates and manages a fine-tuning job for GPT-3.5-turbo
- Uses the fine-tuned model to generate motivational responses
- Provides uplifting and encouraging messages to help users overcome challenges

### [Stock Price Predictor](./stock-price-predictor/README.md)

A command-line tool that generates concise stock performance reports and buy/hold recommendations based on recent price data.

**Key Features:**
- Fetches stock price data from Polygon.io
- Analyzes data using OpenAI's GPT-4o Mini model
- Supports zero-shot and few-shot analysis approaches
- Generates concise reports with buy/hold recommendations

### [Image Generator](./image-generator/README.md)

A command-line tool that generates images from text descriptions using OpenAI's DALL-E models.

**Key Features:**
- Connects to OpenAI's DALL-E 3 API
- Generates high-quality images based on text prompts
- Saves generated images as PNG files locally
- Supports customizable image parameters (size, style)
- Optional Studio Ghibli-style image generation

### [Content Moderator](./content-moderator/README.md)

A command-line tool that checks if text content is safe or contains potentially harmful material using OpenAI's moderation API.

**Key Features:**
- Analyzes text input using OpenAI's moderation API
- Identifies specific categories of harmful content
- Supports multiple input types: single-line text, multi-line text, URLs, and files
- Provides clear feedback on content safety

### [Pollyglot](./pollyglot/README.md)

A command-line tool that translates any given content into a specified language using OpenAI's GPT-3.5-turbo model.

**Key Features:**
- Translates text to any target language using OpenAI's GPT-3.5-turbo
- Simple CLI interface with flags for content and target language
- Secure API key management via environment variables
- Easily extensible and written in idiomatic Go

### [Vector Embeddings](./vector-embeddings/README.md)

A Go project for generating, storing, and searching vector embeddings for text data (e.g., podcast descriptions, movie details) using OpenAI and Supabase. Supports semantic search, recommendations, and AI-powered applications requiring vector representations of text.

**Key Features:**
- Generate embeddings for text data (podcasts, movies) using OpenAI's API.
- Store embeddings and original text in Supabase tables with vector support.
- Perform semantic search over stored documents or movies using vector similarity.
- Chat-based Q&A over documents or movies using OpenAI GPT-4.
- Written in Go with idiomatic project structure and environment-based configuration.

**Setup:**
1. Clone the repository and enter the project directory:
   ```bash
   git clone <repo-url>
   cd vector-embeddings
   ```
2. Install dependencies (requires Go 1.24+):
   ```bash
   go mod tidy
   ```
3. Copy `sample.env` to `.env` and fill in your OpenAI and Supabase credentials:
   ```bash
   cp sample.env .env
   ```
   - `OPEN_API_KEY` — Your OpenAI API key
   - `SUPABASE_PROJECT_URL` — Your Supabase project URL
   - `SUPABASE_API_KEY` — Your Supabase API key (service role or anon, with insert/search permissions)

4. Prepare Supabase by creating the required tables using the SQL scripts in `vector-embeddings/queries/`:
   - `create_documents_table.sql`
   - `create_movies_table.sql`

**Usage Examples:**
- Insert podcast embeddings:
  ```bash
  go run main.go -action=insert-docs
  ```
- Semantic search (podcasts):
  ```bash
  go run main.go -action=search-docs -query="What can I listen to in half an hour?"
  ```
- Semantic search + chat (podcasts):
  ```bash
  go run main.go -action=search-n-chat-docs -query="An episode Elon Musk would enjoy"
  ```
- Insert movie embeddings (chunked):
  ```bash
  go run main.go -action=chunk-n-insert-movies
  ```
- Semantic search + chat (movies):
  ```bash
  go run main.go -action=query-movie -query="Which movie can I take my child to?" -matches=3
  ```

**Command-Line Flags:**
- `-action` (required): One of `insert-docs`, `search-docs`, `search-n-chat-docs`, `chunk-n-insert-movies`, `query-movie`
- `-query` (required for search/chat actions): Query string for semantic search
- `-matches` (optional for `query-movie`): Number of top matches to return (default: 1)

**Security Note:**  
Never commit your real `.env` file or share your actual API keys. Only share `sample.env` as a template.

### [Movie Chatbot](./movie-chatbot/README.md)

A command-line chatbot that recommends movies and answers movie-related questions using OpenAI's GPT-4 and Supabase for semantic search. The bot maintains conversation history and can answer follow-up questions based on previous interactions.

**Key Features:**
- Movie recommendations based on user queries
- Answers follow-up and context-based questions (e.g., "What is my name?", "Which movies have you recommended?")
- Remembers conversation history for personalized responses
- Uses OpenAI GPT-4 and Supabase for semantic search and context retrieval
- Simple CLI interface

### [Pop Choice](./pop-choice/README.md)

A command-line tool that recommends movies based on user interests using OpenAI embeddings and Supabase as a vector database. Supports both single-user and multi-user modes for personalized or group recommendations.

**Key Features:**
- AI-powered movie recommendations using OpenAI embeddings
- Stores and queries movie embeddings with Supabase
- Single-user and multi-user interactive CLI modes
- Easy setup and environment-based configuration

### [reAct](./reAct/README.md)

A Go-based command-line tool that leverages OpenAI's API to generate personalized activity ideas based on your current location and weather. Supports multiple versions of suggestion logic (`v1`, `v2`, and `v3`), with `v3` as the default, and is easily extensible.

**Key Features:**
- AI-powered activity suggestions using OpenAI
- Weather and location-aware recommendations
- Switchable logic versions (`v1`, `v2`, `v3`) with `v3` as the default
- Simple CLI interface with flags for version and query
- Environment-based configuration

**Setup:**
1. Enter the project directory:
   ```bash
   cd reAct
   ```
2. Install dependencies (requires Go 1.24+):
   ```bash
   go mod tidy
   ```
3. Create a `.env` file and add your API keys:
   ```
   OPEN_API_KEY=your-openai-api-key
   WEATHER_STACK_API_KEY=your-weatherstack-api-key
   ```
   Or export them in your shell before running the CLI.

**Security Note:**  
Never commit your real `.env` file or share your actual API keys. Only share `.env` as a template if needed.

## Common Technologies

All projects are built with:
- Go 1.24+ and/or Node.js 18+
- API integration (OpenAI, Hugging Face)
- Environment-based configuration
- Command-line interfaces

## Getting Started

Each project has its own README with detailed instructions for setup and usage. To get started with a specific project:

1. Navigate to the project directory
2. Follow the installation instructions in the project's README
3. Configure the required API keys
4. Run the application

## Repository Structure

```
ai-practice-projects/
├── README.md                      # This file
├── .gitignore                     # Git ignore file
├── hugging-face-inference/        # Hugging Face Inference (Node.js)
│   ├── README.md                  # Project documentation
│   ├── index.js                   # Main application code
│   ├── .env                       # Environment variables template
│   └── package.json               # Project metadata
├── motivational-speaker/          # Motivational Speaker project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── finetunedata.jsonl         # Training data for fine-tuning
│   ├── sample.env                 # Environment variables template
│   └── openai/                    # OpenAI client implementation
│       └── client.go
├── stock-price-predictor/         # Stock Price Predictor project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── sample.env                 # Environment variables template
│   ├── openai/                    # OpenAI client implementation
│   │   └── client.go
│   └── polygon/                   # Polygon.io client implementation
│       ├── client.go
│       └── models.go
├── image-generator/               # Image Generator project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── sample.env                 # Environment variables template
│   └── openai/                    # OpenAI client implementation
│       └── client.go
├── content-moderator/             # Content Moderator project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── sample.env                 # Environment variables template
│   └── openai/                    # OpenAI client implementation
│       └── client.go
├── pollyglot/                     # Pollyglot translation project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── sample.env                 # Environment variables template
│   └── openai/                    # OpenAI client implementation
│       └── client.go
├── vector-embeddings/             # Vector Embeddings project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── sample.env                 # Environment variables template
│   ├── openai/                    # OpenAI client implementation
│   ├── supabase/                  # Supabase client and DB operations
│   ├── constants/                 # Constants and sample data
│   ├── models/                    # Data models
│   └── queries/                   # SQL scripts for DB setup
├── movie-chatbot/                 # Movie Chatbot project
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── .env                       # Environment variables template
│   ├── go.mod                     # Go module definition
│   ├── go.sum                     # Go dependencies lockfile
│   ├── constants/                 # Project constants
│   ├── models/                    # Data models
│   ├── openai/                    # OpenAI client implementation
│   └── supabase/                  # Supabase client and vector search
├── pop-choice/                    # Pop Choice movie recommender
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── .env                       # Environment variables template
│   ├── go.mod                     # Go module definition
│   ├── go.sum                     # Go dependencies lockfile
│   ├── constants/                 # Project constants
│   ├── models/                    # Data models
│   ├── openai/                    # OpenAI client implementation
│   ├── supabase/                  # Supabase client and DB operations
│   └── queries/                   # SQL scripts for DB setup
├── reAct/                         # reAct activity suggestion CLI (Go + OpenAI)
│   ├── README.md                  # Project documentation
│   ├── main.go                    # Main application code
│   ├── .env                       # Environment variables template
│   ├── go.mod                     # Go module definition
│   ├── go.sum                     # Go dependencies lockfile
│   ├── constants/                 # Project constants
│   ├── models/                    # Data models
│   ├── openai/                    # OpenAI client implementation
│   ├── tools/                     # Utility tools
│   ├── utils/                     # Helper utilities
│   └── versions/                  # Versioned suggestion logic
└── prompts/                       # Additional prompt resources
```

## Prerequisites

- Go 1.24 or later
- API keys for the services used by each project:
  - [OpenAI](https://platform.openai.com/) - for all projects
  - [Polygon.io](https://polygon.io/) - for Stock Price Predictor
  - OpenAI DALL-E access - for Image Generator

## License

[TODO]

## Contributing

[TODO]
