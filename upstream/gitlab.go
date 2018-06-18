package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-errors/errors"
)

// Self-hosted Gitlab instances use different domain names
type gitLab struct {
	domain     string
	owner      string
	repository string
}

func (g gitLab) String() string {
	return g.domain + "/" + g.owner + "/" + g.repository
}

func (g gitLab) encoded() string {
	// Gitlab requires the owner + repository to be url-encoded together
	return url.PathEscape(g.owner) + "%2F" + url.PathEscape(g.repository)
}

func (g gitLab) releasesURL() string {
	// API documentation: https://docs.gitlab.com/ee/api/tags.html#list-project-repository-tags
	// Note that the second %s must be url-encoded (see gitLab.encoded())
	return fmt.Sprintf("https://%s/api/v4/projects/%s/repository/tags", g.domain, g.encoded())
}

func (g gitLab) errorWrap(err error) error {
	return errors.WrapPrefix(err, "Failed to obtain Gitlab tag for "+g.String()+" from "+g.releasesURL(), 0)
}

func (g gitLab) errorNotFound() error {
	return errors.Errorf("No Gitlab release found for %s on %s", g, g.releasesURL())
}

// A Gitlab commit descriptor
type gitLabCommit struct {
	// The full commit hash
	Id string `json:"id"`
	// A short commit hash
	ShortId string `json:"short_id"`
	// The pretty commit message
	Title string `json:"title"`
	// When the commit was created
	CreatedAt time.Time `json:"created_at"`
	// An array of commit hashes, but always (??) a single element
	ParentIds []string `json:"parent_ids"`
	// The raw commit message (contains things like '\n')
	Message string `json:"message"`
	// AuthorXXX are from who actually wrote the code, useful if a patch was used
	AuthorName   string    `json:"author_name"`
	AuthorEmail  string    `json:"author_email"`
	AuthoredDate time.Time `json:"authored_date"`
	// CommiterXXX are from who accepted the commit into the branch
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	CommittedDate  time.Time `json:"committed_date"`
}

// The api defines this as "null" unless release notes are added to the tag
// With release notes added, these fields then exist
type gitLabRelease struct {
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

// Describes the individual tags in the returned taglist from the json call
type gitLabTag struct {
	// The actual tag's "version", whatever it was tagged as
	Name string `json:"name"`
	// The message added to the git tag
	Message string `json:"message"`
	// The commit hash
	Target string `json:"target"`
	// A field describing the commit the tag was tagged on
	Commit gitLabCommit `json:"commit"`
	// Note that this isn't like a Github release, and it can be "null"
	Release gitLabRelease `json:"release"`
}

type gitLabMessage struct {
	Message string `json:"message"`
}

func (g gitLab) latestVersion() (Version, error) {
	req, err := http.NewRequest("GET", g.releasesURL(), nil)

	// Obtain Gitlab token for higher request limits, see https://docs.gitlab.com/ee/api/#oauth2-tokens
	token := os.Getenv("GITLAB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if err != nil {
		return "", g.errorWrap(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", g.errorWrap(err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode == http.StatusForbidden {
		var message gitLabMessage
		err = dec.Decode(&message)
		if err == nil && message.Message != "" {
			err = errors.Wrap(message.Message, 0)
		}
		return "", g.errorWrap(err)
	} else if resp.StatusCode == http.StatusNotFound {
		return "", g.errorNotFound()
	}

	// Can't get single tag, has to be an array
	// NOTE: If Gitlab ever adds a "get newest tag" API call then change this
	var taglist []gitLabTag
	err = dec.Decode(&taglist)
	if err != nil {
		return "", g.errorWrap(err)
		// [0] will always be the newest, as its sorted by default
	} else if taglist[0].Name != "" {
		return Version(taglist[0].Name), nil
	}
	return "", g.errorNotFound()
}
