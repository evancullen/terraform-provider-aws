package aws

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"retention_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	conn := meta.(*conns.AWSClient).CloudWatchLogsConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	logGroup, err := lookupCloudWatchLogGroup(conn, name)
	if err != nil {
		return err
	}
	if logGroup == nil {
		return fmt.Errorf("No log group named %s found\n", name)
	}

	d.SetId(name)
	d.Set("arn", logGroup.Arn)
	d.Set("creation_time", logGroup.CreationTime)
	d.Set("retention_in_days", logGroup.RetentionInDays)
	d.Set("kms_key_id", logGroup.KmsKeyId)

	tags, err := tftags.CloudwatchlogsListTags(conn, name)

	if err != nil {
		return fmt.Errorf("error listing tags for CloudWatch Logs Group (%s): %w", name, err)
	}

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
