package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

const (
	dateFMT = "2006-01-02T15:04:05.000Z"
)

var (
	cli = http.Client{Timeout: time.Second * 5}
)

type NewCollection struct {
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	PublishDate     string        `json:"publishDate"`
	Teams           []interface{} `json:"teams"`
	CollectionOwner string        `json:"collectionOwner"`
	ReleaseURI      string        `json:"releaseUri,omitempty"`
}

func GetCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "col",
		Long:  "Do some collections stuff",
	}

	cmd.AddCommand(create())
	return cmd
}

func create() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a new scheduled collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				return err
			}

			auth, err := cmd.Flags().GetString("auth")
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			delay, err := cmd.Flags().GetInt("delay")
			if err != nil {
				return err
			}

			t := time.Now().Add(time.Minute * time.Duration(delay))
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)

			fmt.Println(t)

			publishDate := t.Format(dateFMT)
			fmt.Println(publishDate)

			col := &NewCollection{
				Name:            name,
				Type:            "scheduled",
				PublishDate:     publishDate,
				Teams:           make([]interface{}, 0),
				CollectionOwner: "ADMIN",
				ReleaseURI:      "",
			}

			zebedeeURL := fmt.Sprintf("%s/zebedee/collection", host)
			r, err := newAuthenticatedRequest(auth, http.MethodPost, zebedeeURL, col)
			if err != nil {
				return err
			}

			resp, err := cli.Do(r)
			if err != nil {
				return err
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("incorrect http status returned: expected 200, actual: %d", resp.StatusCode)
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			var body interface{}
			if err := json.Unmarshal(b, &body); err != nil {
				return err
			}

			str, err := json.MarshalIndent(body, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(str))
			return nil
		},
	}

	cmd.Flags().StringP("url", "url", "http://localhost:8081", "the zebedee API url (default http://localhost:8081)")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringP("auth", "a", "", "user api auth token")
	cmd.MarkFlagRequired("auth")

	cmd.Flags().StringP("name", "n", "", "the collection name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().IntP("delay", "d", 5, "the number of minutes in the future to schedule the collection for. Default is 5")

	return cmd
}

func newAuthenticatedRequest(auth, method, url string, entity interface{}) (*http.Request, error) {
	var body io.Reader = nil

	if entity != nil {
		b, err := json.Marshal(entity)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(b)
	}

	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("X-Florence-Token", auth)

	return r, nil
}
