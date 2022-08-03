package appetize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
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
	Url                   string
	FilePath              string
	Platform              string
	ButtonText            string
	Note                  string
	PostSessionButtonText string
	FileType              string
	LaunchUrl             string
	//Timeout               int
	UseLastFrame bool
	DisabledHome bool
	Disabled     bool
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
	FileType              *string `json:"fileType"`
	UseLastFrame          *string `json:"useLastFrame"`
	DisabledHome          *string `json:"disableHome"`
	LaunchUrl             *string `json:"launchUrl"`
	//Timeout               *string `json:"timeout"`
}

type Appetize struct {
	apiToken string
}

func NewAppetizer(apiToken string) *Appetize {
	appetize := new(Appetize)
	appetize.apiToken = apiToken
	return appetize
}

func (appetize *Appetize) CreateApp(appOptions *AppOptions) (*App, error) {
	resp, err := appetize.create(appetizeUrl(appetize.apiToken), appOptions)
	if err != nil {
		return nil, err
	}
	return appetize.read(resp)
}

func (appetize *Appetize) UpdateApp(publicKey string, appOptions *AppOptions) (*App, error) {
	resp, err := appetize.create(appetizeUrlWithPublicKey(appetize.apiToken, publicKey), appOptions)
	if err != nil {
		return nil, err
	}
	return appetize.read(resp)
}

func (appetize *Appetize) DeleteApp(publicKey string) error {
	req, _ := http.NewRequest("DELETE", appetizeUrlWithPublicKey(appetize.apiToken, publicKey), nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code from appetize (%d)", resp.StatusCode)
	}

	return nil
}

func (appetize *Appetize) GetApp(publicKey string) (*App, error) {
	apps, err := appetize.listApps()
	if err != nil {
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
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code from appetize (%d)", resp.StatusCode)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResult ListResult
	_ = json.Unmarshal(bodyBytes, &listResult)
	return listResult.Data, nil
}

func createParams(appOptions *AppOptions) map[string]*string {
	disabled := strconv.FormatBool(appOptions.Disabled)
	disabledHome := strconv.FormatBool(appOptions.DisabledHome)
	useLastFrame := strconv.FormatBool(appOptions.UseLastFrame)
	//timeout := strconv.Itoa(appOptions.Timeout)
	params := map[string]*string{
		"platform":              &appOptions.Platform,
		"disabled":              &disabled,
		"buttonText":            &appOptions.ButtonText,
		"postSessionButtonText": &appOptions.PostSessionButtonText,
		"note":                  &appOptions.Note,
		"fileType":              &appOptions.FileType,
		//"timeout":               &timeout,
		"disableHome":  &disabledHome,
		"useLastFrame": &useLastFrame,
		"launchUrl":    &appOptions.LaunchUrl,
	}

	for key, val := range params {
		if *val == "" {
			params[key] = nil
		}
	}

	return params
}

func (appetize *Appetize) create(url string, appOptions *AppOptions) (*http.Response, error) {
	params := createParams(appOptions)

	if appOptions.Url != "" {
		params["url"] = &appOptions.Url
		jsonValue, _ := json.Marshal(params)

		return http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	} else {
		return uploadFile(url, appOptions.FilePath, params)
	}
}

func (appetize *Appetize) read(resp *http.Response) (*App, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code from appetize (%d)", resp.StatusCode)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var app App
	_ = json.Unmarshal(bodyBytes, &app)
	return &app, nil
}

func NewAppOptions(d *schema.ResourceData) *AppOptions {
	appOptions := &AppOptions{
		Platform:              d.Get("platform").(string),
		ButtonText:            d.Get("button_text").(string),
		PostSessionButtonText: d.Get("post_session_button_text").(string),
		Note:                  d.Get("note").(string),
		FileType:              d.Get("file_type").(string),
		LaunchUrl:             d.Get("launch_url").(string),
		//Timeout:               d.Get("timeout").(int),
		Disabled:     d.Get("disabled").(bool),
		DisabledHome: d.Get("disable_home").(bool),
		UseLastFrame: d.Get("use_last_frame").(bool),
	}

	if v, ok := d.GetOk("file_path"); ok {
		appOptions.FilePath = v.(string)
	}

	if v, ok := d.GetOk("url"); ok {
		appOptions.Url = v.(string)
	}

	return appOptions
}
