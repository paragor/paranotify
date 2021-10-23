# What is it?

Эта штука для посылки сообщений в личный телеграм чатик через бота.

Читает сообщение из stdin

```text
Usage of paranitfy:
  -reply-server
        serve messages and reply user-id
  -token string
        telegram token
  -user-id string
        user who should receive msg
example: 
        paranitfy -token=${TOKEN} -reply-server
        echo this is echo msg | paranitfy -token=${TOKEN} -user-id=${USER}
```
