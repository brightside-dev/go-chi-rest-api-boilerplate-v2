# 🥶 Golang + Chi Web Application Template
The goal of this project is threefold: to learn Golang, to make it the foundation for all my future side projects, and to level up as a programmer.

I chose Golang for my future side projects because working with PHP, which I use daily at my job, has started to feel more like a chore than an enjoyable challenge. 

I've tried a bunch of languages and frameworks over the years, but Golang is the one I’ve had the most fun with.

This project was initially bootstrapped using https://github.com/Melkeydev/go-blueprint

### 🏗️ TODO
1. Implement better project configuration & env variables management
2. Implement DB error logging 
3. Implement cronjob service - implement something like Symfony's command package
4. Implement Mailgun email service
5. Implement Firebase push notifications service
6. Implement web sockets service (Pusher, PubNub etc.)
7. Implement tests

### 🧰 Project Tools & Packages
* MySQL
* Router & Middleware: https://github.com/go-chi/chi
* Database migrations are handled by Goose: https://github.com/pressly/goose
*

### 🚀 Features
**1. REST API**
    * JWT Authentication: https://github.com/go-chi/jwtauth

**2. Admin CMS Dashboard**
    * Session authentication and manager: https://github.com/alexedwards/scs
    * CMS dashboard theme: https://github.com/pro-dev-ph/bootstrap-simple-admin-template
    * Minimal JS to handle necessary CMS animations, charts and data tables

### ⛩️ Folder Structure
```
├── cmd/
│   └── main.go --entry point
├── internal/
│   ├── container/
│   │   └── container.go       --dependency injection container
│   ├── database/
│   │   └── database.go        --db service
│   ├── handler/               --handlers aka controllers
│   │   ├── user_handler.go
│   ├── repository/            --repositories
│   │   ├── user_repository.go
│   ├── server/                
│   │   ├── routes.go          --routes
│   │   └── server.go          --server service
│   ├── template/
│   │   └── template.go        --template service
│   └── model/
│       └── model.go           --entities
├── ui/                        --cms views (html, css, js)
│   ├── assets/
│   │   ├── css/
│   │   ├── img/
│   │   ├── js/
│   │   ├── vendor/
│   ├── html/
│   │   ├── dashboard/
│   │   ├── partials/
│   │   └── base.html
```