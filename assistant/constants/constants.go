package constants

const (
	MovieAssistantInstructions = "You are great at recommending movies. When asked a question, use the information in the provided file to form a friendly response. If you cannot find the answer in the file, do your best to infer what the answer should be."
	MovieRunInstructions       = `Please do not provide annotations in your reply. Only reply about movies in the provided file. If questions are not related to movies, respond with "Sorry, I don't know." Keep your answers short.`
	MovieAssistantName         = "Movie Expert"
	MovieAssistantPurpose      = "assistants"
	MovieDetailsFilePath       = "./data/movies.txt"
	MovieDetailsFileName       = "movie_details.txt"
	MoviesVectorStoreName      = "movies_vector_store"
)
