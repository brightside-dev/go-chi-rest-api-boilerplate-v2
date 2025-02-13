# 🥶 Golang + Chi Web Application Template
## Note: I have now stopped updating this project. I'm currently building a side project using this structure. I may update this in the future based on my learnings on the other project.
The goal of this project is threefold: to learn Golang, to make it the foundation for all my future side projects, and to level up as a programmer.

I chose Golang for my future side projects because working with PHP, which I use daily at my job, has started to feel more like a chore than an enjoyable challenge. 

I've tried a bunch of languages and frameworks over the years, but Golang is the one I’ve had the most fun with.

This project was initially bootstrapped using https://github.com/Melkeydev/go-blueprint

### Setup
1. Install MySQL and create a database
2. Create a .env file and modify DB credentials
3. Run DB migrations `goose up`
4. Start server `go run ./cmd/api`


### 🧰 Project Tools & Packages
* MySQL driver: https://github.com/go-sql-driver/mysql
* Router & Middleware: https://github.com/go-chi/chi
* Goose for DB migrations: https://github.com/pressly/goose
* Godotenv for env variables: https://github.com/joho/godotenv
* Cobra for easy command management: https://github.com/spf13/cobra
* Mailgun for email sends: https://github.com/mailgun/mailgun-go
* APN for iOS push notifications: https://github.com/sideshow/apns2
* Firebase Messaging for Android push notifications: https://firebase.google.com/go

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

**4. Emails** 
    <br>* Local SMTP using Mailcatcher for easy local email development
    <br>* Mailgun service for production email sending

**5. Push Notifications**
    <br>* APN for iOS push notifications
    <br>* Firebase FCM for android push notifications - might refactor to use Firebase to handle all push notifications

**6. Database Logging**
    <br>* DB logger using log/slog

### Useful Resources
* https://go.dev/doc/effective_go
* https://pkg.go.dev/std
* https://lets-go.alexedwards.net/
* https://www.youtube.com/@MelkeyDev
* https://www.youtube.com/@anthonygg_