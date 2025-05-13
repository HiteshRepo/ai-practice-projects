# Vector Embeddings

This project demonstrates how to generate vector embeddings for text data using OpenAI and store them in a Supabase database. It is designed for use cases such as semantic search, recommendation systems, and AI-powered applications that require vector representations of text.

## Features

- Generates embeddings for a set of podcast descriptions using OpenAI's API.
- Stores the embeddings and original text in a Supabase table with vector support.
- Written in Go, using idiomatic project structure and environment-based configuration.

## Project Structure

```
vector-embeddings/
├── .env                # Example environment variables (do not commit secrets)
├── constants/          # Constants used throughout the project
├── go.mod, go.sum      # Go module files
├── main.go             # Main entry point
├── models/             # Data models (Vector, Document, etc.)
├── openai/             # OpenAI client wrapper
├── queries/            # SQL queries (e.g., table creation)
├── sample.env          # Template for environment variables
└── supabase/           # Supabase client and DB operations
```

## Setup

### 1. Clone the Repository

```bash
git clone <repo-url>
cd vector-embeddings
```

### 2. Install Dependencies

Ensure you have Go 1.24+ installed.

```bash
go mod tidy
```

### 3. Configure Environment Variables

Copy `sample.env` to `.env` and fill in your OpenAI and Supabase credentials:

```bash
cp sample.env .env
```

Edit `.env` and set:

- `OPEN_API_KEY` — Your OpenAI API key
- `SUPABASE_PROJECT_URL` — Your Supabase project URL
- `SUPABASE_API_KEY` — Your Supabase API key (service role or anon, with insert permissions)

**Note:** You must load these variables into your environment before running the app:

```bash
source .env
```

### 4. Prepare Supabase

Create the `documents` table in your Supabase project using the SQL in `queries/create_table.sql`:

```sql
create table documents (
  id bigserial primary key,
  content text,
  embedding vector(1536)
);
```

### 5. Run the Application

```bash
go run main.go
```

This will generate embeddings for the sample podcast descriptions and insert them into your Supabase table.

## Code Overview

- **main.go**: Loads environment variables, initializes clients, generates embeddings, and inserts them into Supabase.
- **models/**: Contains data models for vectors and documents.
- **openai/**: Wrapper for OpenAI API client.
- **supabase/**: Wrapper for Supabase client and database operations.
- **constants/**: Contains static data (e.g., podcast descriptions).
- **queries/**: SQL scripts for database setup.

## Environment Variables

See `.env` or `sample.env` for required variables:

- `OPEN_API_KEY`
- `SUPABASE_PROJECT_URL`
- `SUPABASE_API_KEY`

## License

MIT
