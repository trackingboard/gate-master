package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/nlopes/slack"

  "time"
  "github.com/kidoman/embd"
  _ "github.com/kidoman/embd/host/rpi"
)

// MentionlessMagicWords contains all the words to open gate without mentioning bot
var MentionlessMagicWords = []string { "go go gadget gate", "oom" }

// MagicWords contains all the words to open gate if mentioning bot
var MagicWords = []string { "open sesame" }

func main() {
  api := slack.New(os.Getenv("SLACK_TOKEN"))

  rtm := api.NewRTM()
  go rtm.ManageConnection()

  userID := ""

  embd.InitGPIO()
  defer embd.CloseGPIO()
  pin, _ := embd.NewDigitalPin(4)
  pin.SetDirection(embd.Out)
  pin.Write(embd.High)

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

      if(messageToBot) {
        if(contains(MagicWords, botMessage)) {
          openGate(pin, rtm, ev)
        }
        if(ev.Text == "restart") {
          return
        }
      } else {
        if(contains(MentionlessMagicWords, botMessage)) {
          openGate(pin, rtm, ev)
        }
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

func openGate(pin embd.DigitalPin, rtm *slack.RTM, ev *slack.MessageEvent) {
  rtm.SendMessage(rtm.NewOutgoingMessage("Opening gate...", ev.Channel))

  pin.Write(embd.Low)

  time.Sleep(1000 * time.Millisecond)

  pin.Write(embd.High)
}

func contains(s []string, e string) bool {
  for _, a := range s {
    if strings.ToLower(a) == strings.ToLower(e) {
      return true
    }
  }
  return false
}