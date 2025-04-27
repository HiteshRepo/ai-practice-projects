# Hugging Face Inference Demo

A Node.js project demonstrating the use of the [Hugging Face Inference API](https://huggingface.co/docs/huggingface.js/inference/README) for chat completion, sentiment classification, and translation tasks.

## Features

- **Chat Completion**: Generate conversational responses using large language models.
- **Text Classification**: Perform sentiment analysis and emotion detection on text.
- **Translation**: Translate text between languages using state-of-the-art models.

## Requirements

- Node.js v18 or higher
- A Hugging Face API token ([create one here](https://huggingface.co/settings/tokens))

## Installation

1. Clone this repository or copy the `hugging-face-inference` directory.
2. Install dependencies:

   ```bash
   npm install
   ```

3. Create a `.env` file in the project root with your Hugging Face token:

   ```
   HF_TOKEN=your_huggingface_token_here
   ```

## Usage

Run the script with the `-task` flag to specify which functionality to use.

### 1. Chat Completion

Generates a conversational response to a prompt.

```bash
node index.js -task chat-completion
```

**Example Output:**
```
{
  role: 'assistant',
  content: 'Inference refers to the process of applying a trained machine learning model to make predictions, classifications, or decisions based on new and unseen data...'
}
```

### 2. Sentiment Classification

Classifies text as positive or negative, and detects nuanced emotions.

```bash
node index.js -task classify
```

**Example Output:**
```
positive
negative
[
  { label: 'sadness', score: 0.77 },
  { label: 'surprise', score: 0.12 },
  ...
]
```

### 3. Translation

Translates a sample English sentence to Hindi.

```bash
node index.js -task translate
```

**Example Output:**
```
{ translation_text: 'एक एआई इंजीनियर होने के लिए एक रोमांचक समय है' }
```

## Models Used

- **Chat Completion**: `HuggingFaceH4/zephyr-7b-beta`
- **Sentiment Classification**: `cardiffnlp/twitter-roberta-base-sentiment-latest`
- **Emotion Detection**: `j-hartmann/emotion-english-distilroberta-base`
- **Translation**: `facebook/mbart-large-50-many-to-many-mmt`

## Project Structure

```
.
├── .env                # Environment variables (your HF_TOKEN)
├── index.js            # Main script
├── package.json        # Project metadata and dependencies
├── package-lock.json   # Dependency lock file
└── node_modules/       # Installed dependencies
```

## Notes

- Ensure your Hugging Face token has the necessary permissions for inference.
- You can modify the sample texts and models in `index.js` to experiment with other tasks or models.

## License

ISC

---

**Author:** hiteshrepo
