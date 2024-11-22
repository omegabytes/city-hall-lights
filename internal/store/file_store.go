package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
	"time"

	"city-hall-lights/internal/model"
)

var errorUnimplemented = fmt.Errorf("unimplemented")

type FileStore struct {
	path  string
	today time.Time
}

func NewFileStore() FileStore {
	return FileStore{
		path:  "internal/store/events",
		today: time.Now(),
	}
}

func (f *FileStore) Create(events []model.Event) error {
	return writeEventsToFile(f.today, f.path, events)
}

func (f *FileStore) Update(event model.Event) error {
	return errorUnimplemented
}

func (f *FileStore) Read(date time.Time) (*model.Event, error) {
	events, err := readEventsFromFile(f.today, f.path)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if isSameDate(event.StartTimeStamp, date) {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("event not found")
}

func (f *FileStore) List(date time.Time) ([]model.Event, error) {
	events, err := readEventsFromFile(date, f.path)
	if err != nil {
		return []model.Event{}, err
	}
	return events, nil
}

func (f *FileStore) Delete(_ model.Event) error {
	return errorUnimplemented
}

func readEventsFromFile(date time.Time, path string) ([]model.Event, error) {
	filename := generateFilename(date, path)
	err := validateFilename(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var events []model.Event
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return events, nil
}

func writeEventsToFile(date time.Time, path string, events []model.Event) error {
	filename := generateFilename(date, path)
	err := validateFilename(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(events); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

func generateFilename(date time.Time, path string) string {
	currentMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%s/%s.json", path, currentMonth.Format(time.DateOnly))
}

func validateFilename(filename string) error {
	if !strings.Contains(filename, ".json") {
		return fmt.Errorf("invalid format: must be json")
	}
	parts := strings.Split(filename, "/")
	before, _ := strings.CutSuffix(parts[len(parts)-1], ".json")
	_, err := time.Parse(time.DateOnly, before)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}
	return nil
}

func isSameDate(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func (f *FileStore) CheckFileExists() (bool, error) {
	filename := generateFilename(f.today, f.path)
	if _, err := os.OpenFile(filename, os.O_RDONLY, 0777); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("failed to check file: %w", err)
		}
	}
	return true, nil
}

const MAX_IMAGE_BYTES = 1000000

func LoadImageFromFile(path string) (io.Reader, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()

	imageData, err := jpeg.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, imageData)
	if err != nil {
		return nil, err
	}

	if len(buffer.Bytes()) > MAX_IMAGE_BYTES {
		return nil, fmt.Errorf("max image size exceeded: %d bytes", len(buffer.Bytes()))
	}

	return buffer, nil
}

func ReadImageMetadataFromFile(path string) ([]model.ImageMetadata, error) {
	filename := fmt.Sprintf(path)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var imageMetadata []model.ImageMetadata
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&imageMetadata); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return imageMetadata, nil
}
