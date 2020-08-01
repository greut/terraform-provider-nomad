module github.com/terraform-providers/terraform-provider-nomad

go 1.14

exclude (
	github.com/Sirupsen/logrus v1.1.0
	github.com/Sirupsen/logrus v1.1.1
	github.com/Sirupsen/logrus v1.2.0
	github.com/Sirupsen/logrus v1.3.0
	github.com/Sirupsen/logrus v1.4.0
	github.com/Sirupsen/logrus v1.4.1
)

require (
	github.com/google/go-cmp v0.5.0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/nomad/api v0.0.0-20200731161648-a2a727b02e42
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.0
	github.com/hashicorp/vault v1.4.3
	github.com/mitchellh/mapstructure v1.1.2
	github.com/stretchr/testify v1.5.1
	golang.org/x/tools v0.0.0-20200717024301-6ddee64345a6 // indirect
	google.golang.org/genproto v0.0.0-20200720141249-1244ee217b7e // indirect
)
