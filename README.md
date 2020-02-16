# telegram_go_bot

[![Coverage Status](https://coveralls.io/repos/github/azerg/telegram_go_bot/badge.svg?branch=master)](https://coveralls.io/github/azerg/telegram_go_bot?branch=master)

[![Build Status](https://travis-ci.com/azerg/telegram_go_bot.svg?branch=master)](https://travis-ci.com/azerg/telegram_go_bot)

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

