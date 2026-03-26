# Url-Shorterer
Simple web app to practice gin framework and clear architecture

# Motivation
I thought about how this kind of services working, so i decided to create my own url shortener

# Tech
* Go (gin)
* Postgres
* clean architecture

# Installation and start
1. Clone my repo
```
git clone https://github.com/Fista6k/UrlShortener.git

cd UrlShortener
```

2. Create .env file
```
//example .env
DB_NAME=your data base name
DB_USER=your data base user
DB_PASSWORD=your data base password
DB_PORT=your data base port
DB_HOST=localhost
```

3. Run the server
```
go run cmd/main.go
```

4. Check it out
  [http://localhost:8080](http://localhost:8080)
