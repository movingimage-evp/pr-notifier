# PR notifier
A Github action that shames PR's that are open longer than 2 days

## usage
```
name: Shame Bot

on:
  schedule:
    - cron: "0 6 * * *"

jobs:
  shame:
    name: Shame Old PRs
    runs-on: ubuntu-latest
    steps:
      - name: shame old prs
        uses: movingimage-evp/pr-notifier@v1
        env:
          GITHUB_TOKEN: ${{ github.token }}
```
