# Hugging Face Inference UI

This project demonstrates how to use the Hugging Face Inference API for text-to-speech conversion with a simple web UI.

## CORS Issue Fix

The original implementation had CORS issues because it was trying to call the Hugging Face API directly from the browser. This has been fixed by:

1. Creating a Node.js/Express server that acts as a proxy between the UI and the Hugging Face API
2. Modifying the UI code to use the proxy endpoint instead of directly calling the API
3. Adding proper CORS headers to the server

## Setup

1. Install dependencies:
   ```
   npm install
   ```

2. Set your Hugging Face API token as an environment variable:
   ```
   export HF_TOKEN=your_hugging_face_token
   ```

## Running the Application

Start the server:
```
npm run server
```

This will:
1. Start the Express server on port 3000 (or the port specified in the PORT environment variable)
2. Serve the static files from the `ui` directory
3. Create an API endpoint at `/api/text-to-speech` that proxies requests to the Hugging Face API

Then open your browser and navigate to:
```
http://localhost:3000
```

## How It Works

1. The UI makes a POST request to the local server endpoint `/api/text-to-speech` with the text and model parameters
2. The server forwards the request to the Hugging Face API using your API token
3. The server receives the audio data from the Hugging Face API and sends it back to the UI
4. The UI creates a blob URL from the audio data and sets it as the source of the audio element

This approach avoids CORS issues because:
- The UI is served from the same origin as the API endpoint
- The server handles the cross-origin request to the Hugging Face API
- The server adds the necessary CORS headers to allow the UI to access the API
