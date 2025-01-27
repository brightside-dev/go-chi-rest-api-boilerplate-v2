# 🥶 Golang + Chi Web Application Template
The goal of this project is threefold: to learn Golang, to make it the foundation for all my future side projects, and to level up as a programmer.

I chose Golang for my future side projects because working with PHP, which I use daily at my job, has started to feel more like a chore than an enjoyable challenge. 

I've tried a bunch of languages and frameworks over the years, but Golang is the one I’ve had the most fun with.

This project was initially bootstrapped using https://github.com/Melkeydev/go-blueprint

### Setup
1. Install MySQL and create a database
2. Create a .env file and modify DB credentials
3. Run DB migrations `goose up`
4. Start server `go run ./cmd/api `

### 🏗️ TODO
1. Implement Mailgun email service (Mailcatcher for dev)
2. Implement Firebase push notifications service
3. Implement web sockets service (Pusher, PubNub etc.)
4. Implement tests
5. Implement CI/CD using Github

### 🛞 Refactor (when time permits)
* Email service - figure out a good way to create services
* Environment variables

### 🧰 Project Tools & Packages
* MySQL driver: https://github.com/go-sql-driver/mysql
* Router & Middleware: https://github.com/go-chi/chi
* Goose for DB migrations: https://github.com/pressly/goose
* Godotenv for env variables: https://github.com/joho/godotenv
* Cobra for easy command management: https://github.com/spf13/cobra

### 🚀 Features
**1. REST API**
    <br> * Controller(handler)/Repository pattern for separation of concern with minimal complexity
    <br> * JWT Authentication: https://github.com/go-chi/jwtauth

**2. Commands Service**
    <br> * With the use of Cobra, we can integrate a command-line interface (CLI) into our Go web application, allowing us to run custom processes directly from the terminal or programmatically through code. This is especially useful for executing asynchronous or synchronous tasks, either within the same goroutine or in separate goroutines, depending on the application's requirements. Also useful for running certain commands as cron jobs.

**3. Admin CMS Dashboard**
    <br>* Session authentication and manager: https://github.com/alexedwards/scs
    <br>* CMS dashboard theme: https://github.com/pro-dev-ph/bootstrap-simple-admin-template
    <br>* Minimal JS to handle necessary CMS animations, charts and data tables

**4. Emailing** WIP

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

### Useful Resources
* https://go.dev/doc/effective_go
* https://pkg.go.dev/std
* https://lets-go.alexedwards.net/
* https://www.youtube.com/@MelkeyDev
* https://www.youtube.com/@anthonygg_