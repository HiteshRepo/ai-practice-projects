import numpy as np
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity, euclidean_distances

model = SentenceTransformer('all-MiniLM-L6-v2')

# Simulates a RAG scenario: one short query vs docs of varying length but similar meaning
query = "What are the benefits of exercise?"

docs = {
    "short (same meaning)": "Exercise improves health and mood.",
    "long (same meaning)": (
        "Regular physical activity has numerous well-documented benefits including "
        "improved cardiovascular health, better mental health, improved mood, "
        "increased energy levels, stronger muscles, and better sleep quality."
    ),
    "unrelated": "The stock market closed higher yesterday.",
}

query_emb = model.encode([query])
doc_embs  = model.encode(list(docs.values()))

print(f"Query: \"{query}\"\n")
print(f"{'Doc':<22} {'Dot Product':>12} {'Euclidean':>12} {'Cosine':>10}")
print("-" * 60)

for i, (label, _) in enumerate(docs.items()):
    dot = float(np.dot(query_emb[0], doc_embs[i]))
    euc = float(euclidean_distances(query_emb, [doc_embs[i]])[0][0])
    cos = float(cosine_similarity(query_emb, [doc_embs[i]])[0][0])
    print(f"{label:<22} {dot:>12.4f} {euc:>12.4f} {cos:>10.4f}")

print("""
Why cosine wins for RAG:
  - Dot product is skewed by vector magnitude (longer docs score higher regardless of relevance)
  - Euclidean distance is also sensitive to magnitude — same idea phrased verbosely looks "far"
  - Cosine measures the angle between vectors, capturing semantic direction regardless of length
  => A short query reliably retrieves both short and long docs with the same meaning
""")

# Output
# Query: "What are the benefits of exercise?"

# Doc                     Dot Product    Euclidean     Cosine
# ------------------------------------------------------------
# short (same meaning)         0.7055       0.7674     0.7055
# long (same meaning)          0.7528       0.7031     0.7528
# unrelated                   -0.0273       1.4334    -0.0273

# Why cosine wins for RAG:
#   - Dot product is skewed by vector magnitude (longer docs score higher regardless of relevance)
#   - Euclidean distance is also sensitive to magnitude — same idea phrased verbosely looks "far"
#   - Cosine measures the angle between vectors, capturing semantic direction regardless of length
#   => A short query reliably retrieves both short and long docs with the same meaning