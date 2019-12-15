package handlers

import "github.com/riyadennis/identity-server/internal/store/sqlite"

var (
	Idb *sqlite.IdentityDB
)

func Init() {
	Idb = sqlite.PrepareDB("/var/tmp/identity.db")
}
