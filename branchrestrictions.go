package bitbucket

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type BranchRestrictions struct {
	c *Client
}

type BranchRestrictionsBody struct {
	Kind            string `json:"kind"`
	BranchMatchKind string `json:"branch_match_kind"`
	BranchType      string `json:"branch_type"`
	Pattern         string `json:"pattern"`
	Links           struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Value  interface{}                   `json:"value"`
	ID     int                           `json:"id"`
	Users  []branchRestrictionsBodyUser  `json:"users"`
	Groups []branchRestrictionsBodyGroup `json:"groups"`
	Type   string                        `json:"type"`
}

type branchRestrictionsBodyGroup struct {
	Name  string `json:"name"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
		FullSlug string `json:"full_slug"`
		Members  int    `json:"members"`
		Slug     string `json:"slug"`
	} `json:"links"`
}

type branchRestrictionsBodyUser struct {
	Username     string `json:"username"`
	Website      string `json:"website"`
	Display_name string `json:"display_name"`
	UUID         string `json:"uuid"`
	Created_on   string `json:"created_on"`
	Links        struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Repositories struct {
			Href string `json:"href"`
		} `json:"repositories"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
		Followers struct {
			Href string `json:"href"`
		} `json:"followers"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
		Following struct {
			Href string `json:"href"`
		} `json:"following"`
	} `json:"links"`
}

type BranchRestrictionsRes struct {
	Page          int
	Pagelen       int
	MaxDepth      int
	Size          int
	Next          string
	BRestrictions []BranchRestrictionsBody
}

func (b *BranchRestrictions) Gets(bo *BranchRestrictionsOptions) (*BranchRestrictionsRes, error) {

	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.RepoSlug)
	response, err := b.c.executeRaw("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	return decodeBranchRestrictions(bodyString)
}

func (b *BranchRestrictions) Create(bo *BranchRestrictionsOptions) (*BranchRestrictionsBody, error) {
	data := b.BuildBranchRestrictionsBody(bo)
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.RepoSlug)
	response, err := b.c.executeRaw("POST", urlStr, data)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	return decodeBranchRestriction(bodyString)
}

func (b *BranchRestrictions) Get(bo *BranchRestrictionsOptions) (*BranchRestrictionsBody, error) {
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	response, err := b.c.executeRaw("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	return decodeBranchRestriction(bodyString)
}

func (b *BranchRestrictions) Update(bo *BranchRestrictionsOptions) (*BranchRestrictionsBody, error) {
	data := b.BuildBranchRestrictionsBody(bo)
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	response, err := b.c.executeRaw("PUT", urlStr, data)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	return decodeBranchRestriction(bodyString)
}

func (b *BranchRestrictions) Delete(bo *BranchRestrictionsOptions) (interface{}, error) {
	urlStr := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.RepoSlug, bo.ID)
	return b.c.executeRaw("DELETE", urlStr, "")
}

func (b *BranchRestrictions) BuildBranchRestrictionsBody(bo *BranchRestrictionsOptions) string {

	var users []branchRestrictionsBodyUser
	var groups []branchRestrictionsBodyGroup
	for _, u := range bo.Users {
		user := branchRestrictionsBodyUser{
			Username: u,
		}
		users = append(users, user)
	}
	for _, g := range bo.Groups {
		group := branchRestrictionsBodyGroup{
			Name: g,
		}
		groups = append(groups, group)
	}

	body := BranchRestrictionsBody{
		Kind:    bo.Kind,
		Pattern: bo.Pattern,
		Users:   users,
		Groups:  groups,
		Value:   bo.Value,
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func decodeBranchRestrictions(branchRestrictionResStr string) (*BranchRestrictionsRes, error) {
	var branchRestrictResMap map[string]interface{}
	err := json.Unmarshal([]byte(branchRestrictionResStr), &branchRestrictResMap)
	if err != nil {
		return nil, err
	}

	branchRestrictionsArray := branchRestrictResMap["values"].([]interface{})
	var branchRestrictionsSlice []BranchRestrictionsBody
	for _, BranchRestrictionEntry := range branchRestrictionsArray {
		var br BranchRestrictionsBody
		err = mapstructure.Decode(BranchRestrictionEntry, &br)
		if err == nil {
			branchRestrictionsSlice = append(branchRestrictionsSlice, br)
		}
	}

	page, ok := branchRestrictResMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := branchRestrictResMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}

	max_depth, ok := branchRestrictResMap["max_depth"].(float64)
	if !ok {
		max_depth = 0
	}

	size, ok := branchRestrictResMap["size"].(float64)
	if !ok {
		size = 0
	}

	next, ok := branchRestrictResMap["next"].(string)
	if !ok {
		next = ""
	}

	branchRestrictions := BranchRestrictionsRes{
		Page:          int(page),
		Pagelen:       int(pagelen),
		MaxDepth:      int(max_depth),
		Size:          int(size),
		Next:          next,
		BRestrictions: branchRestrictionsSlice,
	}
	return &branchRestrictions, nil
}

func decodeBranchRestriction(branchRestrictionResStr string) (*BranchRestrictionsBody, error) {
	var brResMap map[string]interface{}
	var br BranchRestrictionsBody

	err := json.Unmarshal([]byte(branchRestrictionResStr), &brResMap)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(brResMap, &br)
	if err != nil {
		return nil, err
	}

	return &br, nil
}
