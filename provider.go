package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
        Schema: map[string]*schema.Schema{
			"shell": &schema.Schema{
				Type:     schema.TypeString,
                Optional: true,
                Default: "sh",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cli": resourceServer(),
		},
        ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
    config := Config{
	    Shell: d.Get("shell").(string),
	}

	return &config, nil
}
