package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// File defines a config file
type File struct {
	path       string
	userConfig interface{}
	watcher    *fsnotify.Watcher
	updated    func()
}

// New creates a new config file and starts watching for changes
// filepath is the JSON formatted file to monitor
// c is the interface to attempt to marshal the file into
// updated is called when there are updates to the file
func New(filepath string, c interface{}, updated func()) (*File, error) {
	f := &File{path: filepath, userConfig: c, updated: updated}
	go f.watch(filepath)

	return f, f.loadData()
}

// Close the FileConfig and remove all watchers
func (f *File) Close() {
	f.watcher.Close()
}

// load the data from the config into the defined structure
func (f *File) loadData() error {
	cf, err := os.OpenFile(f.path, os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
	defer cf.Close()

	jd := json.NewDecoder(cf)
	jd.Decode(f.userConfig)

	return nil
}

// watch a file for changes
func (f *File) watch(filepath string) {
	// creates a new file watcher
	var err error
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("ERROR", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					err := f.loadData()
					if err != nil {
						log.Println("error", err)
					}

					if f.updated != nil {
						f.updated()
					}
				}
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = f.watcher.Add(filepath)
	if err != nil {
		log.Println("ERROR", err)
	}

	<-done
}
