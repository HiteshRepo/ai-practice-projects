from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

model = SentenceTransformer('all-MiniLM-L6-v2')

sentences = [
    "I love dogs",
    "I adore puppies",
    "Stock markets fell"
]

embeddings = model.encode(sentences)

for i in range(len(sentences)):
    print(f"'{sentences[i]}' [{embeddings[i].shape[0]}d]: {embeddings[i][:6].round(4)}...")
print()

similarity_matrix = cosine_similarity(embeddings)

for i in range(len(sentences)):
    for j in range(len(sentences)):
        print(f"Similarity between '{sentences[i]}' and '{sentences[j]}': {similarity_matrix[i][j]:.4f}")
    print()

# Output
# Similarity between 'I love dogs' and 'I love dogs': 1.0000
# Similarity between 'I love dogs' and 'I adore puppies': 0.6831
# Similarity between 'I love dogs' and 'Stock markets fell': 0.0197

# Similarity between 'I adore puppies' and 'I love dogs': 0.6831
# Similarity between 'I adore puppies' and 'I adore puppies': 1.0000
# Similarity between 'I adore puppies' and 'Stock markets fell': 0.0349

# Similarity between 'Stock markets fell' and 'I love dogs': 0.0197
# Similarity between 'Stock markets fell' and 'I adore puppies': 0.0349
# Similarity between 'Stock markets fell' and 'Stock markets fell': 1.0000