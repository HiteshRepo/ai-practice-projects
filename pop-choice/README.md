# Pop Choice

Pop Choice is a command-line tool written in Go that recommends movies based on user interests using OpenAI embeddings and Supabase as a vector database. It supports both single-user and multi-user modes, allowing you to get personalized movie recommendations for yourself or a group.

## Features

- **AI-powered Recommendations:** Uses OpenAI's embedding API to understand user interests and match them to movies.
- **Supabase Integration:** Stores and queries movie embeddings using Supabase as a vector database.
- **Single & Multi-user Modes:** Get recommendations for one person or collaboratively for a group.
- **Interactive CLI:** Simple command-line interface for entering interests and receiving recommendations.

## Project Structure

```
pop-choice/
├── .env                  # Environment variables (see below)
├── constants/            # Constants and question lists
├── go.mod, go.sum        # Go module files and dependencies
├── main.go               # Main CLI application
├── models/               # Data models (e.g., vector)
├── openai/               # OpenAI client integration
├── queries/              # SQL queries for Supabase
└── supabase/             # Supabase client and operations
```

## Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/yourrepo.git
cd yourrepo/pop-choice
```

### 2. Install Dependencies

Ensure you have Go 1.24+ installed.

```bash
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the `pop-choice` directory with the following variables:

```
OPEN_API_KEY=your-openai-api-key
SUPABASE_PROJECT_URL=your-supabase-project-url
SUPABASE_API_KEY=your-supabase-api-key
```

**Note:** Do not commit your `.env` file with real credentials.

### 4. Run the Application

#### Setup (populate Supabase with movie embeddings):

```bash
go run main.go -action=setup
```

#### Single User Mode

```bash
go run main.go -action=single-user
```

#### Multi User Mode

```bash
go run main.go -action=multi-user
```

Follow the prompts to enter interests and receive movie recommendations.

## Environment Variables

| Variable                | Description                        |
|-------------------------|------------------------------------|
| `OPEN_API_KEY`          | Your OpenAI API key                |
| `SUPABASE_PROJECT_URL`  | Your Supabase project URL          |
| `SUPABASE_API_KEY`      | Your Supabase API key              |

## Dependencies

- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [Supabase Go](https://github.com/nedpals/supabase-go)
- [caarlos0/env](https://github.com/caarlos0/env) for env parsing

## License

MIT
