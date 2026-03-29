from langchain_text_splitters import CharacterTextSplitter, RecursiveCharacterTextSplitter
from sentence_transformers import SentenceTransformer
import numpy as np
import faiss

# -----------------------------
# Sample document
# -----------------------------
text = """
Section: Drug Usage Guidelines

The drug should be taken after meals. It improves absorption and reduces stomach irritation.

Side Effects:
The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

Dosage:
Take one tablet twice daily. Do not exceed the prescribed dosage.
"""

# -----------------------------
# Embedding model
# -----------------------------
model = SentenceTransformer("all-MiniLM-L6-v2")

def embed(texts):
    return np.array(model.encode(texts))

# -----------------------------
# Simple FAISS retrieval
# -----------------------------
def build_index(chunks):
    vectors = embed(chunks)
    index = faiss.IndexFlatL2(vectors.shape[1])
    index.add(vectors)
    return index, vectors

def search(query, index, chunks, k=2):
    q_vec = embed([query])
    _, indices = index.search(q_vec, k)
    return [chunks[i] for i in indices[0]]

print("Query: Can I take this drug with alcohol?\n")
print()
print("Doc: \n", text)
print()

# -----------------------------
# 1. NO CHUNKING
# -----------------------------
print("\n--- NO CHUNKING ---")
no_chunks = [text]
index, _ = build_index(no_chunks)

results = search("Can I take this drug with alcohol?", index, no_chunks)
print("Retrieved:\n", results[0])


# -----------------------------
# 2. FIXED SIZE CHUNKING
# -----------------------------
print("\n--- FIXED SIZE ---")
fixed_splitter = CharacterTextSplitter(
    chunk_size=100,
    chunk_overlap=0
)

fixed_chunks = fixed_splitter.split_text(text)
index, _ = build_index(fixed_chunks)

results = search("Can I take this drug with alcohol?", index, fixed_chunks)
for r in results:
    print("Retrieved:\n", r)
    print()


# -----------------------------
# 3. FIXED SIZE + OVERLAP
# -----------------------------
print("\n--- FIXED SIZE + OVERLAP ---")
overlap_splitter = CharacterTextSplitter(
    chunk_size=100,
    chunk_overlap=30
)

overlap_chunks = overlap_splitter.split_text(text)
index, _ = build_index(overlap_chunks)

results = search("Can I take this drug with alcohol?", index, overlap_chunks)
for r in results:
    print("Retrieved:\n", r)
    print()


# -----------------------------
# 4. RECURSIVE CHUNKING (BEST DEFAULT)
# -----------------------------
print("\n--- RECURSIVE CHUNKING ---")
recursive_splitter = RecursiveCharacterTextSplitter(
    chunk_size=120,
    chunk_overlap=20
)

recursive_chunks = recursive_splitter.split_text(text)
index, _ = build_index(recursive_chunks)

results = search("Can I take this drug with alcohol?", index, recursive_chunks)
for r in results:
    print("Retrieved:\n", r)
    print()


# Output
# Query: Can I take this drug with alcohol?


# Doc: 
 
# Section: Drug Usage Guidelines

# The drug should be taken after meals. It improves absorption and reduces stomach irritation.

# Side Effects:
# The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

# Dosage:
# Take one tablet twice daily. Do not exceed the prescribed dosage.



# --- NO CHUNKING ---
# Retrieved:
 
# Section: Drug Usage Guidelines

# The drug should be taken after meals. It improves absorption and reduces stomach irritation.

# Side Effects:
# The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

# Dosage:
# Take one tablet twice daily. Do not exceed the prescribed dosage.


# --- FIXED SIZE ---
# Created a chunk of size 106, which is longer than the specified 100
# Retrieved:
#  Side Effects:
# The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

# Retrieved:
#  The drug should be taken after meals. It improves absorption and reduces stomach irritation.


# --- FIXED SIZE + OVERLAP ---
# Created a chunk of size 106, which is longer than the specified 100
# Retrieved:
#  Side Effects:
# The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

# Retrieved:
#  The drug should be taken after meals. It improves absorption and reduces stomach irritation.


# --- RECURSIVE CHUNKING ---
# Retrieved:
#  Side Effects:
# The drug should not be taken with alcohol. Combining it with alcohol may cause liver damage.

# Retrieved:
#  The drug should be taken after meals. It improves absorption and reduces stomach irritation.