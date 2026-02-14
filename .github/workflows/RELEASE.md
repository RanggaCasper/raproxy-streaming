# Release Guide – raproxy-streaming

This document explains how to create a release that will automatically:

* Build multi-platform binaries (Linux / Windows / macOS)
* Generate checksum files
* Upload binaries to GitHub Releases
* (Optional) Deploy to VPS

---

## 1. Preparation

Before creating a release, make sure the project is stable.

Run:

```bash
make test
make build
```

Ensure there are no errors.

---

## 2. Creating a Release Tag

The GitHub Actions workflow will run automatically when you push a tag with the following format:

```
vX.Y.Z
```

Example:

```bash
git add .
git commit -m "feat: add new feature"
git push origin main

git tag v1.0.0
git push origin v1.0.0
```

---

## 3. What Happens After Pushing a Tag

Once the tag is pushed, GitHub Actions will:

1. Build binaries for:

   * linux/amd64
   * linux/arm64
   * windows/amd64
   * darwin/amd64
   * darwin/arm64

2. Generate files with the following format:

   ```
   raproxy-streaming_v1.0.0_linux_amd64
   raproxy-streaming_v1.0.0_linux_amd64.sha256
   ```

3. Create a GitHub Release

4. Upload all binaries to the Release page

---

## 4. Downloading the Binary

Go to:

```
Repository → Releases
```

Download the binary that matches your operating system.

Example (Linux):

```bash
wget https://github.com/<username>/raproxy-streaming/releases/download/v1.0.0/raproxy-streaming_v1.0.0_linux_amd64
chmod +x raproxy-streaming_v1.0.0_linux_amd64
./raproxy-streaming_v1.0.0_linux_amd64
```

---

## 5. Creating a New Version

To release a new version:

```bash
git tag v1.0.1
git push origin v1.0.1
```

---

## 6. Deleting a Release (If Needed)

Delete the local tag:

```bash
git tag -d v1.0.0
```

Delete the remote tag:

```bash
git push origin :refs/tags/v1.0.0
```

Then manually delete the Release from the GitHub Releases page.

---

## 7. Deploying to VPS (Manual)

If the deploy workflow is available:

1. Go to the **Actions** tab
2. Select the **Deploy to VPS** workflow
3. Click **Run workflow**
4. Provide:

   * VPS IP address
   * SSH username
   * SSH password
   * Service name

---

## 8. Best Practices

Recommended for production:

* Use SSH keys instead of passwords
* Store credentials in GitHub Secrets
* Do not commit `.env` files
* Follow Semantic Versioning

Version format:

```
vMAJOR.MINOR.PATCH
```

Examples:

* v1.0.0 → initial release
* v1.1.0 → new feature
* v1.1.1 → bug fix

---

## Quick Summary

To create a new release:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

GitHub Actions will automatically build and publish the release.