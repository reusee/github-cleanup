package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"github.com/reusee/e4"
	"github.com/reusee/toml"
	"golang.org/x/oauth2"
)

var (
	ce, he = e4.Check, e4.Handle
	pt     = fmt.Printf
)

func main() {

	// config
	var oauthConfig struct {
		Token string // oauth2 access token, create one at https://github.com/settings/tokens
	}
	dir, err := os.UserConfigDir()
	ce(err)
	content, err := ioutil.ReadFile(filepath.Join(dir, "github.conf.toml"))
	ce(err)
	ce(toml.Unmarshal(content, &oauthConfig))

	// client
	ctx := context.Background()
	client := github.NewClient(
		oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: oauthConfig.Token,
				},
			),
		),
	)

	// archive inactive repos
	page := 0
	for {
		repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		ce(err)

		for _, repo := range repos {
			if time.Since(repo.UpdatedAt.Time) < time.Hour*24*365*1 {
				continue
			}
			//if *repo.Archived {
			//	continue
			//}
			pt(
				"%s\n\t%v\n\t%v\n",
				*repo.Name,
				*repo.UpdatedAt,
				*repo.HTMLURL,
			)
			t := false
			repo.Archived = &t
			_, _, err := client.Repositories.Edit(ctx, *repo.Owner.Login, *repo.Name, repo)
			ce(err)
		}

		if len(repos) == 0 {
			break
		}
		page++
	}

}
