# Vector Embeddings

This project demonstrates how to generate vector embeddings for text data using OpenAI and store them in a Supabase database. It is designed for use cases such as semantic search, recommendation systems, and AI-powered applications that require vector representations of text.

## Features

- Generates embeddings for a set of podcast descriptions using OpenAI's API.
- Stores the embeddings and original text in a Supabase table with vector support.
- Supports semantic search over stored documents using vector similarity.
- Written in Go, using idiomatic project structure and environment-based configuration.

## Project Structure

```
vector-embeddings/
├── .env                # Environment variables (do not commit real secrets)
├── sample.env          # Template for environment variables (safe to share)
├── constants/          # Constants used throughout the project
├── go.mod, go.sum      # Go module files
├── main.go             # Main entry point
├── models/             # Data models (Vector, Document, etc.)
├── openai/             # OpenAI client wrapper
├── queries/            # SQL queries (e.g., table creation, search)
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
- `SUPABASE_API_KEY` — Your Supabase API key (service role or anon, with insert/search permissions)

**Important Security Note:**  
Never commit your real `.env` file or share your actual API keys. Only share `sample.env` as a template.

**To load environment variables:**

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

### 5. Usage

#### Insert Embeddings

To generate embeddings for the sample podcast descriptions and insert them into your Supabase table:

```bash
go run main.go
# or explicitly:
go run main.go -action=insert
```

#### Semantic Search

To perform a semantic search over your documents using a query string:

```bash
go run main.go -action=search -query="What can I listen to in half an hour?"
```

- The application will generate an embedding for your query and return the most similar documents from Supabase, along with similarity scores.

## Code Overview

- **main.go**: Loads environment variables, initializes clients, supports both "insert" and "search" actions.
- **models/**: Contains data models for vectors and documents.
- **openai/**: Wrapper for OpenAI API client.
- **supabase/**: Wrapper for Supabase client and database operations.
- **constants/**: Contains static data (e.g., podcast descriptions).
- **queries/**: SQL scripts for database setup and vector search.

## Environment Variables

See `.env` or `sample.env` for required variables:

- `OPEN_API_KEY`
- `SUPABASE_PROJECT_URL`
- `SUPABASE_API_KEY`

**Never commit your real `.env` file or share your actual API keys.**

## Security

- Do **not** commit your `.env` file or any real API keys to version control.
- Only use `sample.env` as a template for sharing configuration requirements.

## License

MIT
