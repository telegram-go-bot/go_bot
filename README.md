# telegram_go_bot

[![Coverage Status](https://coveralls.io/repos/github/telegram-go-bot/go_bot/badge.svg?branch=master)](https://coveralls.io/github/telegram-go-bot/go_bot?branch=master) [![Build Status](https://travis-ci.com/telegram-go-bot/go_bot.svg?branch=master)](https://travis-ci.com/telegram-go-bot/go_bot) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com) [![Go Report Card](https://goreportcard.com/badge/github.com/telegram-go-bot/go_bot)](https://goreportcard.com/report/github.com/telegram-go-bot/go_bot) [![HitCount](http://hits.dwyl.com/azerg/githubcom/telegram-go-bot/go_bot.svg)](http://hits.dwyl.com/azerg/githubcom/telegram-go-bot/go_bot) [![Coverage Status](https://coveralls.io/repos/github/telegram-go-bot/go_bot/badge.svg?branch=master)](https://coveralls.io/github/telegram-go-bot/go_bot?branch=master) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/f7237471effa4ee5a07cd85447eaa2e6)](https://www.codacy.com/gh/telegram-go-bot/go_bot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=telegram-go-bot/go_bot&amp;utm_campaign=Badge_Grade)

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
|PHOTO_TO_APP_ID|pho.to dev app id|
|PHOTO_TO_KEY| pho.to dev app key|
|VK_LOGIN|~|
|VK_PASSWORD|~|
|OWNER_ID|Telegram username. is  used to give unique replies for example for debugging|
|PORT|~|

### Todo:
 - [ ] Add UML diagram
 - [x] enable coverity
