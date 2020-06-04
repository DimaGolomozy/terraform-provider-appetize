package main

import (
	"appetize-provider/appetize"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func NewAppetizer(d *schema.ResourceData) *appetize.Appetize {
	apiToken := d.Get("api_token").(string)
	return appetize.NewAppetize(apiToken)
}

func NewAppOptions(d *schema.ResourceData) *appetize.AppOptions {
	appOptions := &appetize.AppOptions{
		Platform:              d.Get("platform").(string),
		ButtonText:            d.Get("button_text").(string),
		PostSessionButtonText: d.Get("post_session_button_text").(string),
		Note:                  d.Get("note").(string),
		FileType:              d.Get("file_type").(string),
		LaunchUrl:             d.Get("launch_url").(string),
		Timeout:               d.Get("timeout").(int),
		Disabled:              d.Get("disabled").(bool),
		DisabledHome:          d.Get("disable_home").(bool),
		UseLastFrame:          d.Get("use_last_frame").(bool),
	}

	if v, ok := d.GetOk("file_path"); ok {
		appOptions.FilePath = v.(string)
	}

	if v, ok := d.GetOk("url"); ok {
		appOptions.Url = v.(string)
	}

	return appOptions
}
