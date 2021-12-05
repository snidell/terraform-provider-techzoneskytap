package skytap

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/snidell/terraform-provider-techzoneskytap/skytap/skytap/utils"
)

func TestAccDataSourceSkytapTemplate_Basic(t *testing.T) {
	name := utils.GetEnv("SKYTAP_TEMPLATE_NAME", "Ubuntu 18.04.1 LTS Desktop Firstboot")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSkytapTemplateConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.skytap_template.foo", "id"),
				),
			},
		},
	})
}

func testAccDataSourceSkytapTemplateConfig_basic(name string) string {
	return fmt.Sprintf(`
data "skytap_template" "foo" {
	name = "%s"
}

output "id" {
  value = "${data.skytap_template.foo.id}"
}`, name)
}

func TestAccDataSourceSkytapTemplate_RegexMostRecent(t *testing.T) {
	name := utils.GetEnv("SKYTAP_TEMPLATE_NAME", "Ubuntu 18.04.1 LTS Desktop Firstboot")
	namePartial := utils.GetEnv("SKYTAP_TEMPLATE_NAME_PARTIAL", "Ubuntu 18.04.1 LTS")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSkytapTemplateConfig_regexMostRecent(namePartial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.skytap_template.foo", "id"),
					resource.TestCheckResourceAttr("data.skytap_template.foo", "name", name),
				),
			},
		},
	})
}

func testAccDataSourceSkytapTemplateConfig_regexMostRecent(partial string) string {
	return fmt.Sprintf(`
data "skytap_template" "foo" {
	name = "%s"
    most_recent = true
}

output "id" {
  value = "${data.skytap_template.foo.id}"
}`, partial)
}
