package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/bamiesking/git-cloud/structs"
	"github.com/bamiesking/git-cloud/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := ".git/cloud/gdrive-token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Set up channels and context
	p := make(chan int)
	c := make(chan string)

	// Spawn goroutine containing HTTP listener
	go utils.AwaitOAuthRedirect(utils.OAuthHandlerGenerator, p, c)

	callbackUrl := fmt.Sprintf("http://localhost:%d/", <-p)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", callbackUrl))

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	utils.Open(authURL)

	tok, err := config.Exchange(context.TODO(), <-c, oauth2.SetAuthURLParam("redirect_uri", callbackUrl))
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func FetchGDrive(handle string) structs.CloudFileInfo {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// TODO: reduce this scope to individual files
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.Get(handle).Fields(googleapi.Field("name,size,modifiedTime")).Do()
	utils.Handle(err)

	file, err := srv.Files.Get(handle).Download()
	utils.Handle(err)

	gitPath, err := utils.GitRepoPath()
	utils.Handle(err)

	cache, err := os.Create(path.Join(gitPath, ".git/cloud/cache", handle))

	utils.Handle(err)
	cacheWriter := io.Writer(cache)
	_, err = io.Copy(cacheWriter, file.Body)
	utils.Handle(err)

	defer file.Body.Close()
	info := structs.CloudFileInfo{}
	if r != nil {
		info.Name = r.Name
		info.Size = r.Size
		t, err := time.Parse(time.RFC3339, r.ModifiedTime)
		if err != nil {
			log.Printf("Failed to parse ModifiedTime for file %s (%s)", r.Name, handle)
			return info
		}
		info.DateModified = t
	}
	return info
}
