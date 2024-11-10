package reverseindex

import "os"

func CleanupDb() {
	os.RemoveAll("./reverse_index_data")
}
