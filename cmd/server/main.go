package main

// Imports
import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {

	// Load env 
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	
	// New socketClient
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
		
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				// Handle mentions events
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					err := handleEventMessage(eventsAPIEvent, client)
					if err != nil {
						log.Fatal(err)
					}
				// Handle commands
				case socketmode.EventTypeSlashCommand:
					command, ok := event.Data.(slack.SlashCommand)
					if !ok {
						log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
						continue
					}
					socketClient.Ack(*event.Request)
					err := handleSlashCommand(command, client)
					if err != nil {
						log.Fatal(err)
					}

				}
			}

		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// Botmentionhandler
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	text := strings.ToLower(event.Text)

	attachment := slack.Attachment{}
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") {
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Color = "#4af030"
	} else {
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Color = "#3d3d3d"
	}

	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

// Handles commands 
func handleSlashCommand(command slack.SlashCommand, client *slack.Client) error {
	switch command.Command {
	case "/p??r":
		return handleParCommand(command, client)
/*     case "/squid":
        return handleSquidCommand(command, client) */
    case "/squid":
        return handleCrabCommand(command, client)
	}

	return nil
}

func handleParCommand(command slack.SlashCommand, client *slack.Client) error {
	// Setup message
	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("Hello %sP??RY", command.Text)
	attachment.Color = "#eb6d54"
	attachment.ImageURL = "https://www.indiewire.com/wp-content/uploads/2021/10/squid-game.png?resize=800,472"

	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func handleSquidCommand(command slack.SlashCommand, client *slack.Client) error {
	// Setup message
	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("Hello %sSQUIDDY", command.Text)
	attachment.Color = "#c67ed6"
	attachment.ImageURL = "https://ichef.bbci.co.uk/news/976/cpsprodpb/BAE1/production/_118314874_squid1.jpg"

	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func handleCrabCommand(command slack.SlashCommand, client *slack.Client) error {
	// Setup message
	attachment := slack.Attachment{}
	attachment.Text = fmt.Sprintf("oh %scrab!", command.Text)
	attachment.Color = "#eb9e34"
	attachment.ImageURL = "http://clipart-library.com/image_gallery2/Crab-Free-Download-PNG.png"


	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}