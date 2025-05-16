# Movie Chatbot

A command-line chatbot that recommends movies and answers movie-related questions using OpenAI's GPT-4 and Supabase for semantic search. The bot maintains conversation history and can answer follow-up questions based on previous interactions.

## Features

- **Movie Recommendations:** Get personalized movie suggestions based on your queries.
- **Follow-up Questions:** Ask about previous recommendations or your own details (e.g., "What is my name?").
- **Contextual Memory:** Remembers conversation history to provide relevant answers.
- **Semantic Search:** Uses vector embeddings and Supabase to find the best-matching movie context.
- **Powered by OpenAI GPT-4:** Generates natural, concise responses.

## Example Conversation

```
> Hi I am hitesh. Recommend me a action movie.
2025/05/16 21:19:38 Hi Hitesh, I recommend you watch "The Expendables." It's filled with high-stakes assignments and shows the unity and courage of the team while displaying their unmatched combat skills.
> Which movie is good for children?
2025/05/16 21:20:00 The "Super Mario Bros. Movie" would be a good choice for children. It's an animated adventure comedy voiced by popular actors such as Chris Pratt and Anya Taylor-Joy.
> What is my name?
2025/05/16 21:20:12 Your name is Hitesh.
> Which movies have been recommended by you till now?
2025/05/16 21:20:34 Based on the conversation history, I have recommended "The Expendables" and "The Super Mario Bros. Movie".
```

## Getting Started

### Prerequisites

- Go 1.24+
- OpenAI API key
- Supabase project and API key

### Setup

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd movie-chatbot
   ```

2. **Configure environment variables:**

   Create a `.env` file in the `movie-chatbot/` directory with the following content:
   ```
   export OPEN_API_KEY=your-openai-api-key
   export SUPABASE_PROJECT_URL=your-supabase-url
   export SUPABASE_API_KEY=your-supabase-api-key
   ```

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Run the chatbot:**
   ```sh
   go run main.go
   ```

## Usage

- Type your questions or requests at the prompt (`>`).
- Press `Ctrl+C` to exit.

## Project Structure

- `main.go` - Entry point and main chatbot logic.
- `constants/` - Constants used throughout the project.
- `models/` - Data models for vectors and database documents.
- `openai/` - OpenAI API client integration.
- `supabase/` - Supabase client and vector search logic.

## License

MIT License
