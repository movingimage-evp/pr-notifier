# PR notifier
Use GitHub Actions to notify Slack that a pull request is opened since `DAYS_BEFORE` days.

## Usage

Add the following YAML to your new GitHub Actions workflow:
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

### Arguments

#### SLACK_API_KEY
The Slack api key. You'll need to create a new Slack bot app which will generate you a token.

#### SLACK_CHANNEL_ID
The Slack channel ID where your notifications should appear.

#### DAYS_BEFORE
The amount days how old your PR is. 


