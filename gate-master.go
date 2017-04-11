package main

import (
  "fmt"
  // "log"
  "os"

  "github.com/nlopes/slack"

  // "github.com/davecgh/go-spew/spew"
)

func main() {
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  rtm := api.NewRTM()
  go rtm.ManageConnection()

  for msg := range rtm.IncomingEvents {
    switch ev := msg.Data.(type) {

    case *slack.MessageEvent:
      // spew.Dump(ev)
      if(ev.Text == "ping") {
        rtm.SendMessage(rtm.NewOutgoingMessage("pong", ev.Channel))
      }

    case *slack.InvalidAuthEvent:
      fmt.Printf("Invalid credentials")
      return

    default:
      // fmt.Printf("Unexpected: %v\n", msg.Data)
    }
  }
}