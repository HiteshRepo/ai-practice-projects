# Content Moderator

A command-line tool that checks if text content is safe or contains potentially harmful material using OpenAI's moderation API.

## Overview

Content Moderator is a Go application that:
1. Takes text input from various sources (currently single-line text, with planned support for multi-line text, URLs, and files)
2. Sends the content to OpenAI's moderation API
3. Analyzes the response to determine if the content is safe
4. If unsafe, identifies specific categories of harmful content (harassment, hate, self-harm, sexual, violence, illicit)

## Prerequisites

- Go 1.24 or later
- API key for:
  - [OpenAI](https://platform.openai.com/) - for content moderation

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd content-moderator
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Configuration

1. Create a `.env` file in the project root directory based on the provided `sample.env`:
   ```
   cp sample.env .env
   ```

2. Edit the `.env` file and add your API key:
   ```
   OPEN_API_KEY=your_openai_api_key_here
   ```

   - Get your OpenAI API key from: https://platform.openai.com/settings/organization/api-keys

## Usage

### Basic Usage

Run the application with a single line of text:

```
go run main.go -content "This is some text to check"
```

### Content Types

The application supports different types of content input (with some planned for future implementation):

```
# Single-line text (default)
go run main.go -type=single_line_text -content "This is some text to check"

# Multi-line text (planned)
go run main.go -type=multi_line_text -content "This is some text to check\nThis is a second line"

# Text from URL (planned)
go run main.go -type=text_from_url -url "https://example.com/article"

# Text from file (planned)
go run main.go -type=text_from_file -file "path/to/file.txt"
```

### Building the Application

To build an executable:

```
go build -o content-moderator
```

Then run it:

```
./content-moderator -content "This is some text to check"
```

## How It Works

1. The application takes text input from the user through command-line flags.
2. It sends this content to OpenAI's moderation API, which analyzes the text for potentially harmful content.
3. The API returns a response indicating whether the content is flagged as unsafe and, if so, in which specific categories.
4. The application processes this response and outputs a message indicating whether the content is safe or unsafe.
5. If the content is unsafe, it lists the specific categories of harmful content that were detected.

## Example Output

### Safe Content:
```
go run main.go -content "The sky is blue and the weather is nice today."

Your content is safe
```

### Unsafe Content:
```
go run main.go -content "I want to hurt myself."

Your content is flagged as unsafe in categories: [Self Harm Self Harm Intent]
```

## Limitations

- Currently only supports single-line text input (support for multi-line text, URLs, and files is planned).
- Relies on OpenAI's moderation API, which may have rate limits depending on your subscription plan.
- The moderation API may not catch all forms of harmful content or may flag content that is actually safe in certain contexts.
