import atexit
import subprocess
import time

from sentence_transformers import SentenceTransformer
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct

CONTAINER = "qdrant-snippet-03"
PORT      = 6333

# ── Docker lifecycle ──────────────────────────────────────────────────────────

def start_qdrant() -> QdrantClient:
    subprocess.run(["docker", "rm", "-f", CONTAINER], capture_output=True)
    subprocess.Popen(
        ["docker", "run", "--rm", "--name", CONTAINER, "-p", f"{PORT}:6333", "qdrant/qdrant"],
        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
    )
    atexit.register(lambda: subprocess.run(["docker", "stop", CONTAINER], capture_output=True))

    client = QdrantClient(host="localhost", port=PORT)
    for _ in range(30):
        try:
            client.get_collections()
            print("Qdrant ready.\n")
            return client
        except Exception:
            time.sleep(0.5)
    raise RuntimeError("Qdrant did not start in time")

# ── Demo ──────────────────────────────────────────────────────────────────────

client = start_qdrant()
model  = SentenceTransformer("all-MiniLM-L6-v2")

docs = [
    "Exercise improves cardiovascular health and mood.",
    "Regular physical activity strengthens muscles and boosts energy.",
    "A balanced diet is essential for good health.",
    "Stock markets closed higher yesterday.",
    "The Federal Reserve raised interest rates again.",
]

query = "What are the benefits of working out?"

embeddings = model.encode(docs)
dim        = embeddings.shape[1]

if not client.collection_exists("demo"):
    client.create_collection(
        collection_name="demo",
        vectors_config=VectorParams(size=dim, distance=Distance.COSINE),
    )

client.upsert(
    collection_name="demo",
    points=[PointStruct(id=i, vector=embeddings[i].tolist(), payload={"text": doc})
            for i, doc in enumerate(docs)],
)

query_vec = model.encode([query])[0].tolist()
results   = client.query_points(collection_name="demo", query=query_vec, limit=3).points

print(f"Query: \"{query}\"\n")
print(f"Top {len(results)} results:")
for r in results:
    print(f"  [{r.score:.4f}] {r.payload['text']}")


# Query: "What are the benefits of working out?"

# Top 3 results:
#   [0.5814] Regular physical activity strengthens muscles and boosts energy.
#   [0.5408] Exercise improves cardiovascular health and mood.
#   [0.2878] A balanced diet is essential for good health.