# Discord Sentiment Analysis Bot

A discord bot that analyzes the sentiment of the channel's chat. It logs all the chats into a MongoDB database and analyzes them to output into the chat.

## Setup

- Go to the developer portal in discord and create an application. Get the token and set this in the config.json file as TOKEN.
- Set your MONGODB_URI in the config.json file

## Installation
go version: 19.4
```bash
go get "github.com/joho/godotenv"
go get "github.com/bwmarrin/discordgo"
go get "github.com/cdipaolo/sentiment"
```

## Execution / Usage

```bash
go run main.go
```
Type "!Sentiment" in the chat to run a full sentiment analysis on all the chats.
