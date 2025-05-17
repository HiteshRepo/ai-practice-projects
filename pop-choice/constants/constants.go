package constants

import "pop-choice/models"

var Movies = []models.Movie{
	{
		Title:       "Avatar: The Way of the Water",
		ReleaseYear: "2022",
		Content:     "Avatar: The Way of Water (3 hr 10 min): Jake Sully lives with his newfound family formed on the extrasolar moon Pandora. Once a familiar threat returns to finish what was previously started, Jake must work with Neytiri and the army of the Na'vi race to protect their home. Action, Adventure, Fantasy film released in 2022. Directed by James Cameron Written by James Cameron, Rick Jaffa and Amanda Silver. Starring Sam Worthington, Zoe Saldana and Sigourney Weaver. Rated 7.6 on IMDB",
	},
	{
		Title:       "The Fabelmans",
		ReleaseYear: "2022",
		Content:     "The Fabelmans (2 hr 31 min): Growing up in post-World War II era Arizona, young Sammy Fabelman aspires to become a filmmaker as he reaches adolescence, but soon discovers a shattering family secret and explores how the power of films can help him see the truth. Drama film released in 2022. Directed by Steven Spielberg. Written by Steven Spielberg and Tony Kushner. Starring Michelle Williams, Gabriel LaBelle & Paul Dano. Rated 7.5 on IMDB",
	},
	{
		Title:       "Troll",
		ReleaseYear: "2022",
		Content:     "Troll (1 hr 41 min): Deep in the Dovre mountain, something gigantic wakes up after a thousand years in captivity. The creature destroys everything in its path and quickly approaches Oslo. Norwegian action, adventure, drama film released in 2022. Directed by Roar Uthaug. Written by Espen Aukan and Roar Uthaug. Starring Ine Marie Wilmann, Kim Falck and Mads Sjøgård Pettersen. Rated 5.8 on IMDB",
	},
	{
		Title:       "Everything Everywhere All at Once",
		ReleaseYear: "2022",
		Content:     "Everything Everywhere All at Once (2 hr 19 min): A middle-aged Chinese immigrant is swept up into an insane adventure in which she alone can save existence by exploring other universes and connecting with the lives she could have led. Action, Adventure, Comedy film released in 2022. Directed by Daniel Kwan and Daniel Scheinert. Written by Daniel Kwan and Daniel Scheinert. Starring: Michelle Yeoh, Stephanie Hsu and Jamie Lee Curtis. Rated 7.8 on IMDB",
	},
	{
		Title:       "Oppenheimer",
		ReleaseYear: "2023",
		Content:     "Oppenheimer (3 hr): The story of American scientist, J. Robert Oppenheimer, and his role in the development of the atomic bomb. Biography, Drama, History film released in 2023. Directed by Christopher Nolan. Written by Christopher Nolan, Kai Bird and Martin Sherwin. Starring Cillian Murphy, Emily Blunt and Matt Damon. Rated 8.5 on IMDB",
	},
	{
		Title:       "Barbie",
		ReleaseYear: "2023",
		Content:     "Barbie (1 hr 54 min): Barbie suffers a crisis that leads her to question her world and her existence. Adventure, Comedy, Fantasy film released in 2023. Directed by Greta Gerwig. Written by Greta Gerwig and Noah Baumbach. Starring Margot Robbie, Ryan Gosling and Issa Rae. Rated 7.0 on IMDB",
	},
	{
		Title:       "Spider-Man: Across the Spider-Verse",
		ReleaseYear: "2023",
		Content:     "Spider-Man: Across the Spider-Verse (2 hr 20 min): Miles Morales catapults across the Multiverse, where he encounters a team of Spider-People charged with protecting its very existence. When the heroes clash on how to handle a new threat, Miles must redefine what it means to be a hero. Animation, Action, Adventure film released in 2023. Directed by Joaquim Dos Santos, Kemp Powers an Justin K. Thompson. Written by Phil Lord, Christopher Miller and Dave Callaham. Starring: Shameik Moore, Hailee Steinfeld and Brian Tyree Henry. Rated 8.7 on IMDB",
	},
	{
		Title:       "Pathaan",
		ReleaseYear: "2023",
		Content:     "Pathaan (2 hr 26 min): An Indian agent races against a doomsday clock as a ruthless mercenary, with a bitter vendetta, mounts an apocalyptic attack against the country. Bollywood action, adventure, triller film released in 2023. Directed by Siddharth Anand. Written by Shridhar Raghavan, Abbas Tyrewala and Siddharth Anand. Starring Shah Rukh Khan, Deepika Padukone and John Abraham. Rated 5.9 on IMDB",
	},
	{
		Title:       "RRR",
		ReleaseYear: "2022",
		Content:     "RRR (3 hr 7 min): A fictitious story about two legendary revolutionaries and their journey away from home before they started fighting for their country in the 1920s. South Indian action, drama film released in 2022. Directed by S. S. Rajamouli. Written by Vijayendra Prasad, S. S. Rajamouli and Sai Madhav Burra. Starring N. T. Rama Rao Jr., Ram Charan and Ajay Devgn. Rated 7.8 on IMDB",
	},
}

var InitialListOfQuestions = []string{
	"What is your favorite movie and why?",
	"Are you in a mood for somethig new or a classic?",
	"Do you wanna have fun or do you want something serious?",
}

var ExtraListOfQuestionsForMultiUser = []string{
	"How many folks are watching?",
	"What is the preferred duration?",
}

const (
	PopChoiceSystemMessage = `You are an enthusiastic movie expert who loves recommending movies to people. Some context about movies, and a list of user interests to gauge user's choice. Your main job is to formulate a short answer to the question using the provided context. If you are unsure and cannot find the answer, say, "Unable to recommend movies based on choices." Please do not make up the answer.`
	MultiUserFilterTmpl    = `Recommend movie to be watched by a group of %s people who would like to watch for a duration of %s. You are an enthusiastic movie expert who loves recommending movies to people. You will be given some context about a movie, and a list of user interests to gauge user's choice for each user. Your main job is to formulate a short answer to the question using the provided context. If you are unsure and cannot find the answer, say, "Unable to recommend movies based on choices." Please do not make up the answer.`
	Temperature            = 1.1
	PresencePenalty        = 0.0
	FrequencyPenalty       = 0.0
)

// Supabase
const (
	PopChoiceTblName              = "pop_choice"
	PopChoiceFunctionName         = "match_pop_choice"
	PopChoiceTblContentColumnName = "content"
)
