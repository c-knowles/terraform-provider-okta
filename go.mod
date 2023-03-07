module github.com/articulate/terraform-provider-okta

go 1.12

require (
	github.com/articulate/oktasdk-go v0.0.0-20190417182045-e41ed7befc56
	github.com/crewjam/saml v0.0.0-20180831135026-ebc5f787b786
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-hclog v0.0.0-20190109152822-4783caec6f2e // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform v0.12.0
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/okta/okta-sdk-golang v0.1.0
	github.com/peterhellberg/link v1.0.0
	github.com/russellhaering/goxmldsig v1.1.1 // indirect
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c // indirect
)

replace github.com/okta/okta-sdk-golang => github.com/articulate/okta-sdk-golang v1.0.4
