package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"github.com/mikkeloscar/aur"
	"golang.org/x/oauth2"
)

func newGitHubClient(accessToken string) *github.Client {
	token := &oauth2.Token{AccessToken: accessToken}
	source := oauth2.StaticTokenSource(token)
	client := oauth2.NewClient(context.TODO(), source)
	return github.NewClient(client)
}

func handlePackage(client *github.Client, pkg *aur.Pkg) error {
	githubURL := regexp.MustCompile("github.com/([^/]+)/([^/]+)")
	match := githubURL.FindSubmatch([]byte(pkg.URL))
	if match == nil {
		return nil
	}
	owner, repo := string(match[1]), string(match[2])
	release, _, err := client.Repositories.GetLatestRelease(context.TODO(), owner, repo)
	if err != nil || release == nil {
		return fmt.Errorf("No release found for %s, %s: %v", pkg.Name, pkg.URL, err)
	}
	fmt.Println(pkg.Name, pkg.URL, pkg.Version, *release.Name, release.CreatedAt)
	return nil
}

func handlePackageForMaintainer(maintainer string, client *github.Client) {
	packages, err := aur.SearchByMaintainer(maintainer)
	if err != nil {
		panic(err)
	}
	for _, pkg := range packages {
		err := handlePackage(client, &pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func main() {
	user := flag.String("user", "", "AUR username")
	token := flag.String("token", "", "GitHub personal access token")
	flag.Parse()
	if *user == "" {
		fmt.Fprintln(os.Stderr, "-user is required")
		flag.Usage()
		os.Exit(1)
	}
	handlePackageForMaintainer(*user, newGitHubClient(*token))
}
