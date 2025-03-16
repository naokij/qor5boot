package main

import (
	"github.com/naokij/qor5boot/admin"
)

func main() {
	db := admin.ConnectDB()
	tbs := admin.GetNonIgnoredTableNames(db)
	admin.EmptyDB(db, tbs)
	admin.InitDB(db, tbs)
	return
}
