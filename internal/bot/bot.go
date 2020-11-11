package bot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"city-hall-lights/internal/model"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
	"github.com/tailscale/go-bluesky"
)

func CreateAndSendPost(event *model.Event) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	blueskyHandle := os.Getenv("BLUESKY_IDENTIFIER")
	blueskyAppkey := os.Getenv("BLUESKY_APP_PASSWORD")

	client, err := bluesky.Dial(ctx, bluesky.ServerBskySocial)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.Login(ctx, blueskyHandle, blueskyAppkey)
	switch {
	case errors.Is(err, bluesky.ErrMasterCredentials):
		panic("You're not allowed to use your full-access credentials, please create an appkey")
	case errors.Is(err, bluesky.ErrLoginUnauthorized):
		panic("Username of application password seems incorrect, please double check")
	case err != nil:
		fmt.Println(err)
		panic("Something else went wrong, please look at the returned error")
	}

	post := &bsky.FeedPost{
		Text:      event.Description,
		CreatedAt: time.Now().Local().Format(time.RFC3339),
	}

	err = client.CustomCall(func(c *xrpc.Client) error {
		input := &atproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Record: &util.LexiconTypeDecoder{
				Val: post,
			},
			Repo: blueskyHandle,
		}
		_, err := atproto.RepoCreateRecord(context.Background(), c, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return
}
