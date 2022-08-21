data "sops-content" "test" {
  content = file("test-fixtures/test.yml")
  format = "yaml"
}

locals {
  decrypted         = data.sops-content.test.decrypted
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
