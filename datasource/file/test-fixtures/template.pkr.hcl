data "sops-file" "test" {
  file = "test-fixtures/test.yml"
  format = "yaml"
}

locals {
  decrypted         = data.sops-file.test.decrypted
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo decrypted value: ${local.decrypted}",
    ]
  }
}
