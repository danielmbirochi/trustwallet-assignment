package inmemorydb_test

import (
	"bytes"
	"testing"

	"github.com/danielmbirochi/trustwallet-assignment/internal/state/inmemorydb"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestInMemoryDB(t *testing.T) {
	t.Run("KeyValueStorerOperations", func(t *testing.T) {
		db := inmemorydb.New()
		defer db.Close()

		key := "0x388c818ca8b9251b393131c08a736a67ccb19297"
		value := []byte("some value")

		{
			testID := 0
			if got, err := db.Has(key); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check if the key-value store has an entry : %s", Failed, testID, err)
			} else if got {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check if the key-value store has an entry : Expected false. Got %t", Failed, testID, got)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to check if the key-value store has an entry", Success, testID)
		}

		{
			testID := 1
			if err := db.Put(key, [][]byte{value}); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to write an entry in the key-value store : %s", Failed, testID, err)
			}

			if got, err := db.Has(key); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check for a written key : %s", Failed, testID, err)
			} else if !got {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check for a written key : Expected true. Got %t", Failed, testID, got)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to check for a written key", Success, testID)
		}

		{
			testID := 2
			if got, err := db.Get(key); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able query an entry in the key-value store : %s", Failed, testID, err)
			} else if !bytes.Equal(got[0], value) {
				t.Fatalf("\t%s\tTest %d:\tShould be able query an entry in the key-value store : Expected %s. Got %s", Failed, testID, value, got[0])
			}
			t.Logf("\t%s\tTest %d:\tShould be able query an entry in the key-value store", Success, testID)
		}

		{
			testID := 3
			if err := db.Delete(key); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete an entry in the key-value store : %s", Failed, testID, err)
			}

			if got, err := db.Has(key); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check if the key-value store has an entry : %s", Failed, testID, err)
			} else if got {
				t.Fatalf("\t%s\tTest %d:\tShould be able to check if the key-value store has an entry : Expected false. Got %t", Failed, testID, got)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete an entry in the key-value store", Success, testID)
		}
	})
}
