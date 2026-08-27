## 1.0.1 (August 20, 2026)

SECURITY:

* provider: Address CVE-2026-39821 by upgrading Go to 1.25.13 and golang.org/x/net to v0.58.0. CVSS 9.6 Critical.
* provider: Address CVE-2026-56865 by upgrading Go to 1.25.13. CVSS 8.4 High.
* provider: Address CVE-2026-56864 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-33818 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-46600 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-56862 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-56859 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-56860 by upgrading Go to 1.25.13. CVSS 7.5 High.
* provider: Address CVE-2026-56858 by upgrading Go to 1.25.13. CVSS 6.1 Medium.
* provider: Address CVE-2026-56853 by upgrading Go to 1.25.13. CVSS 5.3 Medium.
* provider: Address CVE-2026-25680 by upgrading golang.org/x/net to v0.58.0. CVSS 6.5 Medium.
* provider: Address CVE-2026-42506 by upgrading golang.org/x/net to v0.58.0. CVSS 6.1 Medium.
* provider: Address CVE-2026-42502 by upgrading golang.org/x/net to v0.58.0. CVSS 6.1 Medium.
* provider: Address CVE-2026-25681 by upgrading golang.org/x/net to v0.58.0. CVSS 6.1 Medium.
* provider: Address CVE-2026-27136 by upgrading golang.org/x/net to v0.58.0. CVSS 6.1 Medium.
* provider: Address CVE-2026-46595 by upgrading golang.org/x/crypto to v0.55.0. CVSS 10.0 Critical.
* provider: Address CVE-2026-42508 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39834 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39833 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39832 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39831 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39830 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39829 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-46597 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39828 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39827 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-39835 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-46598 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2025-47914 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2025-58181 by upgrading golang.org/x/crypto to v0.55.0.
* provider: Address CVE-2026-1229 by upgrading github.com/cloudflare/circl to v1.6.5.

## 1.0.0 (November 25, 2025).

FEATURES:

* **New Data Source:** `runtime_registration_token` - Retrieves the current, active registration token
* **New Resource:** `runtime_registration_token` - Manages registration tokens for cluster node registration

NOTES:

* Initial release of the JFrog Runtime provider
* Supports Runtime API v1 endpoints for registration token management
