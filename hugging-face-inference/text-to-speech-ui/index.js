const text = "It's an exciting time to be an A.I. engineer."
const model = "facebook/mms-tts"

// Use the proxy endpoint instead of directly calling the Hugging Face API
async function generateSpeech() {
    try {
        const response = await fetch('/api/text-to-speech', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                text,
                model,
            }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const audioBlob = await response.blob();
        const audioElement = document.getElementById('speech');
        const speechUrl = URL.createObjectURL(audioBlob);
        audioElement.src = speechUrl;
    } catch (error) {
        console.error('Error generating speech:', error);
        alert('Failed to generate speech. See console for details.');
    }
}

// Call the function when the page loads
generateSpeech();
