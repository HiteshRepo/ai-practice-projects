import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { HfInference } from '@huggingface/inference';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, 'ui')));

const hf = new HfInference(process.env.HF_TOKEN);

app.post('/api/text-to-speech', async (req, res) => {
  try {
    const { text, model } = req.body;
    
    if (!text || !model) {
      return res.status(400).json({ error: 'Text and model are required' });
    }
    
    const response = await hf.textToSpeech({
      inputs: text,
      model,
    });
    
    const buffer = await response.arrayBuffer();
    const audioBuffer = Buffer.from(buffer);
    
    res.set('Content-Type', 'audio/wav');
    res.set('Content-Length', audioBuffer.length);
    
    res.send(audioBuffer);
  } catch (error) {
    console.error('Error in text-to-speech:', error);
    res.status(500).json({ error: error.message });
  }
});

app.post('/api/color-photo', async (req, res) => {
  try {
    const { model } = req.body;
    
    if (!model) {
      return res.status(400).json({ error: 'model is required' });
    }

    const oldImagePath = "./color-photo-ui/black-n-white.jpg"
    const oldImageDescription = `An elderly couple walks together on a gravel path with green 
grass and trees on each side. Wearing neutral-colored clothes, they face away
 from the camera as they carry their bags.`

    const blackAndWhiteImageResponse = await fetch(oldImagePath)
    const blackAndWhiteImageBlob = await blackAndWhiteImageResponse.blob()
    
    const newImageBlob = await hf.imageToImage({
      inputs: blackAndWhiteImageBlob,
      model,
      parameters: {
        prompt: oldImageDescription,
        negative_prompt: "Black and white photo. text, bad anatomy, blurry, low quality",
        strength: 0.85,
      },
    });

    const imageArrayBuffer = await newImageBlob.arrayBuffer();
    const imageBuffer = Buffer.from(imageArrayBuffer);
    
    res.set('Content-Type', 'image/jpg');
    res.set('Content-Length', imageBuffer.length);

    res.send(imageBuffer);
  } catch (error) {
    console.error('Error in color-photo:', error);
    res.status(500).json({ error: error.message });
  }
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
