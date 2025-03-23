# Motivational Speaker

A command-line tool that generates positive and encouraging responses to user inputs using a fine-tuned OpenAI model.

## Overview

Motivational Speaker is a Go application that:
1. Uploads training data to OpenAI (if not already uploaded)
2. Creates and manages a fine-tuning job for GPT-3.5-turbo
3. Uses the fine-tuned model to generate motivational responses
4. Provides uplifting and encouraging messages to help users overcome challenges

## Prerequisites

- Go 1.24 or later
- API key for:
  - [OpenAI](https://platform.openai.com/) - for fine-tuning and generating responses

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd motivational-speaker
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

Run the application:

```
go run main.go
```

The application will:
1. Check if the training data has been uploaded to OpenAI
2. Upload the data if necessary
3. Check if a fine-tuning job has been created and completed
4. Create a new fine-tuning job if necessary
5. Use the fine-tuned model to generate a motivational response to a sample input

### Customizing the Application

To modify the application for different use cases:

1. Edit the `finetunedata.jsonl` file to include different training examples
2. Modify the user message in `main.go` to test different inputs:
   ```go
   messages := []openai.ChatCompletionMessageParamUnion{
       openai.UserMessage("Your custom message here"),
   }
   ```

### Building the Application

To build an executable:

```
go build -o motivational-speaker
```

Then run it:

```
./motivational-speaker
```

## How It Works

1. The application first checks if the training data file (`finetunedata.jsonl`) has been uploaded to OpenAI.
2. If the file hasn't been uploaded, it uploads the file to OpenAI's servers.
3. It then checks if a fine-tuning job has been created for this training data.
4. If no job exists or the previous job failed, it creates a new fine-tuning job.
5. The application monitors the status of the fine-tuning job until it completes successfully.
6. Once the fine-tuned model is ready, it uses the model to generate a motivational response to a user input.
7. The response is displayed in the console.

## Training Data

The application uses a JSONL file (`finetunedata.jsonl`) containing examples of motivational responses to various user inputs. Each example includes:

- A system message defining the bot's role as a motivational assistant
- A user message expressing a concern, doubt, or negative feeling
- An assistant message providing a positive, encouraging response

You can customize this file to include different examples and train the model to respond in different ways.

## Example Output

```
2025/03/23 13:37:00 training data already uploaded file-abc123
2025/03/23 13:37:05 fine tune job retrieved ft-xyz789
2025/03/23 13:37:05 fine tune job in succeed state
2025/03/23 13:37:05 going to use fine tuned model ft:gpt-3.5-turbo-0125:personal::xyz789
2025/03/23 13:37:10 Life is full of purpose and meaning, even when it doesn't feel that way. Your unique talents, passions, and experiences all contribute to your personal journey. Start by exploring what brings you joy and fulfillment. Set small, achievable goals that align with your values. Remember, finding purpose is a process, not a destination. Each step you take brings clarity. You have incredible potential waiting to be discovered. Trust yourself and keep moving forward – your purpose will reveal itself along the way!
```

## Limitations

- Fine-tuning OpenAI models requires credits on your OpenAI account.
- The fine-tuning process may take several minutes to complete.
- The quality of the responses depends on the quality and variety of the training examples provided.
- The application currently only handles a single hardcoded user input. For a more interactive experience, you would need to modify the code to accept user input from the command line.
