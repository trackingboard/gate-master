package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/nlopes/slack"
  "github.com/stianeikeland/go-rpio"
)

func main() {
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  rtm := api.NewRTM()
  go rtm.ManageConnection()

  userID := ""

  _ = rpio.Open()

  pin := rpio.Pin(8)

  for msg := range rtm.IncomingEvents {
    switch ev := msg.Data.(type) {

    case *slack.ConnectedEvent:
      fmt.Printf("Logged in as: %s\n", ev.Info.User.ID)
      userID = ev.Info.User.ID

    case *slack.MessageEvent:
      messageToBot := strings.Contains(ev.Text, "<@"+userID+"> ")
      botMessage := strings.Replace(ev.Text, "<@"+userID+"> ", "", 1)

      if(botMessage == "ping" && messageToBot) {
        rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
      }

      if(botMessage == "switch on" && messageToBot) {
        // rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
        pin.High()
      }

      if(botMessage == "switch off" && messageToBot) {
        // rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
        pin.Low()
      }

    case *slack.InvalidAuthEvent:
      fmt.Printf("Invalid credentials")
      return

    default:
      // Ignore other events..
      // fmt.Printf("Unexpected: %v\n", msg.Data)
    }
  }
}