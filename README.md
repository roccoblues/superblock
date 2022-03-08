# SuperBlock

Don't like a bad tweet? Block the tweet author, and every single person who liked it. Like [megablock](https://megablock.xyz/) but actually working.

## Usage

```
❯ superblock https://twitter.com/AlenaAlex16/status/1500577207357480968
2022/03/08 16:19:11 	Logged in as Dennis Schön
2022/03/08 16:19:11 	Fetched tweet: Please raise your hand if you voted for Trump🤚

I want everyone to follow you!👍
2022/03/08 16:19:11 	Blocked tweet author Alena Alex 👍🏽🇺🇸 (1329087600787767300) created 1 year ago
2022/03/08 16:19:12 	Blocked user Greg Hawkins (1499886090769084418) created 3 days ago
2022/03/08 16:19:12 	Blocked user Todd Fairl (1442000952429453316) created 5 months ago
2022/03/08 16:19:12 	Blocked user Brya For Colorado (1365175314859573249) created 1 year ago
2022/03/08 16:19:12 	Blocked user gloria vance (1939530716) created 8 years ago
2022/03/08 16:19:13 	Blocked user Theresa (814255960550150145) created 5 years ago
2022/03/08 16:19:13 	Blocked user Dan Schotte (2322837576) created 8 years ago
2022/03/08 16:19:13 	Blocked user Lianna Shanklin (1324027613766189058) created 1 year ago
2022/03/08 16:19:13 	Blocked user Jay (354958028) created 10 years ago
2022/03/08 16:19:13 	Blocked user John B (1485143206769610753) created 1 month ago
2022/03/08 16:19:14 	Blocked user John Rose (2194349587) created 8 years ago
2022/03/08 16:19:14 	Blocked user Lynn emerson (1170048816663453698) created 2 years ago
...
```

## Install

`go install github.com/roccoblues/superblock@latest`

## Setup

### Create Twitter App

The Twitter API requires OAuth for all of its functionality, so you'll need a registered Twitter application. Follow the instructions
[here](https://developer.twitter.com/en/docs/twitter-api/getting-started/getting-access-to-the-twitter-api) to create one.

Once you've created an application, make sure you enable __OAuth 1.0a__ access and set the App permission to __Read and Write__.

The app API key and secret can be specified as command line arguments `--api-key / --api-key-secret` or as environment variables `SUPERBLOCK_API_KEY` and `SUPERBLOCK_API_KEY_SECRET`.

### Generate access token

To generate and access token run `superblock --authorize` and past the PIN shown after you authorized the app. Example:

```
$ superblock --authorize
Paste your PIN here:
Set the following environment variables or pass them with --token and --token-secret

SUPERBLOCK_TOKEN=myToken
SUPERBLOCK_TOKEN_SECRET=myTokenSecret
```

