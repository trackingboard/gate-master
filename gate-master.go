package main

import (
  "fmt"
  "os"
  "strings"
  "github.com/nlopes/slack"

  "time"
  "github.com/kidoman/embd"
  _ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

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

      if(botMessage == "open sesame" && messageToBot) {
        rtm.SendMessage(rtm.NewOutgoingMessage("Opening gate...", ev.Channel))

        pin.Write(embd.Low)

        time.Sleep(1000 * time.Millisecond)

        pin.Write(embd.High)

        // embd.InitGPIO()

        // embd.SetDirection(4, embd.Out)
        // embd.DigitalWrite(4, embd.Low)

        // time.Sleep(2000 * time.Millisecond)

        // embd.CloseGPIO()
        // embd.DigitalWrite(4, embd.Low)
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