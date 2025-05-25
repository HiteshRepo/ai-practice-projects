# Bill Price Calculator

A Go application that leverages OpenAI's GPT-4 Vision API to analyze images and answer questions about them. This project demonstrates multimodal AI capabilities by processing images and providing intelligent responses to queries.

## Features

- **Image Analysis**: Uses OpenAI's GPT-4 Vision model to analyze and understand image content
- **Multiple Use Cases**: Supports various image analysis scenarios including:
  - Food comparison (cheese types)
  - Menu price calculation for bulk orders
- **Base64 Image Encoding**: Automatically encodes images for API transmission
- **Environment Configuration**: Secure API key management through environment variables

## Project Structure

```
bill-price-calculator/
├── .env                    # Environment variables (API key)
├── go.mod                  # Go module dependencies
├── go.sum                  # Go module checksums
├── main.go                 # Main application logic
├── images/                 # Sample images for testing
│   ├── cheese-1.jpeg      # First cheese sample
│   ├── cheese-2.jpeg      # Second cheese sample
│   └── menu.png           # Restaurant menu sample
└── openai/                # OpenAI client package
    └── client.go          # OpenAI client initialization
```

## Dependencies

- **OpenAI Go SDK**: `github.com/openai/openai-go v1.1.0`
- **Environment Parser**: `github.com/caarlos0/env v3.5.0+incompatible`
- **Error Handling**: `github.com/pkg/errors v0.9.1`

## Setup

### Prerequisites

- Go 1.24.1 or later
- OpenAI API key with GPT-4 Vision access

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd bill-price-calculator
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
# Copy the .env file and update with your OpenAI API key
export OPEN_API_KEY=your_openai_api_key_here
```

## Usage

### Running the Application

```bash
# Source environment variables
source .env

# Run the application
go run main.go
```

### Use Cases

The application currently demonstrates two main use cases:

1. **Cheese Comparison**: 
   - Query: "What's the difference between these two types of cheese?"
   - Images: `cheese-1.jpeg`, `cheese-2.jpeg`

2. **Menu Price Calculation**:
   - Query: "I want to order one of each item on this menu for my company party. How much would that cost?"
   - Image: `menu.png`

### Adding Custom Use Cases

You can extend the application by modifying the `UseCases` map in `main.go`:

```go
var UseCases = map[string][]string{
    "Your question here": {"path/to/image1.jpg", "path/to/image2.jpg"},
    // Add more use cases as needed
}
```

## Code Structure

### Main Components

- **`main()`**: Entry point that processes all use cases
- **`VisualizeImage()`**: Handles OpenAI API communication for image analysis
- **`base64ImageFile()`**: Converts image files to base64 encoding
- **`NewOpenAiClient()`**: Initializes the OpenAI client with API key

### Key Features

- **Error Handling**: Comprehensive error handling with detailed logging
- **Flexible Image Support**: Supports multiple image formats (JPEG, PNG)
- **Batch Processing**: Processes multiple use cases in sequence
- **Modular Design**: Separated OpenAI client logic for reusability

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OPEN_API_KEY` | OpenAI API key for GPT-4 Vision access | Yes |

## API Usage

The application uses OpenAI's Chat Completions API with the `gpt-4o` model, which supports vision capabilities. Images are encoded in base64 format and sent as part of the chat completion request.

## Example Output

```
Usecase: What's the difference between these two types of cheese?, Response: [AI analysis of cheese differences]

Usecase: I want to order one of each item on this menu for my company party. How much would that cost?, Response: [AI calculation of total menu costs]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Security Notes

- Never commit your `.env` file with real API keys
- Keep your OpenAI API key secure and rotate it regularly
- Monitor your OpenAI API usage to avoid unexpected charges

## Troubleshooting

### Common Issues

1. **API Key Error**: Ensure your OpenAI API key is correctly set in the environment
2. **Image Not Found**: Verify image paths are correct relative to the project root
3. **Module Not Found**: Run `go mod download` to install dependencies

### Support

For issues and questions, please open an issue in the repository or refer to the [OpenAI API documentation](https://platform.openai.com/docs).
