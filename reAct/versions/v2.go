package versions

import (
	"context"
	"fmt"
	"log"
	"react/constants"
	"react/models"
	"react/tools"
	"react/utils"

	"github.com/openai/openai-go"
)

/**
 * Goal - build an agent that can answer any questions that might require knowledge about
 * my current location and the current weather at my location.
 */

/**
PLAN:
1. Design a well-written ReAct prompt
2. Build a loop for my agent to run in.
3. Parse any actions that the LLM determines are necessary
	a. Split the string on the newline character \n
	b. Search through the array of strings for one that has "Action:"
		regex to use: /^Action: (\w+): (.*)$/
	c. Parse the action (function and parameter) from the string
	d. Calling the function
	e. Add an "Observation" message with the results of the function call
4. End condition - final Answer is given
*/

func V2(
	ctx context.Context,
	openaiClient openai.Client,
	query string,
) {
	messages := []openai.ChatCompletionMessageParamUnion{}

	systemMessage := openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(constants.WellWrittenReActSystemPrompt),
			},
		},
	}

	userMessage := openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: openai.String(query),
			},
		},
	}

	messages = append(messages, systemMessage)
	messages = append(messages, userMessage)

	for i := 0; i < constants.MaxIterations; i++ {
		log.Printf("Iteration #%d", i+1)

		resp, err := openaiClient.Chat.Completions.New(
			ctx,
			openai.ChatCompletionNewParams{
				Model:    "gpt-4",
				Messages: messages,
			},
		)
		if err != nil {
			log.Fatalln(err)
		}

		content := resp.Choices[0].Message.Content

		log.Println(content)

		assistantMessage := openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(content),
				},
			},
		}

		messages = append(messages, assistantMessage)

		action, err := utils.ActionExtractor(content)
		if err != nil && err == utils.ErrNoActionFound {
			return
		}

		if err != nil {
			log.Fatalln(err)
		}

		actionResponse := invokeValidAction(action)

		assistantMessage = openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(fmt.Sprintf("Observation: %s", actionResponse)),
				},
			},
		}

		messages = append(messages, assistantMessage)
	}
}

func invokeValidAction(action utils.Action) string {
	switch action.FunctionName {
	case "getCurrentWeather":
		log.Println("calling function getCurrentWeather")

		if len(action.Arguments) != 1 {
			log.Fatalln("invalid number of arguments for getCurrentWeather action")
		}

		weather := tools.
			GetNewHardCodedTool().
			GetCurrentWeather(models.Location{Address: action.Arguments[0]})

		return weather.ToString()

	case "getLocation":
		log.Println("calling function getLocation")

		if len(action.Arguments) != 1 {
			log.Fatalln("invalid number of arguments for getLocation action")
		}

		if !utils.IsEmpty(action.Arguments[0]) {
			log.Fatalln("getLocation action requires no args")
		}

		loc := tools.
			GetNewHardCodedTool().
			GetLocation()

		return loc.ToString()
	}

	log.Fatalf("unknown action invoked: %s", action.FunctionName)

	return ""
}

/*
Response:

2025/05/18 22:06:46 Iteration #1
2025/05/18 22:06:48 Thought: Since I do not know the user's exact location, I should check it to get relevant information on activities.
Action: getLocation: "NONE"
2025/05/18 22:06:48 calling function getLocation
2025/05/18 22:06:48 Iteration #2
2025/05/18 22:06:49 Thought: Now that I know the user is in Bhubaneswar, Odisha, we should get the current weather to determine appropriate activities.
Action: getCurrentWeather: "Delta Square, Bhubaneswar, Odisha, India"
2025/05/18 22:06:49 calling function getCurrentWeather
2025/05/18 22:06:49 Iteration #3
2025/05/18 22:06:56 Thought: With the current weather being sunny and warm at 35 degrees Celsius, the activities should be preferably outdoors and ones that allow the user to stay cool. In Bhubaneswar, the user could visit places of interest or indulge in delightful activities fitting the weather.

Answer: Here are a few ideas for activities to do in Bhubaneswar, Odisha given the sunny weather:

1. Visit the Nandankanan Zoological Park. The park has a white tiger safari and a botanical garden. Don't forget your water bottle and a hat to protect yourself from the sun.
2. Explore the Dhauli Giri Hills. Located on the banks of the river Daya, it's a great place for a short trek.
3. Take a stroll in the Ekamra Kanan Botanical Gardens. It's the perfect place to relax under the shades of beautiful trees.
4. Visit Udaygiri & Khandagiri caves. The caves are considered to be one of the earliest groups of Jain rock-cut shelters.
5. Try out local Odisha cuisine at a nice open-air restaurant.

Remember to stay hydrated and wear sunscreen. Enjoy your time at these beautiful places.
*/
