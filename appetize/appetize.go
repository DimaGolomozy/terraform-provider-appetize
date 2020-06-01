package appetize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func appetizeUrl(apiToken string) string {
	return fmt.Sprintf("https://%v@api.appetize.io/v1/apps", apiToken)
}

func appetizeUrlWithPublicKey(apiToken string, publicKey string) string {
	return appetizeUrl(apiToken) + "/" + publicKey
}

type AppOptions struct {
	FilePath              string
	Platform              string
	ButtonText            string
	Note                  string
	PostSessionButtonText string
	Disabled              bool
}

type ListResult struct {
	HasMore string `json:"hasMore"`
	Data    []App  `json:"data"`
}

type App struct {
	PublicKey             string  `json:"publicKey"`
	PrivateKey            string  `json:"privateKey"`
	Platform              string  `json:"platform"`
	Disabled              string  `json:"disabled"`
	Name                  *string `json:"name"`
	ButtonText            *string `json:"buttonText"`
	Note                  *string `json:"note"`
	PostSessionButtonText *string `json:"postSessionButtonText"`
}

type Appetize struct {
	apiToken string
}

func NewAppetize(apiToken string) *Appetize {
	appetize := new(Appetize)
	appetize.apiToken = apiToken
	return appetize
}

func (appetize *Appetize) CreateApp(appOptions *AppOptions) (*App, error) {
	params := createParams(appOptions)
	resp, err := uploadFile(appetizeUrl(appetize.apiToken), appOptions.FilePath, params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var app App
		_ = json.Unmarshal(bodyBytes, &app)
		return &app, nil
	}

	return nil, nil
}

func (appetize *Appetize) UpdateApp(publicKey string, appOptions *AppOptions) (*App, error) {
	params := createParams(appOptions)
	jsonValue, _ := json.Marshal(params)

	resp, err := http.Post(appetizeUrlWithPublicKey(appetize.apiToken, publicKey), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var app App
		_ = json.Unmarshal(bodyBytes, &app)
		return &app, nil
	}

	return nil, nil
}

func (appetize *Appetize) DeleteApp(publicKey string) error {
	req, _ := http.NewRequest("DELETE", appetizeUrlWithPublicKey(appetize.apiToken, publicKey), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code from appetize (%d)", resp.StatusCode)
	}
	return nil
}

func (appetize *Appetize) GetApp(publicKey string) (*App, error) {
	apps, err := appetize.listApps()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for i := 0; i < len(apps); i++ {
		if apps[i].PublicKey == publicKey {
			return &apps[i], nil
		}
	}

	return nil, nil
}

func (appetize *Appetize) listApps() ([]App, error) {
	resp, err := http.Get(appetizeUrl(appetize.apiToken))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var listResult ListResult
		_ = json.Unmarshal(bodyBytes, &listResult)
		return listResult.Data, nil
	}

	return nil, nil
}

func createParams(appOptions *AppOptions) map[string]*string {
	disabled := strconv.FormatBool(appOptions.Disabled)
	params := map[string]*string{
		"platform":              &appOptions.Platform,
		"disabled":              &disabled,
		"buttonText":            &appOptions.ButtonText,
		"postSessionButtonText": &appOptions.PostSessionButtonText,
		"note":                  &appOptions.Note,
	}

	for key, val := range params {
		if *val == "" {
			params[key] = nil
		}
	}

	return params
}
