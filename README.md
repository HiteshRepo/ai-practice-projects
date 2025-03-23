# AI Practice Projects

This repository contains a collection of Go applications that leverage AI capabilities through OpenAI's API for various practical use cases.

## Projects

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

## Common Technologies

Both projects are built with:
- Go 1.24+
- OpenAI API integration
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
└── prompts/                       # Additional prompt resources
```

## Prerequisites

- Go 1.24 or later
- API keys for the services used by each project:
  - [OpenAI](https://platform.openai.com/) - for all projects
  - [Polygon.io](https://polygon.io/) - for Stock Price Predictor

## License

[TODO]

## Contributing

[TODO]
