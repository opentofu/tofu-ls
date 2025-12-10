terraform {
  required_providers {
    test = {
      source = "test/test"
    }
  }
  required_version = "1.11.0"
}

ephemeral "" "eph1" {

}
