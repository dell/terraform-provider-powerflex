package powerflex

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSDSResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config: ProviderConfigForTesting + createSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "2"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "156"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rfcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "performance_profile", "Compact"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   "10.247.100.232",
						"role": "all",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   "10.10.10.1",
						"role": "sdcOnly",
					}),
				),
			},
			// update sds name
			// update sds ips from all, sdcOnly to sdsOnly, all
			// increase rmcache
			// disable rfcache
			// Enable high performance profile
			{
				Config: ProviderConfigForTesting + updateSDSTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "2"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "256"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rfcache_enabled", "false"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "performance_profile", "HighPerformance"),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   "10.247.100.232",
						"role": "sdsOnly",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("powerflex_sds.sds", "ip_list.*", map[string]string{
						"ip":   "10.10.10.2",
						"role": "sdcOnly",
					}),
				),
			},
			// disable sds rmcache
			// re-enable rfcache
			// Disable high performance profile
			{
				Config: ProviderConfigForTesting + updateSDSTest2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_02"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "2"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_size_in_mb", "256"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rmcache_enabled", "false"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "rfcache_enabled", "true"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "performance_profile", "Compact"),
				),
			},
		},
	})
}

func TestAccSDSResourceDuplicateIP(t *testing.T) {
	createSDSTestMany := `
		resource "powerflex_sds" "sds" {
			name = "Tf_SDS_01"
			ip_list = [
				{
					ip = "10.247.100.232"
					role = "sdsOnly"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				},
				{
					ip = "10.10.10.2"
					role = "sdcOnly"
				}
			]
			protection_domain_id = "4eeb304600000000"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config: ProviderConfigForTesting + createSDSTestMany,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerflex_sds.sds", "name", "Tf_SDS_01"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "protection_domain_id", "4eeb304600000000"),
					resource.TestCheckResourceAttr("powerflex_sds.sds", "ip_list.#", "3"),
				),
			},
		},
	})
}

func TestAccSDSResourceCreateWithoutIP(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "4eeb304600000000"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*ip_list.*`),
			},
		},
	})
}

func TestAccSDSResourceCreateWithBadRole(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "4eeb304600000000"
			ip_list = [
				{
					ip = "10.10.10.1"
					role = "invalidRole"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*role.*`),
			},
		},
	})
}

func TestAccSDSResourceCreateWithBadPerformanceProfile(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			protection_domain_id = "4eeb304600000000"
			ip_list = [
				{
					ip = "10.247.100.232"
					role = "all"
				}
			]
			performance_profile = "inv"
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*performance_profile.*`),
			},
		},
	})
}

func TestAccSDSResourceCreateWithoutPD(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			name = "Sds123"
			ip_list = [
				{
					ip = "10.247.100.232"
					role = "all"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*protection_domain.*`),
			},
		},
	})
}

func TestAccSDSResourceCreateWithoutName(t *testing.T) {
	createInvalidConfig := `
		resource "powerflex_sds" "invalid" {
			protection_domain_id = "4eeb304600000000"
			ip_list = [
				{
					ip = "10.247.100.232"
					role = "all"
				},
				{
					ip = "10.10.10.1"
					role = "sdcOnly"
				}
			]
		}
		`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create sds test
			{
				Config:      ProviderConfigForTesting + createInvalidConfig,
				ExpectError: regexp.MustCompile(`.*name.*`),
			},
		},
	})
}

var createSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_01"
	ip_list = [
		{
			ip = "10.247.100.232"
			role = "all"
		},
		{
			ip = "10.10.10.1"
			role = "sdcOnly"
		}
	]
	performance_profile = "Compact"
	rmcache_enabled = true
	rmcache_size_in_mb = 156
	# num_of_io_buffers = 4
	rfcache_enabled = true
	drl_mode = "NonVolatile"
	protection_domain_id = "4eeb304600000000"
}
`

var updateSDSTest = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_02"
	ip_list = [
		{
			ip = "10.247.100.232"
			role = "sdsOnly"
		},
		{
			ip = "10.10.10.2"
			role = "sdcOnly"
		}
	]
	drl_mode = "Volatile"
	performance_profile = "HighPerformance"
	rmcache_size_in_mb = 256
	rmcache_enabled = true
	rfcache_enabled = false
	protection_domain_id = "4eeb304600000000"
}
`

var updateSDSTest2 = `
resource "powerflex_sds" "sds" {
	name = "Tf_SDS_02"
	ip_list = [
		{
			ip = "10.247.100.232"
			role = "sdsOnly"
		},
		{
			ip = "10.10.10.2"
			role = "sdcOnly"
		}
	]
	performance_profile = "Compact"
	rmcache_enabled = false
	rfcache_enabled = true
	protection_domain_id = "4eeb304600000000"
}
`
