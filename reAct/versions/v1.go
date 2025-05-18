package versions

import (
	"context"
	"log"

	"github.com/openai/openai-go"
)

func V1(
	ctx context.Context,
	openaiClient openai.Client,
	query string,
) {
	resp, err := openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: "gpt-4",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(query),
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resp.Choices[0].Message.Content)
}

/*
Response:

As an AI, I'm sorry, I do not have real-time access to data such as your location or current weather. However, I can suggest a few general ideas which you can perhaps consider based on your current conditions.

If the weather is sunny/warm:
1. Visit a local park for a picnic
2. Go biking or for a hike
3. Explore outdoor attractions in your city
4. Enjoy some beach activities if you are near a coast
5. Try outdoor water sports such as surfing, sailing or swimming

If the weather is cool:
1. Have a cozy indoor gathering with friends
2. Visit a museum or art exhibition
3. Go to a movie theater or watch a play
4. Try indoor games like bowling, trampolining or escape rooms
5. Do some indoor ice-skating or visit an indoor ski slope

If the weather is rainy:
1. Visit an art gallery or a museum
2. Have a cozy movie marathon at home
3. Read a book at a local coffee shop
4. Try indoor rock climbing
5. Maybe it's time for some shopping at the mall.

Remember to check COVID-19 safety measures in your area before planning your activities.
*/
