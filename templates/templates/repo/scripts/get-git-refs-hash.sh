# Created by Cursor
#!/bin/bash
set -euo pipefail

# Extract Git module sources with refs from fogg.yml and resolve them to commit SHAs
# This creates a hash that changes when any Git ref points to a different commit

cd "$(dirname "$0")/.."

# Extract all Git module sources with refs from fogg.yml
git_refs=$(grep -E "module_source:.*\?ref=" fogg.yml | \
    sed -E 's/.*module_source: *//' | \
    grep -E "\?ref=" | \
    sort -u)

if [ -z "$git_refs" ]; then
    echo "No Git refs found in fogg.yml"
    exit 1
fi

# For each Git ref, resolve to commit SHA
commit_shas=""
while IFS= read -r git_ref; do
    if [ -n "$git_ref" ]; then
        # Parse the repository URL and ref
        full_ref="$git_ref"
        repo_url=$(echo "$full_ref" | cut -d'?' -f1)
        ref=$(echo "$full_ref" | sed -E 's/.*\?ref=//')
        
        # Convert SSH URLs to HTTPS for git ls-remote
        if [[ "$repo_url" == git@github.com:* ]]; then
            repo_url=$(echo "$repo_url" | sed 's/git@github.com:/https:\/\/github.com\//')
        elif [[ "$repo_url" == github.com/* ]]; then
            repo_url="https://$repo_url"
        fi
        
        # Remove the module path part for repository URL (everything after the repo name //)
        if [[ "$repo_url" == *"//"* ]] && [[ "$repo_url" == *"github.com"* ]]; then
            # For GitHub URLs, keep everything up to the repository name
            repo_url=$(echo "$repo_url" | sed -E 's|^(https://github.com/[^/]+/[^/]+)//.*|\1|')
        fi
        
        echo "Resolving $repo_url ref $ref..." >&2
        
        # Use git ls-remote to get the commit SHA for this ref
        sha=""
        
        # Try as a tag first
        sha=$(git ls-remote "$repo_url" "refs/tags/$ref" 2>/dev/null | head -n1 | cut -f1 || true)
        
        # If tag doesn't exist, try as a branch
        if [ -z "$sha" ]; then
            sha=$(git ls-remote "$repo_url" "refs/heads/$ref" 2>/dev/null | head -n1 | cut -f1 || true)
        fi
        
        # If still empty, try without refs prefix (for some edge cases)
        if [ -z "$sha" ]; then
            sha=$(git ls-remote "$repo_url" "$ref" 2>/dev/null | head -n1 | cut -f1 || true)
        fi
        
        if [ -n "$sha" ]; then
            commit_shas="${commit_shas}${repo_url}=${ref}=${sha}\n"
            echo "  -> $sha" >&2
        else
            echo "  -> Could not resolve ref $ref in $repo_url" >&2
            exit 1
        fi
    fi
done <<< "$git_refs"

# Create a hash of all the commit SHAs
echo -e "$commit_shas" | sha256sum | cut -d' ' -f1
