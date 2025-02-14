package okta

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/okta/okta-sdk-golang/okta"
)

func TestAccappSamlApplicationImport(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSaml)
	config := mgr.GetFixtures("import.tf", ri, t)
	resourceName := buildResourceFQN(appSaml, ri)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("failed to import into resource into state")
					}
					if s[0].Attributes["preconfigured_app"] != "pagerduty" {
						return errors.New("failed to set required properties when import existing infrastructure")
					}
					return nil
				},
			},
		},
	})
}

// Ensure conditional require logic causes this plan to fail
func TestAccOktaAppSamlApplicationConditionalRequire(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigMissingFields(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSaml, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("missing conditionally required fields, reason: Custom SAML applications must contain these fields*"),
			},
		},
	})
}

// Ensure conditional require logic causes this plan to fail
func TestAccOktaAppSamlApplicationInvalidUrl(t *testing.T) {
	ri := acctest.RandInt()
	config := buildTestSamlConfigInvalidUrl(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSaml, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("config is invalid: failed to validate url, \"123\""),
			},
		},
	})
}

// Test creation of a custom SAML app.
func TestAccOktaAppSamlApplication(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSaml)
	config := mgr.GetFixtures("custom_saml_app.tf", ri, t)
	updatedConfig := mgr.GetFixtures("custom_saml_app_updated.tf", ri, t)
	resourceName := buildResourceFQN(appSaml, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSaml, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttrSet(resourceName, "http_post_binding"),
					resource.TestCheckResourceAttrSet(resourceName, "http_redirect_binding"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata"),
					resource.TestCheckResourceAttrSet(resourceName, "entity_key"),
					resource.TestCheckResourceAttrSet(resourceName, "entity_url"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
				),
			},
		},
	})
}

func TestAccOktaAppSamlApplicationAllFields(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSaml)
	config := mgr.GetFixtures("custom_saml_app.tf", ri, t)
	allFields := mgr.GetFixtures("custom_saml_app_all_fields.tf", ri, t)
	resourceName := buildResourceFQN(appSaml, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSaml, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: allFields,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_template", "${source.login}"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "response_signed", "true"),
					resource.TestCheckResourceAttr(resourceName, "assertion_signed", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "RSA_SHA1"),
					resource.TestCheckResourceAttr(resourceName, "digest_algorithm", "SHA1"),
					resource.TestCheckResourceAttr(resourceName, "honor_force_authn", "true"),
					resource.TestCheckResourceAttr(resourceName, "authn_context_class_ref", "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.name", "Attr One"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.namespace", "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.0.values.0", "val"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.name", "Attr Two"),
					resource.TestCheckResourceAttr(resourceName, "attribute_statements.1.namespace", "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified"),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "sso_url", "http://google.com"),
					resource.TestCheckResourceAttr(resourceName, "recipient", "http://here.com"),
					resource.TestCheckResourceAttr(resourceName, "destination", "http://its-about-the-journey.com"),
					resource.TestCheckResourceAttr(resourceName, "audience", "http://audience.com"),
					resource.TestCheckResourceAttr(resourceName, "subject_name_id_format", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"),
				),
			},
		},
	})
}

// Add and remove groups/users
func TestAccOktaAppSamlApplicationUserGroups(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSaml)
	config := mgr.GetFixtures("saml_app_with_groups_and_users.tf", ri, t)
	updatedConfig := mgr.GetFixtures("saml_app_with_groups_and_users_updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSaml)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(appSaml, createDoesAppExist(okta.NewSamlApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "key_years_valid", "3"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSamlApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttr(resourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
				),
			},
		},
	})
}

func buildTestSamlConfigMissingFields(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         		= "%s"
  status 	    	    = "INACTIVE"
}
`, appSaml, name, name)
}

func buildTestSamlConfigInvalidUrl(rInt int) string {
	name := buildResourceName(rInt)

	return fmt.Sprintf(`
resource "%s" "%s" {
  label         		= "%s"
  status 	    	    = "INACTIVE"
  sso_url      			= "123"
}
`, appSaml, name, name)
}
