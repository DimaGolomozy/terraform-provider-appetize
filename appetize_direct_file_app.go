package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppetizeDirectFileApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppetizeDirectFileAppCreate,
		Read:   resourceAppetizeAppRead,
		Delete: resourceAppetizeAppDelete,

		Schema: func() map[string]*schema.Schema {
			s := resourceAppetizeApp().Schema

			delete(s, "url")
			s["file_path"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}

			// Everything is ForceNew
			for _, v := range s {
				v.ForceNew = true
			}

			return s
		}(),
	}
}

func resourceAppetizeDirectFileAppCreate(d *schema.ResourceData, m interface{}) error {
	appetizer := NewAppetizer(d)
	app, err := appetizer.CreateApp(NewAppOptions(d))
	if err != nil {
		return err
	}

	d.SetId(app.PublicKey)
	return resourceAppetizeAppRead(d, m)
}
