package main

import (
	"errors"
	"regexp"
	"strings"
)

type Subscriptions struct {
	Client *APIClient
}

var sanitizeUsername = regexp.MustCompile("[^a-zA-Z0-9_-]")

func (x *Subscriptions) GetSubscriptions() ([]API_Common_Account, error) {
	readSub := &API_Subscriptions{}
	err := x.Client.GetAPI("https://apiv2.fansly.com/api/v1/subscriptions", readSub)
	if err != nil {
		return nil, err
	}
	if !readSub.Success {
		return nil, errors.New("response is not successful")
	}

	accts := make([]string, 0)
	for _, sub := range readSub.Response.Subscriptions {
		accts = append(accts, sub.AccountId)
	}

	if len(accts) > 0 {
		readAcct := &API_Account{}
		err := x.Client.GetAPI("https://apiv2.fansly.com/api/v1/account?ids="+strings.Join(accts, ","), readAcct)
		if err != nil {
			return nil, err
		}
		if !readAcct.Success {
			return nil, errors.New("response is not successful")
		}

		for i := range readAcct.Response {
			readAcct.Response[i].Username = sanitizeUsername.ReplaceAllString(readAcct.Response[i].Username, "")
		}

		return readAcct.Response, nil
	}

	return nil, nil
}
