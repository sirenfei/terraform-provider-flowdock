provider "flowdock" {
  version = "1.1.2"
  api_token = "xxxxxx"
}
resource "flowdock_invitation" "gyles-ops" {
  count = "${length(var.user_flows)}"
  email = "mickey.mouse@gmail.com"
  flow = "${var.user_flows[count.index]}"
  org = "smart-mouse"
  message = "Mickey Mouse"
}

resource "flowdock_invitation" "aaaaa-robot_1_test-terraform" {
    org = "test-terraform"
    flow = "flow1"
    email = "aaaaa.robot@gmail.com"
    message = "hello aaaa"
}
