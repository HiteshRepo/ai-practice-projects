from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
from sentence_transformers import SentenceTransformer
from sentence_transformers import CrossEncoder

model = SentenceTransformer("all-MiniLM-L6-v2")
reranker = CrossEncoder("cross-encoder/ms-marco-MiniLM-L-6-v2")

documents = [
    "Go routines enable concurrent programming in Go.",
    "Python supports multiprocessing and threading.",
    "Databricks is used for big data analytics and machine learning.",
    "Vector databases store embeddings for similarity search.",
    "BM25 is a ranking function used in search engines.",
    "Transformers provide state-of-the-art NLP models."
]

query = "How does concurrency work in Go?"

print("Query: \n",query)
print()
print("Docs: \n", documents)
print()

vectorizer = TfidfVectorizer()
doc_vectors = vectorizer.fit_transform(documents)
query_vector = vectorizer.transform([query])

scores = cosine_similarity(query_vector, doc_vectors).flatten()

sparse_results = sorted(
    list(zip(documents, scores)),
    key=lambda x: x[1],
    reverse=True
)

print("\n--- Sparse Retrieval (TF-IDF) ---")
for doc, score in sparse_results[:3]:
    print(f"{score:.4f} - {doc}")


doc_embeddings = model.encode(documents)
query_embedding = model.encode([query])

scores = cosine_similarity(query_embedding, doc_embeddings).flatten()

dense_results = sorted(
    list(zip(documents, scores)),
    key=lambda x: x[1],
    reverse=True
)

print("\n--- Dense Retrieval (Embeddings) ---")
for doc, score in dense_results[:3]:
    print(f"{score:.4f} - {doc}")


alpha = 0.5  # weight between sparse and dense

hybrid_scores = alpha * scores + (1 - alpha) * cosine_similarity(query_vector, doc_vectors).flatten()

hybrid_results = sorted(
    list(zip(documents, hybrid_scores)),
    key=lambda x: x[1],
    reverse=True
)

print("\n--- Hybrid Retrieval ---")
for doc, score in hybrid_results[:3]:
    print(f"{score:.4f} - {doc}")

# Take top-k from hybrid
top_k_docs = [doc for doc, _ in hybrid_results[:5]]

pairs = [(query, doc) for doc in top_k_docs]
rerank_scores = reranker.predict(pairs)

reranked_results = sorted(
    list(zip(top_k_docs, rerank_scores)),
    key=lambda x: x[1],
    reverse=True
)

print("\n--- Re-ranked Results ---")
for doc, score in reranked_results:
    print(f"{score:.4f} - {doc}")

# Output
# Query: 
#  How does concurrency work in Go?

# Docs: 
#  ['Go routines enable concurrent programming in Go.', 'Python supports multiprocessing and threading.', 'Databricks is used for big data analytics and machine learning.', 'Vector databases store embeddings for similarity search.', 'BM25 is a ranking function used in search engines.', 'Transformers provide state-of-the-art NLP models.']


# --- Sparse Retrieval (TF-IDF) ---
# 0.7017 - Go routines enable concurrent programming in Go.
# 0.2010 - BM25 is a ranking function used in search engines.
# 0.0000 - Python supports multiprocessing and threading.

# --- Dense Retrieval (Embeddings) ---
# 0.8303 - Go routines enable concurrent programming in Go.
# 0.3284 - Python supports multiprocessing and threading.
# 0.1527 - Databricks is used for big data analytics and machine learning.

# --- Hybrid Retrieval ---
# 0.7660 - Go routines enable concurrent programming in Go.
# 0.1642 - Python supports multiprocessing and threading.
# 0.1622 - BM25 is a ranking function used in search engines.

# --- Re-ranked Results ---
# 2.1895 - Go routines enable concurrent programming in Go.
# -11.0048 - Python supports multiprocessing and threading.
# -11.2840 - Vector databases store embeddings for similarity search.
# -11.3016 - BM25 is a ranking function used in search engines.
# -11.4187 - Databricks is used for big data analytics and machine learning.