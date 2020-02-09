package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
)

func getClient(accessToken string) *github.Client {
	// Create authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getLatestMatchingAsset(client *github.Client, config repoConfig) *github.ReleaseAsset {
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), config.Owner, config.Name)
	if err != nil {
		fmt.Printf("Failed to get latest release from %s/%s : %v", config.Owner, config.Name, err)
	} else {
		for _, asset := range release.Assets {
			match, _ := regexp.MatchString(config.AssetPattern, *asset.Name)
			if match {
				return &asset
			}
		}
	}
	return nil
}

func downloadAsset(client *github.Client, config repoConfig, asset *github.ReleaseAsset) {
	// Check if download directory exists, create if not
	if _, err := os.Stat(config.TargetDir); os.IsNotExist(err) {
		os.Mkdir(config.TargetDir, 755)
	}
	assetID := *asset.ID
	assetPath := path.Join(config.TargetDir, *asset.Name)
	// Check if we already have the asset
	if _, err := os.Stat(assetPath); os.IsNotExist(err) {
		fmt.Printf("Starting download of %s\n", *asset.Name)
		httpClient := http.Client{}
		assetReader, _, err := client.Repositories.DownloadReleaseAsset(
			context.Background(), config.Owner, config.Name, assetID, &httpClient)
		defer assetReader.Close()
		if err != nil {
			fmt.Printf("Failed to start asset download: %v\n", err)
		} else {
			out, err := os.Create(assetPath)
			defer out.Close()
			if err != nil {
				fmt.Printf("Failed to create asset file: %v\n", err)
			} else {
				_, err := io.Copy(out, assetReader)
				if err != nil {
					fmt.Printf("Failed to write asset to file: %v\n", err)
				} else {
					fmt.Printf("Successfully downloaded: %s\n", assetPath)
				}
			}
		}
	} else {
		fmt.Printf("Asset already exists at %s\n", assetPath)
	}

}

func GetLatestReleases(conf finalConfig) {
	client := getClient(conf.GithubAccessToken)
	for _, repo := range conf.Repositories {
		asset := getLatestMatchingAsset(client, repo)
		if asset != nil {
			downloadAsset(client, repo, asset)
		}
	}
}
