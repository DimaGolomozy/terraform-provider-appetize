package appetize

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAppetizeApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppetizeAppCreate,
		Read:   resourceAppetizeAppRead,
		Update: resourceAppetizeAppUpdate,
		Delete: resourceAppetizeAppDelete,

		Schema: map[string]*schema.Schema{
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPETIZE_API_TOKEN", nil),
			},
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"platform": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ios", "android"}, false),
			},
			"file_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			//"timeout": &schema.Schema{
			//	Type:     schema.TypeInt,
			//	Optional: true,
			//},
			"disable_home": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_last_frame": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"launch_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"button_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"post_session_button_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppetizeAppCreate(d *schema.ResourceData, m interface{}) error {
	appetizer := NewAppetizer(d)
	app, err := appetizer.CreateApp(NewAppOptions(d))
	if err != nil {
		return err
	}

	d.SetId(app.PublicKey)
	return resourceAppetizeAppRead(d, m)
}

func resourceAppetizeAppRead(d *schema.ResourceData, m interface{}) error {
	appetizer := NewAppetizer(d)
	app, err := appetizer.GetApp(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

	if app == nil {
		fmt.Printf("cant find app with id (%s)", d.Id())
		d.SetId("")
	} else {
		d.Set("public_key", app.PublicKey)
		d.Set("private_key", app.PrivateKey)
		d.Set("platform", app.Platform)
		d.Set("disabled", app.Disabled)

		if app.Name != nil {
			d.Set("name", app.Name)
		}
		if app.ButtonText != nil {
			d.Set("button_text", app.ButtonText)
		}
		if app.PostSessionButtonText != nil {
			d.Set("post_session_button_text", app.PostSessionButtonText)
		}
		if app.Note != nil {
			d.Set("note", app.Note)
		}
		if app.FileType != nil {
			d.Set("file_type", app.FileType)
		}
		if app.UseLastFrame != nil {
			d.Set("use_last_frame", app.UseLastFrame)
		}
		if app.DisabledHome != nil {
			d.Set("disabled_home", app.DisabledHome)
		}
		if app.LaunchUrl != nil {
			d.Set("launch_url", app.LaunchUrl)
		}
		//if app.Timeout != nil {
		//	d.Set("timeout", app.Timeout)
		//}
	}

	return nil
}

func resourceAppetizeAppUpdate(d *schema.ResourceData, m interface{}) error {
	isChanged := d.HasChanges("url", "platform", "disabled", "note", "button_text", "post_session_button_text")
	if isChanged {
		appetizer := NewAppetizer(d)
		_, err := appetizer.UpdateApp(d.Id(), NewAppOptions(d))
		if err != nil {
			return err
		}
	}
	return resourceAppetizeAppRead(d, m)
}

func resourceAppetizeAppDelete(d *schema.ResourceData, m interface{}) error {
	appetizer := NewAppetizer(d)
	err := appetizer.DeleteApp(d.Id())
	if err != nil {
		return err
	}

	return nil
}
