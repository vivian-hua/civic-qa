module github.com/team-ravl/civic-qa/services/logAggregator

go 1.15

require (
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/team-ravl/civic-qa/services/common v0.0.0
	github.com/urfave/negroni v1.0.0
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.3
)

replace github.com/team-ravl/civic-qa/services/common v0.0.0 => ../common
