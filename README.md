## Kora

> Backend repository for kora (web app). 

### Introduction
Kora is a side project i developed for learning purposes. It's a clone of the popular question/answer thread app [Quora](https://quora.com). This is the backend repository built with Go/fiber & MongoDB, deployed on heroku. The [front-end repository](https://github.com/Mayowa-Ojo/kora-client) is hosted separately.

### Features
While this is not a 100% working clone as the app is pretty large with a lot of moving parts, I tried to implement most of the core features of quora.

- Email and password authentication
- Ask/answer questions
- Create spaces, topics, comments
- Edit profile
- Follow/unfollow users, spaces
- Share posts
- Basic suggestions logic
- Simple shorturl service
- Role-based access

### Development
Requires Go >=1.13 & MongoDB >=4.2
```shell
$ mkdir <folder>
$ git clone https://github.com/Mayowa-Ojo/kora-server.git .
$ touch .env
$ go run server.go || air 

```

I do intend to keep developing this app (adding new features and fixing bugs) as I enjoyed building this one.

### Bug report
If you find any bugs in this app, [open an issue](https://github.com/Mayowa-Ojo/kora-server/issues/new)