from openai import OpenAI
client = OpenAI()

def ask_llm(prompt):
    response = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[{"role": "user", "content": prompt}],
        temperature=0.7
    )
    return response.choices[0].message.content

question = "What did our Q3 report say about APAC revenue?"

print(ask_llm(question))

context = """
Q3 Report Summary:
- North America revenue grew by 12%
- Europe revenue declined by 3%
"""

prompt = f"""
Context:
{context}

Question: What did our Q3 report say about APAC revenue?
Answer:
"""

print(ask_llm(prompt))

def grounded_prompt(context, question):
    return f"""
You are a helpful assistant. Answer the user's question using ONLY the information provided.

If the answer is not in the context, say:
"I don't have enough information to answer this."

Context:
---------
{context}
---------

Question: {question}

Answer:
"""

print(ask_llm(grounded_prompt(context, question)))

def citation_prompt(context, question):
    return f"""
Answer the question using ONLY the context below.

For every claim, include an exact quote from the context.

If the answer is not present, say:
"I don't have enough information."

Context:
---------
{context}
---------

Question: {question}

Answer (with quotes):
"""

print(ask_llm(citation_prompt(context, question)))


retrieved_chunks = [
    "Company overview and mission statement",
    "Employee satisfaction survey results"
]

context = "\n".join(retrieved_chunks)

print(ask_llm(grounded_prompt(context, question)))

def faithfulness_check(context, answer):
    prompt = f"""
Given this context:
{context}

And this answer:
{answer}

Does the answer contain claims not supported by the context?
Reply YES or NO and explain.
"""
    return ask_llm(prompt)


# Simulate a bad answer
bad_answer = "APAC revenue grew by 8% in Q3."

print(faithfulness_check(context, bad_answer))

from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity

docs = [
    "Q3 Report: North America revenue grew 12%",
    "Q3 Report: Europe revenue declined 3%",
    "Q2 Report: APAC revenue grew 8%"
]

def retrieve(query, docs):
    vectorizer = TfidfVectorizer()
    vectors = vectorizer.fit_transform([query] + docs)
    scores = cosine_similarity(vectors[0:1], vectors[1:]).flatten()
    return [docs[i] for i in scores.argsort()[-2:]]

retrieved = retrieve(question, docs)
context = "\n".join(retrieved)

print("Retrieved:\n", context)
print(ask_llm(grounded_prompt(context, question)))