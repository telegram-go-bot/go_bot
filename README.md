# BillyBot - telegram bot written in golang

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Coverage Status](https://coveralls.io/repos/github/telegram-go-bot/go_bot/badge.svg?branch=master)](https://coveralls.io/github/telegram-go-bot/go_bot?branch=master) [![Build Status](https://travis-ci.com/telegram-go-bot/go_bot.svg?branch=master)](https://travis-ci.com/telegram-go-bot/go_bot) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com) [![Go Report Card](https://goreportcard.com/badge/github.com/telegram-go-bot/go_bot)](https://goreportcard.com/report/github.com/telegram-go-bot/go_bot) [![HitCount](http://hits.dwyl.com/azerg/githubcom/telegram-go-bot/go_bot.svg)](http://hits.dwyl.com/azerg/githubcom/telegram-go-bot/go_bot) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/f7237471effa4ee5a07cd85447eaa2e6)](https://www.codacy.com/gh/telegram-go-bot/go_bot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=telegram-go-bot/go_bot&amp;utm_campaign=Badge_Grade)[![Test Coverage](https://api.codeclimate.com/v1/badges/dfe7d60958d6a1c7ef5b/test_coverage)](https://codeclimate.com/github/telegram-go-bot/go_bot/test_coverage)[![Maintainability](https://api.codeclimate.com/v1/badges/dfe7d60958d6a1c7ef5b/maintainability)](https://codeclimate.com/github/telegram-go-bot/go_bot/maintainability)
[![DeepSource](https://deepsource.io/gh/telegram-go-bot/go_bot.svg/?label=active+issues&show_trend=true&token=AF1566g708U-Z6yFOl9NPj90)](https://deepsource.io/gh/telegram-go-bot/go_bot/?ref=repository-badge)

funny telegram bot in go.

### Environmental variables required:

| Key  | Value (description)  |
| ------------ | ------------ |
|BOT_UIDS|list of bot names, separated by ",". Is used to let bot know when he is triggered|
|GOOGLE_SEARCH_API_KEY|aBc-DeFgHiJkl...|
|GOOGLE_SEARCH_ENGINE_ID|1234:abcde...|
|HEROKU_BOT_ID| telegram bot id|
|HEROKU_BASE_URL|public url of bots service. For debugging i like to put tunneled ngroks url here|
|DATABASE_URL|postgres db external url|
|VK_LOGIN|~|
|VK_PASSWORD|~|
|OWNER_ID|Telegram username. is  used to give unique replies for example for debugging|
|PORT|~|

### Todo:
  - [ ] Add UML diagram
  - [x] enable coverity
  - [ ] Fix coveralls
  - [ ] + db
