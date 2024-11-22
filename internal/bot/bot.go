package bot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"city-hall-lights/internal/model"
	"city-hall-lights/internal/store"
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

	imageMeta, err := store.ReadImageMetadataFromFile("internal/store/images/attribution.json")
	if err != nil {
		panic(err)
	}

	selectedImage := model.ImageMetadata{}
	for i, entry := range imageMeta {
		if entry.FileName == fmt.Sprintf("%s.jpg", event.Color) {
			selectedImage = imageMeta[i]
		}
	}

	blob, err := uploadBlob(client, fmt.Sprintf("internal/store/images/%s", selectedImage.FileName))
	if err != nil {
		panic(err)
	}
	imageEmbed := buildImageEmbed(selectedImage.AltText, blob)
	post := buildPost(event, imageEmbed)
	err = sendPost(client, blueskyHandle, post)
	if err != nil {
		panic(err)
	}

	return
}

func buildImageEmbed(altText string, blob *util.LexBlob) *bsky.FeedPost_Embed {
	return &bsky.FeedPost_Embed{
		EmbedImages: &bsky.EmbedImages{
			LexiconTypeID: "app.bsky.embed.images",
			Images: []*bsky.EmbedImages_Image{
				{
					Alt:   altText,
					Image: blob,
				},
			},
		},
	}
}

func uploadBlob(client *bluesky.Client, path string) (*util.LexBlob, error) {
	buffer, err := store.LoadImageFromFile(path)
	if err != nil {
		return nil, err
	}

	var blob *util.LexBlob
	err = client.CustomCall(func(c *xrpc.Client) error {
		// input := &atproto.RepoUploadBlob{
		// 	Collection: "app.bsky.feed.post",
		// 	Record: &util.LexiconTypeDecoder{
		// 		Val: post,
		// 	},
		// 	Repo: blueskyHandle,
		// }
		output, err := atproto.RepoUploadBlob(context.Background(), c, buffer)
		if err != nil {
			return err
		}
		blob = output.Blob
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func buildPost(event *model.Event, embed *bsky.FeedPost_Embed) *bsky.FeedPost {
	return &bsky.FeedPost{
		Text:      event.Description,
		CreatedAt: time.Now().Local().Format(time.RFC3339),
		Embed:     embed,
	}
}

func sendPost(client *bluesky.Client, blueskyHandle string, post *bsky.FeedPost) error {
	return client.CustomCall(func(c *xrpc.Client) error {
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
}

func createImageEmbed(client *bluesky.Client, blueskyHandle string, post *bsky.FeedPost) error {
	return client.CustomCall(func(c *xrpc.Client) error {
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
}
