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

A Go project that generates vector embeddings for text data (e.g., podcast descriptions) using OpenAI and stores them in a Supabase database with vector support. Supports both inserting new embeddings and performing semantic search over stored documents.

**Key Features:**
- Generates embeddings for a set of podcast descriptions using OpenAI's API
- Stores embeddings and original text in a Supabase table with vector support
- Supports semantic search over stored documents using vector similarity
- Includes SQL for table setup and idiomatic Go project structure
- Environment-based configuration for easy setup
- Secure handling of API keys via environment variables (never commit real secrets)

**Usage Examples:**
- Insert sample embeddings:
  ```bash
  cd vector-embeddings
  go run main.go
  # or explicitly:
  go run main.go -action=insert
  ```
- Perform semantic search:
  ```bash
  go run main.go -action=search -query="What can I listen to in half an hour?"
  ```

**Security Note:**  
Never commit your real `.env` file or share your actual API keys. Use `sample.env` as a template.

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
