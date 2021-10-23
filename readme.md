# What is it?

Эта штука для посылки сообщений в личный телеграм чатик через бота.

Читает сообщение из stdin

```text
Usage of paranotify:
  -reply-server
        serve messages and reply user-id
  -token string
        telegram token
  -user-id string
        user who should receive msg
example: 
        paranotify -token=${TOKEN} -reply-server
        echo this is echo msg | paranotify -token=${TOKEN} -user-id=${USER}
```
