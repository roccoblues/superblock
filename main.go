package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/oauth1"
	twoauth "github.com/dghubble/oauth1/twitter"
	"github.com/peterbourgon/ff/v3"
	"github.com/roccoblues/superblock/twitter"
)

const (
	outOfBand = "oob"
)

type application struct {
	config *oauth1.Config
	client *twitter.Client
	logger *log.Logger
}

func main() {
	fs := flag.NewFlagSet("superblock", flag.ExitOnError)
	fs.Usage = func() { usage(fs) }
	var (
		authorizeFlag = fs.Bool("authorize", false, "fetch access token and secret")
		apiKey        = fs.String("api-key", "", "twitter app client key (also via SUPERBLOCK_API_KEY)")
		apiKeySecret  = fs.String("api-key-secret", "", "twitter app client secret (also via SUPERBLOCK_API_KEY_SECRET)")
		token         = fs.String("token", "", "twitter access token (also via SUPERBLOCK_TOKEN)")
		tokenSecret   = fs.String("token-secret", "", "twitter access token secret (also via SUPERBLOCK_TOKEN_SECRET)")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("SUPERBLOCK"))

	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "--api-key is required")
		os.Exit(1)
	}
	if *apiKeySecret == "" {
		fmt.Fprintln(os.Stderr, "--api-key-secret is required")
		os.Exit(1)
	}

	app := application{
		logger: log.New(os.Stdout, "\t", log.Ldate|log.Ltime|log.Lmsgprefix),
		config: &oauth1.Config{
			ConsumerKey:    *apiKey,
			ConsumerSecret: *apiKeySecret,
			CallbackURL:    outOfBand,
			Endpoint:       twoauth.AuthorizeEndpoint,
		},
	}

	if *authorizeFlag {
		err := app.authorize()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(os.Args) != 2 {
		usage(fs)
		os.Exit(1)
	}

	if *token == "" {
		fmt.Fprintf(os.Stderr, "--token is required.\nTo generate one run: superblock --authorize\n")
		os.Exit(1)
	}
	if *tokenSecret == "" {
		fmt.Fprintf(os.Stderr, "--token-secret is required.\nTo generate one run: superblock --authorize\n")
		os.Exit(1)
	}

	accessToken := oauth1.NewToken(*token, *tokenSecret)
	clientOpts := []twitter.ClientOption{
		twitter.EnableRateLimitRetry(),
		twitter.WithLogger(app.logger),
	}
	app.client = twitter.NewClient(app.config.Client(oauth1.NoContext, accessToken), clientOpts...)

	if err := app.block(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage(fs *flag.FlagSet) {
	fmt.Fprintf(fs.Output(), "Usage: superblock [args] <tweet url>\n\n")
	fmt.Fprintln(fs.Output(), "Parameters:")
	fs.PrintDefaults()
}
