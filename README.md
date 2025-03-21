## Documentation

### Go API Resource Project

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

This repository contains a collection of Go programs and libraries that
demonstrate the language, standard libraries, and tools.

📌 Prasyarat
Sebelum memulai, pastikan Anda memiliki:

Golang (Minimal versi 1.24)

Git

Database MySQL

### Run Project di Local

Clone terlebih dahulu project ini

```bash
  git clone https://github.com/ddfrmnsh-dev/test-be-golang.git
```

Masuk kedalam direktori

```bash
  cd test-be-golang
```

Install depedensi go

```bash
  go mod tidy
```

Jalankan server

```bash
  go run .
```

## Environment Variables

Untuk menjalankan project ini, anda harus menyesuaikan DBConfig dan APIConfig di dalam `config/config.go` file

`Port` Default Port 3306

`Database` Nama Database

`Username` Database

`Password` Database

Sesuaikan dengan Environment lokal anda terkait koneksi database.

Setelah itu cek port Golang jika `ApiPort: "8888"` sudah terpakai pada lokal anda atau sedang berjalan maka ganti port tersebut.

## API Documentation

Silahkan Import Postman Collection untuk mengetes API.
test-be-ordent\API Backend Ordent.postman_collection.json

Berikut design arsitektur project dan ERD
![Screenshot](https://github.com/ddfrmnsh-dev/test-be-ordent/blob/main/image/ERD%20Library.drawio.png)
