# PR notifier
A GitHub action that shames PR's that are open longer `DAYS_BEFORE` days

## Usage
```
name: PR notifier

on:
  schedule:
    - cron: "30 10 * * 1-5"

jobs:
  pr-notify:
    name: Notify old PRs
    runs-on: ubuntu-latest
    steps:
      - name: notify old prs
        uses: movingimage-evp/pr-notifier@v1
        env:
          SLACK_API_KEY: REPLACE_ME
          SLACK_CHANNEL_ID: REPLACE_ME
          DAYS_BEFORE: -2
          GITHUB_TOKEN: ${{ github.token }}
```
