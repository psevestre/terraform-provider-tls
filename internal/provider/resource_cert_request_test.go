package provider

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCertRequest(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
                    resource "tls_cert_request" "test1" {
                        subject {
                            common_name = "example.com"
                            organization = "Example, Inc"
                            organizational_unit = "Department of Terraform Testing"
                            street_address = ["5879 Cotton Link"]
                            locality = "Pirate Harbor"
                            province = "CA"
                            country = "US"
                            postal_code = "95559-1227"
                            serial_number = "2"
                        }

                        dns_names = [
                            "example.com",
                            "example.net",
                        ]

                        ip_addresses = [
                            "127.0.0.1",
                            "127.0.0.2",
                        ]

                        uris = [
                            "spiffe://example-trust-domain/workload",
                            "spiffe://example-trust-domain/workload2",
                        ]

                        key_algorithm = "RSA"
                        private_key_pem = <<EOT
%s
EOT
                    }
                    output "key_pem_1" {
                        value = "${tls_cert_request.test1.cert_request_pem}"
                    }
                `, testPrivateKeyPEM),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["key_pem_1"].Value

					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"key_pem_1\" is not a string")
					}

					if !strings.HasPrefix(got, "-----BEGIN CERTIFICATE REQUEST----") {
						return fmt.Errorf("key is missing CSR PEM preamble")
					}
					block, _ := pem.Decode([]byte(got))
					csr, err := x509.ParseCertificateRequest(block.Bytes)
					if err != nil {
						return fmt.Errorf("error parsing CSR: %s", err)
					}
					if expected, got := "2", csr.Subject.SerialNumber; got != expected {
						return fmt.Errorf("incorrect subject serial number: expected %v, got %v", expected, got)
					}
					if expected, got := "example.com", csr.Subject.CommonName; got != expected {
						return fmt.Errorf("incorrect subject common name: expected %v, got %v", expected, got)
					}
					if expected, got := "Example, Inc", csr.Subject.Organization[0]; got != expected {
						return fmt.Errorf("incorrect subject organization: expected %v, got %v", expected, got)
					}
					if expected, got := "Department of Terraform Testing", csr.Subject.OrganizationalUnit[0]; got != expected {
						return fmt.Errorf("incorrect subject organizational unit: expected %v, got %v", expected, got)
					}
					if expected, got := "5879 Cotton Link", csr.Subject.StreetAddress[0]; got != expected {
						return fmt.Errorf("incorrect subject street address: expected %v, got %v", expected, got)
					}
					if expected, got := "Pirate Harbor", csr.Subject.Locality[0]; got != expected {
						return fmt.Errorf("incorrect subject locality: expected %v, got %v", expected, got)
					}
					if expected, got := "CA", csr.Subject.Province[0]; got != expected {
						return fmt.Errorf("incorrect subject province: expected %v, got %v", expected, got)
					}
					if expected, got := "US", csr.Subject.Country[0]; got != expected {
						return fmt.Errorf("incorrect subject country: expected %v, got %v", expected, got)
					}
					if expected, got := "95559-1227", csr.Subject.PostalCode[0]; got != expected {
						return fmt.Errorf("incorrect subject postal code: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(csr.DNSNames); got != expected {
						return fmt.Errorf("incorrect number of DNS names: expected %v, got %v", expected, got)
					}
					if expected, got := "example.com", csr.DNSNames[0]; got != expected {
						return fmt.Errorf("incorrect DNS name 0: expected %v, got %v", expected, got)
					}
					if expected, got := "example.net", csr.DNSNames[1]; got != expected {
						return fmt.Errorf("incorrect DNS name 0: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(csr.IPAddresses); got != expected {
						return fmt.Errorf("incorrect number of IP addresses: expected %v, got %v", expected, got)
					}
					if expected, got := "127.0.0.1", csr.IPAddresses[0].String(); got != expected {
						return fmt.Errorf("incorrect IP address 0: expected %v, got %v", expected, got)
					}
					if expected, got := "127.0.0.2", csr.IPAddresses[1].String(); got != expected {
						return fmt.Errorf("incorrect IP address 0: expected %v, got %v", expected, got)
					}

					if expected, got := 2, len(csr.URIs); got != expected {
						return fmt.Errorf("incorrect number of URIs: expected %v, got %v", expected, got)
					}
					if expected, got := "spiffe://example-trust-domain/workload", csr.URIs[0].String(); got != expected {
						return fmt.Errorf("incorrect URI 0: expected %v, got %v", expected, got)
					}
					if expected, got := "spiffe://example-trust-domain/workload2", csr.URIs[1].String(); got != expected {
						return fmt.Errorf("incorrect URI 1: expected %v, got %v", expected, got)
					}

					return nil
				},
			},
			{
				Config: fmt.Sprintf(`
                    resource "tls_cert_request" "test2" {
                        subject {
						serial_number = "42"
						}

                        key_algorithm = "RSA"
                        private_key_pem = <<EOT
%s
EOT
                    }
                    output "key_pem_2" {
                        value = "${tls_cert_request.test2.cert_request_pem}"
                    }
                `, testPrivateKeyPEM),
				Check: func(s *terraform.State) error {
					gotUntyped := s.RootModule().Outputs["key_pem_2"].Value

					got, ok := gotUntyped.(string)
					if !ok {
						return fmt.Errorf("output for \"key_pem_2\" is not a string")
					}

					if !strings.HasPrefix(got, "-----BEGIN CERTIFICATE REQUEST----") {
						return fmt.Errorf("key is missing CSR PEM preamble")
					}
					block, _ := pem.Decode([]byte(got))
					csr, err := x509.ParseCertificateRequest(block.Bytes)
					if err != nil {
						return fmt.Errorf("error parsing CSR: %s", err)
					}
					if expected, got := "42", csr.Subject.SerialNumber; got != expected {
						return fmt.Errorf("incorrect subject serial number: expected %v, got %v", expected, got)
					}
					if expected, got := "", csr.Subject.CommonName; got != expected {
						return fmt.Errorf("incorrect subject common name: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.Organization); got != expected {
						return fmt.Errorf("incorrect subject organization: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.OrganizationalUnit); got != expected {
						return fmt.Errorf("incorrect subject organizational unit: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.StreetAddress); got != expected {
						return fmt.Errorf("incorrect subject street address: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.Locality); got != expected {
						return fmt.Errorf("incorrect subject locality: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.Province); got != expected {
						return fmt.Errorf("incorrect subject province: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.Country); got != expected {
						return fmt.Errorf("incorrect subject country: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.Subject.PostalCode); got != expected {
						return fmt.Errorf("incorrect subject postal code: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.DNSNames); got != expected {
						return fmt.Errorf("incorrect number of DNS names: expected %v, got %v", expected, got)
					}
					if expected, got := 0, len(csr.IPAddresses); got != expected {
						return fmt.Errorf("incorrect number of IP addresses: expected %v, got %v", expected, got)
					}

					return nil
				},
			},
		},
	})
}
