
# provider.tf

# Specify the provider and access details
provider "vault" {
  address = "https://vault.hashi.co"
}

provider "aws" {
  profile                 = "your-account"
  region                  = var.aws_region

  access_key = data.vault_aws_access_credentials.telemetry.access_key
  secret_key = data.vault_aws_access_credentials.telemetry.secret_key
  token      = data.vault_aws_access_credentials.telemetry.security_token
}
