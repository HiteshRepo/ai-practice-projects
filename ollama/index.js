import ollama from "ollama";
import express from "express";
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const port = 3000;

app.use(cors());
app.use(express.json());
app.use(express.static(path.join(__dirname, 'ui')));

app.get('/', async (req, res) => {
    const question = req.query.question;
    if (!question) {
        res.status(200).send("Ask something via the `?question=` parameter");  
    } else {
        const response = await ollama.chat({
            model: 'qwen2.5:1.5b',
            messages: [{ role: 'user', content: question }],
        });
        res.status(200).send(response.message.content);
    }
});

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});