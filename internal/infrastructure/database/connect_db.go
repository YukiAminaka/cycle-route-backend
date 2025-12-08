package database

import "github.com/YukiAminaka/cycle-route-backend/internal/infrastructure/database/dbgen"


var (
	query     *dbgen.Queries
)

func GetQuery() *dbgen.Queries {
	return query
}

func SetQuery(q *dbgen.Queries) {
	query = q
}