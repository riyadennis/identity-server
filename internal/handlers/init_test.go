package handlers_test

import (
	"github.com/riyadennis/identity-server/internal/store/sqlite"
	"os"
	"testing"

	"github.com/riyadennis/identity-server/internal"
)

func TestMain(m *testing.M){
	err := sqlite.Setup("/var/tmp/identityTest.db")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	internal.Server(":8085")
	os.Exit(m.Run())
}

