package drivers

import (
	"encoding/json"
	"fmt"
	"github.com/jcelliott/lumber"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const version = "1.0.0" // Version of the driver package

type (
	// Logger interface for logging messages
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Debug(string, ...interface{})
		Info(string, ...interface{})
		Trace(string, ...interface{})
	}

	// Driver struct manages database operations
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

// Connect initializes a new Driver instance with the default directory
func Connect() (*Driver, error) {
	dir := "./" // Default directory for database storage

	db, err := New(dir, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// New creates a new Driver instance with optional settings
func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}

	// Apply provided options if available
	if options != nil {
		opts = *options
	}

	// Default logger setup
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	// Check if the directory already exists
	if _, err := stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'...\n", dir)

	return &driver, os.MkdirAll(dir, 0755) // Create the directory if it doesn't exist
}

// Write saves a record to the specified collection and resource
func (d *Driver) Write(collection string, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection name - no place to save data!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource name - unable to save record!")
	}

	resource = strings.ReplaceAll(resource, " ", "_") // Replace spaces with underscores

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tmpPath := finalPath + ".tmp"

	// Ensure the collection directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Serialize the data to JSON
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	// Write the data to a temporary file, then rename it to the final file
	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, finalPath)
}

// Read fetches a record from the specified collection and resource
func (d *Driver) Read(collection string, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection name - unable to read the data!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource name - unable to read record!")
	}

	resource = strings.ReplaceAll(resource, " ", "_") // Normalize the resource name

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	// Read the JSON file and deserialize it into the provided interface
	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

// ReadAll retrieves all records from a collection
func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection name - unable to read data!")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := os.ReadDir(dir)

	var records []string

	for _, file := range files {
		record, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(record))
	}

	return records, nil
}

// Delete removes a collection or a specific resource within it
func (d *Driver) Delete(collection string, resource string) error {
	if collection == "" {
		return fmt.Errorf("Missing collection name - unable to delete record!")
	}

	if resource != "" {
		resource = strings.ReplaceAll(resource, " ", "_")
	}

	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to find file or directory named %v\n", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

// getOrCreateMutex provides a mutex for a collection, creating it if necessary
func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

// stat checks if a file or directory exists and retrieves its information
func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}

	return fi, err
}
