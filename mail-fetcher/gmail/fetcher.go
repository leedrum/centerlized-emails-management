package gmail

import (
	"fmt"
	"os"

	"google.golang.org/api/gmail/v1"
)

type GmailFetcher struct {
	svc *gmail.Service
}

func NewFetcher(svc *gmail.Service) *GmailFetcher {
	return &GmailFetcher{svc: svc}
}

func (g *GmailFetcher) FetchLatestEmails(user string, max int64) ([]string, error) {
	msgs, err := g.svc.Users.Messages.List(user).MaxResults(max).Do()
	if err != nil {
		return nil, err
	}

	var files []string
	for _, m := range msgs.Messages {
		msg, err := g.svc.Users.Messages.Get(user, m.Id).Format("raw").Do()
		if err != nil {
			fmt.Printf("Error fetching message %s: %v\n", m.Id, err)
			continue
		}

		emlFile := fmt.Sprintf("/tmp/%s.eml", m.Id)
		if err := os.WriteFile(emlFile, []byte(msg.Raw), 0644); err != nil {
			return nil, err
		}
		files = append(files, emlFile)
	}

	return files, nil
}
