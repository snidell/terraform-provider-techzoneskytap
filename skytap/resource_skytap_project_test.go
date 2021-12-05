package skytap

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/snidell/terraform-provider-techzoneskytap/skytap/skytap/utils"
)

func init() {
	resource.AddTestSweepers("skytap_project", &resource.Sweeper{
		Name: "skytap_project",
		F:    testSweepSkytapProject,
	})
}

func testSweepSkytapProject(region string) error {
	meta, err := sharedClientForRegion(region)
	if err != nil {
		return err
	}

	client := meta.projectsClient
	ctx := context.TODO()

	log.Printf("[INFO] Retrieving list of projects")
	projects, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving list of projects: %v", err)
	}

	for _, p := range projects.Value {
		if shouldSweepAcceptanceTestResource(*p.Name) {
			log.Printf("destroying project %s", *p.Name)
			if err := client.Delete(ctx, *p.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccSkytapProject_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSkytapProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSkytapProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSkytapProjectExists("skytap_project.foo"),
					resource.TestCheckResourceAttr("skytap_project.foo", "name", fmt.Sprintf("tftest-project-%d", rInt)),
					resource.TestCheckResourceAttrSet("skytap_project.foo", "summary"),
					resource.TestCheckResourceAttr("skytap_project.foo", "auto_add_role_name", ""),
					resource.TestCheckResourceAttr("skytap_project.foo", "show_project_members", "true"),
				),
			},
		},
	})
}

func TestAccSkytapProject_AddEnvironment(t *testing.T) {
	templateID, _, _ := setupEnvironment()
	uniqueSuffix := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSkytapProjectDestroy,
		Steps: []resource.TestStep{
			{
				// Create project with environment
				Config: testAccSkytapProjectConfig_withEnvironment(templateID, uniqueSuffix, "skytap_environment.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSkytapProjectExists("skytap_project.foo"),
					resource.TestCheckTypeSetElemAttrPair("skytap_project.foo", "environment_ids.*", "skytap_environment.foo", "id"),
				),
			},
			{
				// Remove it
				Config: testAccSkytapProjectConfig_withEnvironment(templateID, uniqueSuffix, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSkytapProjectExists("skytap_project.foo"),
					resource.TestCheckResourceAttr("skytap_project.foo", "environment_ids.%", "0"),
				),
			},
			{
				// Add it back again
				Config: testAccSkytapProjectConfig_withEnvironment(templateID, uniqueSuffix, "skytap_environment.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSkytapProjectExists("skytap_project.foo"),
					resource.TestCheckTypeSetElemAttrPair("skytap_project.foo", "environment_ids.*", "skytap_environment.foo", "id"),
				),
			},
		},
	})
}

// Verifies the Project exists
func testAccCheckSkytapProjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %q", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Project ID is set")
		}

		// retrieve the connection established in Provider configuration
		client := testAccProvider.Meta().(*SkytapClient).projectsClient
		ctx := context.TODO()

		// Retrieve our project by referencing its state ID for API lookup

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("project (%s) is not an integer: %v", rs.Primary.ID, err)
		}

		_, err = client.Get(ctx, id)
		if err != nil {
			if utils.ResponseErrorIsNotFound(err) {
				return fmt.Errorf("project (%d) was not found - does not exist", id)
			}

			return fmt.Errorf("error retrieving project (%d): %v", id, err)
		}

		return nil
	}
}

// Verifies the Project has been destroyed
func testAccCheckSkytapProjectDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	client := testAccProvider.Meta().(*SkytapClient).projectsClient
	ctx := context.TODO()

	// loop through the resources in state, verifying each project
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "skytap_project" {
			continue
		}

		// Retrieve our project by referencing it's state ID for API lookup
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("project (%s) is not an integer: %v", rs.Primary.ID, err)
		}

		_, err = client.Get(ctx, id)
		if err != nil {
			if utils.ResponseErrorIsNotFound(err) {
				return nil
			}

			return fmt.Errorf("error waiting for project (%d) to be destroyed: %s", id, err)
		}

		return fmt.Errorf("project still exists: %d", id)
	}

	return nil
}

func testAccSkytapProjectConfig_basic(rInt int) string {
	return fmt.Sprintf(`
      resource "skytap_project" "foo" {
	    name = "tftest-project-%d"
	    summary = "This is a project created by the skytap terraform provider acceptance test"
      }`, rInt)
}

func testAccSkytapProjectConfig_withEnvironment(envTemplateID string, uniqueSuffix int, projectEnvIDs string) string {
	return fmt.Sprintf(`
	resource "skytap_environment" "foo" {
 		template_id = "%s"
 		name 		= "%s-environment-%d"
 		description = "This is an environment to support a skytap project terraform provider acceptance test"
 	}
	resource "skytap_project" "foo" {
	    name = "tftest-project-%d"
	    summary = "This is a project created by the skytap terraform provider acceptance test"
		environment_ids = [%s]
      }
	`, envTemplateID, vmEnvironmentPrefix, uniqueSuffix, uniqueSuffix, projectEnvIDs)
}
