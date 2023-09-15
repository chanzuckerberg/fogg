data "external" "git_sha" {
  program = [
    "./soft_git_log"
  ]
}

output "git_sha" {
  value = data.external.git_sha.result.sha
}