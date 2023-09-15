data "external" "git_sha" {
  try(
  program = [
    "git",
    "hash-object",
    "*",
    "--pretty=format:{ \"sha\": \"%H\" }",
    "-1",
    "HEAD"
  ], "blah")
}

output "git_sha" {
  value = try(data.external.git_sha.result.sha, "blah")
