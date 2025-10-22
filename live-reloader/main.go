package main

import (
	"crypto/sha256"
	"io/fs"
	"os"
	fp "path/filepath"
	"sync"
	"time" 
	"log"
)

// Utility function for handling errors
func check(err error) {
	if err != nil {
		log.Println("Error Occured: ", err)
	}
}

// # LiveReloader
// 
// dirpath: is the path to dir to be watched
//
// The filehashes: concurrent hashmap will
// have the key as the file path and
// value as the file content hash
//
// # This is a live reloader application
//
// Will scan through files in a directory or anything
// changes with in the directory then reloads in real
// time, We'll spawn goroutines to look at each file
// concurrently, then we'll use a concurrent hashmap
// Where the key is the filepath and the value is the
// hash of the contents in the given file
type LiveReloader struct {
	dirpath    string
	changes    chan string
	filehashes *sync.Map
}

// Initialize the LiveReloader
func New(dir string) LiveReloader {
	// Initialize the filehashes concurrent hashmap
	var filehashes sync.Map
	changes := make(chan string, 10)

	return LiveReloader{
		dirpath:    dir,
		changes:    changes,
		filehashes: &filehashes,
	}
}

// Returns the dir where the live reloading will happen
func (livereloader *LiveReloader) GetDir() string {
	return livereloader.dirpath
}

// Read the changes channel for changed files with evidence of their file paths
// Takes in a function and executes it each time a file's content is changed
func (livereloader *LiveReloader) GetFileChangesFromChannel(changefunc func(filepath string)) {
	for {
		filepath := <-livereloader.changes
		changefunc(filepath)
	}
}

// Goroutine that will read file, hash the contents and then
// cache the its hash in a concurrent map where key is the 
// path to file, and the value is its content hash, 
// using sha256 for the hashing
func MonitorFileChanges(livereloader *LiveReloader, path string) {
	// Read the file
	c, err := os.ReadFile(path)
	check(err)

	// The hashing step
	h := sha256.New()
	h.Write(c)
	bs := h.Sum(nil)
 
	// Check if the key exists in concurrent map
	hashvalue, exists := livereloader.filehashes.Load(path)

	// Check if the previous hash is the same, and hashvalue from map
	// Is equal to current content hash...
	if exists {
		if hashvalue != string(bs) {
			livereloader.changes <- path
		}
	}

	// Store the path and hash in concurrent hashmap
	livereloader.filehashes.Store(path, string(bs))
}

// This will be used in the filepath dir walker
func VisitDir(livereloader *LiveReloader, path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Check whether the given path is a directory or not
	if !d.IsDir() {
		// Spawn a goroutine for each file path
		go MonitorFileChanges(livereloader, path)
	}

	return nil
}

func main() {
	// Intialize the LiveReloader
	livereloader := New("../")

	// Start the changes reader
	go livereloader.GetFileChangesFromChannel(func(filepath string) {
		log.Println("Changed: ", filepath)
		log.Println("Rerunning application")
	})
 
	// Initialize the ticker
	ticker := time.NewTicker(time.Second * 1)
	done := make(chan bool)

	log.Println("Starting the Live Reloader")

	// Run the live reloader
	// Add a ticker with an interval of 1 second
	for {
		select {
		case <-done:
			continue
		case <-ticker.C:
			// Start the walkdir
			fp.WalkDir(livereloader.GetDir(), func(path string, d fs.DirEntry, err error) error {
				return VisitDir(&livereloader, path, d, err)
			})
		}
	}
}
