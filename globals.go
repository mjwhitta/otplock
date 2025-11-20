package otplock

import _ "embed"

// Version is the package version
const Version string = "1.2.15"

var (
	//go:embed tmpls/advanced_dashboard.html
	advancedDashboard string

	//go:embed tmpls/advanced_new.html
	advancedNew string

	//go:embed tmpls/advanced_response.html
	advancedResp string

	//go:embed tmpls/error_page.html
	errPg string

	//go:embed tmpls/not_found.html
	notFound string

	//go:embed tmpls/simple_dashboard.html
	simpleDashboard string

	//go:embed tmpls/simple_response.html
	simpleResp string
)
