## Todo with HTMX, Go Backend, and PostgreSQL in Docker

In this project, I am developing a Todo application that leverages the power of HTMX for dynamic user interactions. The backend is implemented in Go, integrated with a Dockerized PostgreSQL instance.

### HTMX

[HTMX](https://htmx.org/) is a library that facilitates dynamic web applications by allowing you to update parts of a page with data from the server, providing a smoother and more interactive user experience without the need for complex JavaScript frameworks.

### Go

[Go](https://golang.org/), often referred to as Golang, is a programming language designed for simplicity, efficiency, and reliability. It excels in creating fast and scalable server-side applications, making it well-suited for building robust and high-performance software systems.

### PostgreSQL

[PostgreSQL](https://www.postgresql.org/) is an open-source relational database management system known for its reliability and extensibility.

### Docker

[Docker](https://www.docker.com/) is an open platform for developing, shipping, and running applications. Docker enables you to separate your applications from your infrastructure so you can deliver software quickly

## Usage

```console
$ git clone https://github.com/mahonyodhran/go-htmx-todo.git
```

I have a .env file locally which contains DB_CONN, so create this first.

```console
$ cd go-htmx-todo
$ touch .env
```
Add DB_CONN="< your connection string >"
