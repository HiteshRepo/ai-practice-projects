import time

import faiss
import numpy as np
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("all-MiniLM-L6-v2")

corpus = [
    # health & exercise
    "Exercise improves cardiovascular health and mood.",
    "Regular physical activity strengthens muscles and boosts energy.",
    "A daily walk can significantly reduce stress levels.",
    "Weight training helps prevent bone density loss.",
    "Swimming is a low-impact full-body workout.",
    # nutrition
    "A balanced diet is essential for good health.",
    "Eating vegetables daily reduces the risk of chronic disease.",
    "Processed foods are linked to increased inflammation.",
    # finance
    "Stock markets closed higher yesterday.",
    "The Federal Reserve raised interest rates again.",
    "Inflation is affecting consumer spending patterns.",
    "Investing early maximizes compound interest over time.",
    # technology & AI
    "Machine learning models require large amounts of training data.",
    "Neural networks are inspired by the structure of the human brain.",
    "Transformers revolutionized natural language processing.",
    "Cloud computing has transformed software infrastructure.",
    # nature & climate
    "Climate change is accelerating glacier melting worldwide.",
    "Rainforests are home to more than half of Earth's species.",
    "Ocean temperatures are rising due to greenhouse gas emissions.",
    "Wildfires are becoming more frequent and severe globally.",
]

query = "What exercises are good for building strength?"

# Encode and normalize (cosine similarity = dot product on normalized vectors)
print("Encoding corpus...")
embeddings = model.encode(corpus, normalize_embeddings=True).astype("float32")
dim = embeddings.shape[1]

# Tile to simulate a large-scale database
SCALE       = 500
large_db    = np.tile(embeddings, (SCALE, 1))
num_vectors = large_db.shape[0]
print(f"Database: {num_vectors:,} vectors × {dim}d\n")

query_emb = model.encode([query], normalize_embeddings=True).astype("float32")

K = 5

# ── Exact Search (brute-force numpy) ─────────────────────────────────────────
t0           = time.perf_counter()
scores       = large_db @ query_emb[0]
exact_idx    = np.argsort(-scores)[:K]
exact_time   = time.perf_counter() - t0

# ── ANN Search (FAISS HNSW) ───────────────────────────────────────────────────
# HNSW builds a navigable small-world graph; query traverses a fraction of nodes
index = faiss.IndexHNSWFlat(dim, 32)   # 32 edges per node
index.hnsw.efConstruction = 64         # depth during build (quality vs build time)
index.hnsw.efSearch        = 32        # depth during search (quality vs query time)

t0          = time.perf_counter()
index.add(large_db)
build_time  = time.perf_counter() - t0

t0          = time.perf_counter()
_, ann_raw  = index.search(query_emb, K)
ann_time    = time.perf_counter() - t0

ann_idx = ann_raw[0]

# ── Print results ─────────────────────────────────────────────────────────────
def sentence(idx: int) -> str:
    return corpus[idx % len(corpus)]

print(f'Query: "{query}"\n')

print("Exact Search (brute-force):")
for rank, i in enumerate(exact_idx, 1):
    print(f"  {rank}. [{scores[i]:.4f}] {sentence(i)}")
print(f"  Time: {exact_time * 1000:.2f}ms\n")

print("ANN Search (FAISS HNSW):")
for rank, i in enumerate(ann_idx, 1):
    print(f"  {rank}. {sentence(i)}")
print(f"  Index build: {build_time * 1000:.2f}ms")
print(f"  Query time:  {ann_time * 1000:.2f}ms\n")

exact_sentences = {sentence(i) for i in exact_idx}
ann_sentences   = {sentence(i) for i in ann_idx}
overlap         = exact_sentences & ann_sentences

print(f"Overlap: {len(overlap)}/{K} ({len(overlap) / K * 100:.0f}% recall)")
print(f"Speedup: {exact_time / ann_time:.1f}x faster at query time\n")

print("""
Key insight:
  - Exact search is O(n) — scans every vector, guaranteed correct results
  - HNSW builds a graph index; queries traverse only a small subgraph
  - At 10k+ vectors ANN is orders of magnitude faster with ~95-99% recall
  - RAG systems use ANN because exact search over millions of docs is too slow
""")


# Output
# Encoding corpus...
# Database: 10,000 vectors × 384d

# Query: "What exercises are good for building strength?"

# Exact Search (brute-force):
#   1. [0.4900] Regular physical activity strengthens muscles and boosts energy.
#   2. [0.4900] Regular physical activity strengthens muscles and boosts energy.
#   3. [0.4900] Regular physical activity strengthens muscles and boosts energy.
#   4. [0.4900] Regular physical activity strengthens muscles and boosts energy.
#   5. [0.4900] Regular physical activity strengthens muscles and boosts energy.
#   Time: 10.21ms

# ANN Search (FAISS HNSW):
#   1. Regular physical activity strengthens muscles and boosts energy.
#   2. Regular physical activity strengthens muscles and boosts energy.
#   3. Regular physical activity strengthens muscles and boosts energy.
#   4. Regular physical activity strengthens muscles and boosts energy.
#   5. Regular physical activity strengthens muscles and boosts energy.
#   Index build: 9407.97ms
#   Query time:  0.71ms

# Overlap: 1/5 (20% recall)
# Speedup: 14.5x faster at query time


# Key insight:
#   - Exact search is O(n) — scans every vector, guaranteed correct results
#   - HNSW builds a graph index; queries traverse only a small subgraph
#   - At 10k+ vectors ANN is orders of magnitude faster with ~95-99% recall
#   - RAG systems use ANN because exact search over millions of docs is too slow