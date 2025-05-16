# Vector Embeddings

A Go project for generating, storing, and searching vector embeddings for text data (e.g., podcast descriptions, movie details) using OpenAI and Supabase. Supports semantic search, recommendations, and AI-powered applications requiring vector representations of text.

## Features

- Generate embeddings for text data (podcasts, movies) using OpenAI's API.
- Store embeddings and original text in Supabase tables with vector support.
- Perform semantic search over stored documents or movies using vector similarity.
- Chat-based Q&A over documents using OpenAI GPT-4.
- Written in Go with idiomatic project structure and environment-based configuration.

## Project Structure

```
vector-embeddings/
├── .env                # Environment variables (do not commit real secrets)
├── sample.env          # Template for environment variables (safe to share)
├── constants/          # Static data (e.g., podcast descriptions, table names)
├── go.mod, go.sum      # Go module files
├── main.go             # Main entry point and CLI
├── models/             # Data models (Vector, Document, etc.)
├── openai/             # OpenAI API client wrapper
├── queries/            # SQL scripts for DB setup and vector search
├── supabase/           # Supabase client and DB operations
└── langchain/          # Document chunking and text processing utilities
```

## Setup

### 1. Clone the Repository

```bash
git clone <repo-url>
cd vector-embeddings
```

### 2. Install Dependencies

Requires Go 1.24+.

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

**Important:**  
Never commit your real `.env` file or share your actual API keys. Only share `sample.env` as a template.

To load environment variables:

```bash
source .env
```

### 4. Prepare Supabase

Create the required tables in your Supabase project using the SQL in `queries/`:

#### Documents Table

```sql
-- queries/create_documents_table.sql
create table documents (
  id bigserial primary key,
  content text,
  embedding vector(1536)
);
```

#### Movies Table

```sql
-- queries/create_movies_table.sql
create table movies (
  id bigserial primary key,
  content text,
  embedding vector(1536)
);
```

## Usage

Run the application with different actions using the `-action` flag. Some actions also require `-query` and/or `-matches`.

### Actions

#### 1. Insert Podcast Embeddings

Generate embeddings for sample podcast descriptions and insert them into the `documents` table.

```bash
go run main.go -action=insert-docs
```

#### 2. Semantic Search (Podcasts)

Perform a semantic search over podcast documents using a query string.

```bash
go run main.go -action=search-docs -query="What can I listen to in half an hour?"
```

- Returns the most similar documents from Supabase, with similarity scores.

#### 3. Semantic Search + Chat (Podcasts)

Performs a semantic search and then uses GPT-4 to answer the query based on the most relevant document.

```bash
go run main.go -action=search-n-chat-docs -query="An episode Elon Musk would enjoy"
```

#### 4. Insert Movie Embeddings (Chunked)

Splits movie details into chunks, generates embeddings, and inserts them into the `movies` table.

```bash
go run main.go -action=chunk-n-insert-movies
```

#### 5. Semantic Search + Chat (Movies)

Performs a semantic search over movies and uses GPT-4 to answer the query based on the top matches.

```bash
go run main.go -action=query-movie -query="Which movie can I take my child to?" -matches=3
```

- `-matches` (optional, default: 1): Number of top matches to use for the answer.

### Command-Line Flags

- `-action` (required): One of `insert-docs`, `search-docs`, `search-n-chat-docs`, `chunk-n-insert-movies`, `query-movie`
- `-query` (required for search/chat actions): Query string for semantic search
- `-matches` (optional for `query-movie`): Number of top matches to return (default: 1)

## Code Overview

- **main.go**: CLI entry point, parses flags, dispatches actions.
- **constants/**: Static data (podcast descriptions, table names, etc.).
- **langchain/**: Text chunking utilities for large documents.
- **models/**: Data models for vectors and database rows.
- **openai/**: OpenAI API client wrapper for embeddings and chat.
- **supabase/**: Supabase client and database operations.
- **queries/**: SQL scripts for table creation and vector search.

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
