package repository_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dre2004/go-nexus-client/nexus3/schema/repository"
	"github.com/dre2004/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceRepositoryDockerGroupConfig() string {
	return `
data "nexus_repository_docker_group" "acceptance" {
	name   = nexus_repository_docker_group.acceptance.id
}`
}

func TestAccDataSourceRepositoryDockerGroup(t *testing.T) {
	repoHosted := testAccResourceRepositoryDockerHosted()
	repoGroup := repository.DockerGroupRepository{
		Name:   fmt.Sprintf("acceptance-%s", acctest.RandString(10)),
		Online: true,
		Storage: repository.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		Docker: repository.Docker{
			ForceBasicAuth: true,
			V1Enabled:      true,
		},
		Group: repository.GroupDeploy{
			MemberNames: []string{repoHosted.Name},
		},
	}
	dataSourceName := "data.nexus_repository_docker_group.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				// TODO: Need to figure out why this is happening
				Config: testAccResourceRepositoryDockerHostedConfig(repoHosted) + testAccResourceRepositoryDockerGroupConfig(repoGroup) + testAccDataSourceRepositoryDockerGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "id", repoGroup.Name),
							resource.TestCheckResourceAttr(dataSourceName, "name", repoGroup.Name),
							resource.TestCheckResourceAttr(dataSourceName, "online", strconv.FormatBool(repoGroup.Online)),
						),
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "docker.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "docker.0.force_basic_auth", strconv.FormatBool(repoGroup.Docker.ForceBasicAuth)),
							resource.TestCheckResourceAttr(dataSourceName, "docker.0.v1_enabled", strconv.FormatBool(repoGroup.Docker.V1Enabled)),
						),
						resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "storage.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "storage.0.blob_store_name", repoGroup.Storage.BlobStoreName),
							resource.TestCheckResourceAttr(dataSourceName, "storage.0.strict_content_type_validation", strconv.FormatBool(repoGroup.Storage.StrictContentTypeValidation)),
							resource.TestCheckResourceAttr(dataSourceName, "group.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "group.0.member_names.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "group.0.member_names.0", repoGroup.Group.MemberNames[0]),
						),
					),
				),
			},
		},
	})
}
