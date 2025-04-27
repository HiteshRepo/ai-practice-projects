import { HfInference } from '@huggingface/inference'

const taskFlagIdx = process.argv.indexOf('-task');

if (taskFlagIdx == -1) {
    console.log("flag task not provided")
    process.exit(1) 
}

if (process.argv.length <= taskFlagIdx+1) {
    console.log("flag task value not provided")
    process.exit(1) 
}

const task = process.argv[taskFlagIdx+1]


// Create your Hugging Face Token: https://huggingface.co/settings/tokens
// Set your Hugging Face Token: https://scrimba.com/dashboard#env
// Learn more: https://scrimba.com/links/env-variables
const hf = new HfInference(process.env.HF_TOKEN)

// Hugging Face Inference API docs: https://huggingface.co/docs/huggingface.js/inference/README

if (task == "chat-completion") {
    const textToGenerate = "The definition of machine learning inference is "
    const model = "HuggingFaceH4/zephyr-7b-beta"

    const response = await hf.chatCompletion({
        messages: 
        [
            {
                content: textToGenerate,
                role: "user"
            }
        ],
        model
    })

    console.log(response.choices[0].message)
    // {
    //     role: 'assistant',
    //     content: 'Inference refers to the process of applying a trained machine learning model to make predictions, classifications, or decisions based on new and unseen data. This is in contrast to the training process, which involves feeding large amounts of labeled data into the model to learn how to make accurate predictions. During inference, the model uses the insights it gained during training to quickly and accurately predict an outcome based on input data. Inference can be performed on a single machine or distributed across a network of computers in parallel, allowing for faster and more efficient computations on large datasets.'
    //   }
}

if (task == "classify") {
    const positiveTextToClassify = "I just bought a new camera. It's the best camera I've ever owned!"
    const negativeTextToClassify = "I just bought a new camera. It's been a real disappointment."

    const model = "cardiffnlp/twitter-roberta-base-sentiment-latest"

    let response = await hf.textClassification({
        model,
        inputs: positiveTextToClassify
    })

    console.log(response[0].label) // positive

    response = await hf.textClassification({
        model,
        inputs: negativeTextToClassify
    })

    console.log(response[0].label) // negative

    const modelForNuancedResponse = "j-hartmann/emotion-english-distilroberta-base"

    response = await hf.textClassification({
        model: modelForNuancedResponse,
        inputs: negativeTextToClassify
    })

    console.log(response)   
    // [
    //     { label: 'sadness', score: 0.7744333148002625 },
    //     { label: 'surprise', score: 0.12495166063308716 },
    //     { label: 'neutral', score: 0.051190637052059174 },
    //     { label: 'disgust', score: 0.023046262562274933 },
    //     { label: 'joy', score: 0.013148864731192589 },
    //     { label: 'anger', score: 0.010783486068248749 },
    //     { label: 'fear', score: 0.0024457769468426704 }
    //   ] 
}

if (task == "translate") {
    const textToTranslate = "It's an exciting time to be an AI engineer"

    const model = "facebook/mbart-large-50-many-to-many-mmt"

    let response = await hf.translation({
        inputs: textToTranslate,
        model,
        parameters: {
            src_lang: "en_XX",
            tgt_lang: "hi_IN",
        }
    })

    console.log(response) // { translation_text: 'एक एआई इंजीनियर होने के लिए एक रोमांचक समय है' }
}