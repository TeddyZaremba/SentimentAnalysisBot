package main

import (
	"github.com/bwmarrin/discordgo"

	"github.com/cdipaolo/sentiment"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"time"
	"encoding/json"	
	"fmt"
	"log"
	"io/ioutil"
)

var (
	Token string
	URI string

	config *configStruct
)

type configStruct struct {
	Token string `json: "Token"`
	URI string `json: "URI"`
}

type UserMessage struct {
	time string
	Username string 
	Message string
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token

	return nil
}


var BotId string


func Start() {
	//
	// Create the discord bot
	//
	err := ReadConfig()

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return 
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	BotId = u.ID 
	//
	//


	//
	// Connect to MongoDB
	//
	uri := config.URI
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' config variable")
	}

	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	userColl := client.Database("DiscordBotDB").Collection("DiscordBotCollection")


	//
	//

	//
	// Creates the sentiment analysis model
	//
	model, _ := sentiment.Restore()
	var analysis *sentiment.Analysis
	//
	//

	//
	// Event handler for incoming chats
	//
	messageHandler := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == BotId { // if the bot sent the message then ignore
			return
		}
	
		// Gives a summary of the chat's sentiment
		if m.Content == "!Sentiment"{ 
			cursor, err := userColl.Find(context.TODO(), bson.D{})
			if err != nil {
				panic(err)
			}
			var results []UserMessage
			if err = cursor.All(context.TODO(), &results); err != nil {
				panic(err)
			}
			userMessages := make(map[string]string)
			for _, result := range results {
				fmt.Printf("%+v %+v\n", result.Username, result.Message)

				val, ok := userMessages[result.Username]
				if ok {
					userMessages[result.Username] = val + " " + result.Message
				} else {
					userMessages[result.Username] = result.Message
				}				
			}
			sentimentBotMessage := ""
			for user, element := range userMessages {
				analysis = model.SentimentAnalysis(element, sentiment.English)  
				if analysis.Score == 1 {
					sentimentBotMessage += user + " is happy\n"
				} else {   
					sentimentBotMessage += user + " is unhappy\n"
				}  
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, sentimentBotMessage)
		} else {
			doc := bson.D{{"time",primitive.NewDateTimeFromTime(time.Now().UTC())},{"Username", m.Author.Username},{"Message", m.Content}}

			_, err = userColl.InsertOne(context.TODO(), doc)
			if err != nil {
				fmt.Println(err.Error())
				return 
			}
		}

	}
	// 
	//
	

	//
	// Adds the event handler to the bot and runs the bot
	// 
	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return 
	}

	fmt.Println("Bot is running fine!")
	//
	//
}

func main() {

	Start()

	<-make(chan struct{})
	return
}