package reverseindex

import "os"

func CleanupDb() {
	os.RemoveAll("./reverse_index_data")
	os.RemoveAll("./pos_index_data")
	os.RemoveAll("./rev")
}
