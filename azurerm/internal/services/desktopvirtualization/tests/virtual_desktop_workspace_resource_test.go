package tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
)

func TestAccAzureRMDesktopVirtualizationWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDesktopVirtualizationWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_complete(data, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDesktopVirtualizationWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_workspace", "test")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDesktopVirtualizationWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDesktopVirtualizationWorkspace_basic(data, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDesktopVirtualizationWorkspaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDesktopVirtualizationWorkspace_requiresImport(data, clientID, clientSecret),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_desktop_workspace"),
			},
		},
	})
}

func testCheckAzureRMDesktopVirtualizationWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.DesktopVirtualizationWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Desktop Workspace %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Bad: Get virtualDesktopWorspaceClient: %+v", err)
	}
}

func testCheckAzureRMDesktopVirtualizationWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_desktop_workspace" {
			continue
		}

		log.Printf("[WARN] azurerm_virtual_desktop_workspace still exists in state file.")

		id, err := parse.DesktopVirtualizationWorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return fmt.Errorf("Virtual Desktop Workspace still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccAzureRMDesktopVirtualizationWorkspace_basic(data acceptance.TestData, clientID string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_workspace" "test" {
	name                 = "acctws%d"
	location             = azurerm_resource_group.test.location
	resource_group_name  = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDesktopVirtualizationWorkspace_complete(data acceptance.TestData, clientID string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_workspace" "test" {
	name                 = "acctws%d"
	location             = azurerm_resource_group.test.location
	resource_group_name  = azurerm_resource_group.test.name
	friendly_name        = "acceptance test"
	description          = "acceptance test by creating acctws%d"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDesktopVirtualizationWorkspace_requiresImport(data acceptance.TestData, clientID string, clientSecret string) string {
	template := testAccAzureRMDesktopVirtualizationWorkspace_basic(data, clientID, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_workspace" "import" {
	name                 = azurerm_virtual_desktop_workspace.test.name
	location             = azurerm_virtual_desktop_workspace.test.location
	resource_group_name  = azurerm_virtual_desktop_workspace.test.resource_group_name
}
`, template)
}
