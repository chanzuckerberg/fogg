data "external" "git_sha" {
  program = [
    "git",
    "log",
    "--pretty=format:{ \"sha\": \"%H\" }",
    "-1",
    "HEAD"
  ]
}

output "git_sha" {
  value = try(data.external.git_sha.result.sha, "blah")
