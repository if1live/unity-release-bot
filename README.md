# Unity Release Bot

Twitter bot to notify new unity release. https://twitter.com/unityreleasebot
[![Build Status](https://travis-ci.org/if1live/unity-release-bot.svg?branch=master)](https://travis-ci.org/if1live/unity-release-bot)

## Features
* Read RSS Feed
  * Patch release : https://unity3d.com/unity/qa/patch-releases/latest.xml
  * beta : https://unity3d.com/unity/beta/latest.xml
* Read stable release from download page



## Daemon
use bot.service

```
$ sudo systemctl daemon-reload
$ sudo systemctl stop unity-release-bot
$ sudo systemctl start unity-release-bot
```
