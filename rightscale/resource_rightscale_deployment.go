package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"

	"fmt"
	"sync"

	"gopkg.in/rightscale/rsc.v5/cm15"
	"gopkg.in/rightscale/rsc.v5/rsapi"
)

var mutex = &sync.Mutex{}

func resourceRightScaleDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceRightScaleDeploymentCreate,
		Read:   resourceRightScaleDeploymentRead,
		Delete: resourceRightScaleDeploymentDelete,

		Schema: map[string]*schema.Schema{
			/*			"href": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
			*/
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, //FIX
			},
		},
	}
}

func resourceRightScaleDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	client := meta.(*cm15.API)

	deploymentLocator := client.DeploymentLocator("/api/deployments")
	deployment, err := deploymentLocator.Create(&cm15.DeploymentParam{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})

	if err != nil {
		fmt.Errorf("[RIGHTSCALE] DEPLOYMENT CREATE ERROR: %s", err.Error())
	}

	// Set this resource id to RightScale HREF
	d.SetId(string(deployment.Href))

	mutex.Unlock()
	return resourceRightScaleDeploymentRead(d, meta)
}

func resourceRightScaleDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()

	client := meta.(*cm15.API)
	deployment, err := client.DeploymentLocator(d.Id()).Show(rsapi.APIParams{})

	if err != nil {
		fmt.Printf("[RIGHTSCALE] DEPLOYMENT READ ERROR %s", err.Error())
	}

	fmt.Printf("[RIGHTSCALE] DESCRIPTION: %s", deployment.Description)
	return nil
}

func resourceRightScaleDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
