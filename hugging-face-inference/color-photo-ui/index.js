const model = "ghoskno/Color-Canny-Controlnet-model"

// Use the proxy endpoint instead of directly calling the Hugging Face API
async function colorPhoto() {
    try {
        const response = await fetch('/api/color-photo', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                model,
            }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const newImageBase64 = await blobToBase64(newImageBlob)
        const newImage = document.getElementById("new-image")
        newImage.src = newImageBase64
    } catch (error) {
        console.error('Error coloring photo:', error);
        alert('Failed to color photo. See console for details.');
    }
}

// Call the function when the page loads
colorPhoto();
