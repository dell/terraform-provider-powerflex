terraform {
  required_version = "> 0.8.0"
}

resource "null_resource" "health_check" {

connection {
    type     = "ssh"
    user     = "XXXX"
    password = "XXXX"
    host     = "XXXX"
  }


provisioner "file" {
    source      = "gateway_installer.sh"
    destination = "/tmp/gateway_installer.sh"
}

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/gateway_installer.sh",
      "/tmp/gateway_installer.sh",
    ]
  }
}