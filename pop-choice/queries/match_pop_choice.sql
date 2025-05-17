-- Create a function to search for pop_choice
create or replace function match_pop_choice (
  query_embedding vector(1536),
  match_threshold float,
  match_count int
)
returns table (
  id bigint,
  content text,
  similarity float
)
language sql stable
as $$
  select
    pop_choice.id,
    pop_choice.content,
    1 - (pop_choice.embedding <=> query_embedding) as similarity
  from pop_choice
  where 1 - (pop_choice.embedding <=> query_embedding) > match_threshold
  order by similarity desc
  limit match_count;
$$;