package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Users []User

type User struct {
	Guid                  string `json:"guid"`
	Admin                 bool   `json:"admin"`
	Active                bool   `json:"active"`
	DefaultSpaceGUID      string `json:"default_space_guid"`
	Username              string `json:"username"`
	SpacesURL             string `json:"spaces_url"`
	OrgsURL               string `json:"organizations_url"`
	ManagedOrgsURL        string `json:"managed_organizations_url"`
	BillingManagedOrgsURL string `json:"billing_managed_organizations_url"`
	AuditedOrgsURL        string `json:"audited_organizations_url"`
	ManagedSpacesURL      string `json:"managed_spaces_url"`
	AuditedSpacesURL      string `json:"audited_spaces_url"`
	c                     *Client
}

type UserResource struct {
	Meta   Meta `json:"metadata"`
	Entity User `json:"entity"`
}

type UserResponse struct {
	Count     int            `json:"total_results"`
	Pages     int            `json:"total_pages"`
	NextUrl   string         `json:"next_url"`
	Resources []UserResource `json:"resources"`
}

func (c *Client) ListUsersByQuery(query url.Values) (Users, error) {
	var users []User
	requestUrl := "/v2/users?" + query.Encode()
	for {
		userResp, err := c.getUserResponse(requestUrl)
		if err != nil {
			return []User{}, err
		}
		for _, user := range userResp.Resources {
			user.Entity.Guid = user.Meta.Guid
			user.Entity.c = c
			users = append(users, user.Entity)
		}
		requestUrl = userResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) ListUsers() (Users, error) {
	return c.ListUsersByQuery(nil)
}

func (u Users) GetUserByUsername(username string) User {
	for _, user := range u {
		if user.Username == username {
			return user
		}
	}
	return User{}
}

func (c *Client) getUserResponse(requestUrl string) (UserResponse, error) {
	var userResp UserResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return UserResponse{}, fmt.Errorf("Error requesting users %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return UserResponse{}, fmt.Errorf("Error reading user request %v", err)
	}
	err = json.Unmarshal(resBody, &userResp)
	if err != nil {
		return UserResponse{}, fmt.Errorf("Error unmarshalling user %v", err)
	}
	return userResp, nil
}
