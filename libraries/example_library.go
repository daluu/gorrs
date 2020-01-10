package libraries

import (
	"errors"
	"fmt"
	"io/ioutil"
)

//ExampleRemoteLibrary to be used with Robot Framework's remote server.
type ExampleRemoteLibrary struct{}

//CountItemsInDirectory the number of items in the directory specified by `path`.
func (lib *ExampleRemoteLibrary) CountItemsInDirectory(path string) (int, error) {
	fileCount := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fileCount, err
	}
	fileCount = len(files)
	return fileCount, err
}

//StringsShouldBeEqual ...
func (lib *ExampleRemoteLibrary) StringsShouldBeEqual(str1 string, str2 string) error {
	fmt.Printf("Comparing '%s' to '%s'.", str1, str2)
	if str1 != str2 {
		return errors.New("Given strings are not equal.")
	} else {
		return nil
	}
}

//optional extra keyword below, following phrrs (PHP robot framework remote server)
//comment out if it interferes with running example remote library tests against gorrs

//TruthOfLife ...
func (lib *ExampleRemoteLibrary) TruthOfLife() int {
	return 42
}

//TruthOfLife ...
func (lib *ExampleRemoteLibrary) ReturnArray() []interface{} {
	var testArray []interface{}
	testArray = append(testArray, "string")
	testArray = append(testArray, 1)
	testArray = append(testArray, 1.1)
	return testArray
}
