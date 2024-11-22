package store

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"city-hall-lights/internal/model"
	"github.com/stretchr/testify/require"
)

func TestFileStore_Create(t *testing.T) {
	type args struct {
		date   time.Time
		events []model.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "unimplemented",
			args: args{
				events: []model.Event{},
			},
			wantErr: true,
			err:     fmt.Errorf("unimplemented"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStore{}
			if err := f.Create(tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStore_Delete(t *testing.T) {
	type args struct {
		event model.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "unimplemented",
			args: args{
				event: model.Event{},
			},
			wantErr: true,
			err:     fmt.Errorf("unimplemented"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStore{}
			if err := f.Delete(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStore_List(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStore{}
			got, err := f.List(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileStore_Read(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStore{}
			if _, err := f.Read(tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStore_readEventsFromFile(t *testing.T) {
	type args struct {
		date time.Time
		path string
	}
	tests := []struct {
		name      string
		args      args
		want      []model.Event
		wantErr   bool
		errString string
	}{
		{
			name: "reading multiple valid events from correctly formated file succeeds",
			args: args{
				date: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
				path: "file-test-fixtures/success-cases",
			},
			want: []model.Event{
				{
					DateString:     "Tuesday, November 5, 2024",
					StartTimeStamp: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
					Color:          "red/white/blue",
					Description:    "in recognition of Election Day 2024",
					RawEventString: "Tuesday, November 5, 2024 – red/white/blue – in recognition of Election Day 2024",
				},
				{
					DateString:     "Wednesday, November 6, 2024",
					StartTimeStamp: time.Date(2024, 11, 6, 0, 0, 0, 0, time.UTC),
					Color:          "teal",
					Description:    "in recognition of the Alzheimer Foundations’ annual “Light the World Teal” Campaign",
					RawEventString: "Wednesday, November 6, 2024 – teal – in recognition of the Alzheimer Foundations’ annual “Light the World Teal” Campaign",
				},
			},
		},
		{
			name: "reading from non-existing file fails",
			args: args{
				date: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
				path: "file-test-fixtures/failure-cases",
			},
			wantErr:   true,
			errString: "no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readEventsFromFile(tt.args.date, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("readEventsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Contains(t, err.Error(), tt.errString)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestFileStore_writeEventsToFile(t *testing.T) {
	type args struct {
		date   time.Time
		path   string
		events []model.Event
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		errString string
		cleanupFn func(string) error
	}{
		{
			name: "writing a valid event to correctly formated file succeeds",
			args: args{
				date: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
				path: "file-test-fixtures/success-cases/outputs",
				events: []model.Event{
					{
						DateString:     "Tuesday, November 5, 2024",
						StartTimeStamp: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
						Color:          "red/white/blue",
						Description:    "in recognition of Election Day 2024",
						RawEventString: "Tuesday, November 5, 2024 – red/white/blue – in recognition of Election Day 2024",
					},
				},
			},
			cleanupFn: func(filename string) error {
				return os.Remove(filename)
			},
		},
		{
			name: "writing multiple valid events to correctly formated file succeeds",
			args: args{
				date: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
				path: "file-test-fixtures/success-cases/outputs",
				events: []model.Event{
					{
						DateString:     "Tuesday, November 5, 2024",
						StartTimeStamp: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
						Color:          "red/white/blue",
						Description:    "in recognition of Election Day 2024",
						RawEventString: "Tuesday, November 5, 2024 – red/white/blue – in recognition of Election Day 2024",
					},
					{
						DateString:     "Wednesday, November 6, 2024",
						StartTimeStamp: time.Date(2024, 11, 6, 0, 0, 0, 0, time.UTC),
						Color:          "teal",
						Description:    "in recognition of the Alzheimer Foundations’ annual “Light the World Teal” Campaign",
						RawEventString: "Wednesday, November 6, 2024 – teal – in recognition of the Alzheimer Foundations’ annual “Light the World Teal” Campaign",
					},
				},
			},
			cleanupFn: func(filename string) error {
				return os.Remove(filename)
			},
		},
		{
			name: "writing to existing file fails",
			args: args{
				date:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				path:   "file-test-fixtures/failure-cases",
				events: []model.Event{},
			},
			wantErr:   true,
			errString: "file exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeEventsToFile(tt.args.date, tt.args.path, tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeEventsToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Contains(t, err.Error(), tt.errString)
				return
			}
			require.NoError(t, err)
			got, err := readEventsFromFile(tt.args.date, tt.args.path)
			require.NoError(t, err)
			require.EqualValues(t, model.Events{Events: tt.args.events}, model.Events{Events: got})

			if tt.cleanupFn != nil {
				require.NoError(t, tt.cleanupFn(generateFilename(tt.args.date, tt.args.path)))
			}
		})
	}
}

func TestFileStore_Update(t *testing.T) {
	type args struct {
		event model.Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "unimplemented",
			args: args{
				event: model.Event{},
			},
			wantErr: true,
			err:     fmt.Errorf("unimplemented"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileStore{}
			if err := f.Update(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				require.EqualError(t, err, tt.err.Error())
			}
		})
	}
}

func TestNewFileStore(t *testing.T) {
	tests := []struct {
		name string
		want FileStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateFilename(t *testing.T) {
	type args struct {
		date time.Time
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generates filename from date and path",
			args: args{
				date: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
				path: "./file-test-fixtures",
			},
			want: "./file-test-fixtures/2024-11-05.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateFilename(tt.args.date, tt.args.path); got != tt.want {
				t.Errorf("generateFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateFilename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid filename passes",
			args: args{
				filename: "2024-11-05.json",
			},
			wantErr: false,
		},
		{
			name: "valid filename with directory passes",
			args: args{
				filename: "file-test-fixtures/2024-11-05.json",
			},
			wantErr: false,
		},
		{
			name: "dot directory filename passes",
			args: args{
				filename: "./file-test-fixtures/2024-11-05.json",
			},
			wantErr: false,
		},
		{
			name: "multiple directory filename passes",
			args: args{
				filename: "path/to/my/very/special/file/2024-11-05.json",
			},
			wantErr: false,
		},
		{
			name: "incorrect suffix fails",
			args: args{
				filename: "file-test-fixtures/2024-11-05.yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFilename(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("validateFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isSameDate(t *testing.T) {
	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same date returns true",
			args: args{
				t1: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "different date returns false",
			args: args{
				t1: time.Date(2024, 11, 6, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameDate(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("isSameDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadImageFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    io.Reader
		wantErr bool
	}{
		{
			name: "valid jpeg image file returns reader",
			args: args{
				path: "images/purple.jpg",
			},
			wantErr: false,
		},
		// todo
		// {
		// 	name: "fails for png image file",
		// 	args: args{
		// 		path: "images/blue-pink-white.png",
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "fails when file not found",
		// 	args: args{
		// 		path: "images/i-dont-exist.jpg",
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "fails when image is too large",
		// 	args: args{
		// 		path: "images/blue-pink-white.png",
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadImageFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadImageFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_readImageMetadataFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.ImageMetadata
		wantErr bool
	}{
		{
			name: "valid image metadata file returns metadata",
			args: args{
				path: "images/attribution.json",
			},
			want: []model.ImageMetadata{
				{
					FileName: "blue-pink-white.jpg",
					AltText:  "San Francisco City Hall illuminated at night with vibrant pink, blue, and white lights highlighting the dome and façade. The structure stands out against a dark black sky, showcasing its detailed architectural elements and grandeur.",
					Attribution: model.Attribution{
						Creator:    "Kae Ng",
						Title:      "Brown Concrete Dome Tent",
						SourceURL:  "https://unsplash.com/photos/brown-concrete-dome-tent-sAn35VwjzEk",
						LicenseURL: "",
					},
				},
				{
					FileName: "orange.jpg",
					AltText:  "San Francisco City Hall illuminated in warm orange lights at night, highlighting the dome and front façade. A pathway lined with trees leads to the building, where a single person is seen walking towards the entrance. The surrounding sky is dark, emphasizing the vibrant glow of the building’s architectural details.",
					Attribution: model.Attribution{
						Creator:    "Gurpreet Singh",
						Title:      "The Hunt for Orange October",
						SourceURL:  "https://www.flickr.com/photos/zoxcleb/5127493349",
						LicenseURL: "https://creativecommons.org/licenses/by-sa/2.0",
					},
				},
				{
					FileName: "purple.jpg",
					AltText:  "",
					Attribution: model.Attribution{
						Creator:    "Brittany Murphy/The Chronicle",
						Title:      "San Francisco City Hall illuminated in vibrant purple lights under a dramatic cloudy night sky. The building’s grand dome and façade are highlighted by the lighting, while large glowing white sculptures resembling rabbits are displayed in the foreground, surrounded by a crowd of onlookers.",
						SourceURL:  "https://s.hdnux.com/photos/45/53/35/9877325/12/960x0.webp",
						LicenseURL: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadImageMetadataFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadImageMetadataFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadImageMetadataFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
