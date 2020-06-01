package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
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
		//return fmt.Errorf("Error launching source instance: %s", err)
		return err
	}

	if app == nil {
		return fmt.Errorf("Error launching source instance: %s", err)
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
	}

	return nil
}

func resourceAppetizeAppUpdate(d *schema.ResourceData, m interface{}) error {
	//appetizer := NewAppetizer(d)
	//app, err := appetizer.CreateApp(NewAppOptions(d))
	//if err != nil {
	//	//return fmt.Errorf("Error launching source instance: %s", err)
	//	return err
	//}
	//
	//if app == nil {
	//	return fmt.Errorf("Error launching source instance: %s", err)
	//}
	//
	//d.SetId(app.PublicKey)
	return resourceAppetizeAppRead(d, m)
}

func resourceAppetizeAppDelete(d *schema.ResourceData, m interface{}) error {
	appetizer := NewAppetizer(d)
	err := appetizer.DeleteApp(d.Id())
	if err != nil {
		return fmt.Errorf("error deleting app (%s): %s", d.Id(), err)
	}

	return nil
}
