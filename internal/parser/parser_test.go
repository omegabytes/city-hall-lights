package parser

import (
	"reflect"
	"testing"
	"time"

	"city-hall-lights/internal/model"
	"github.com/stretchr/testify/require"
)

func TestParseEvent(t *testing.T) {
	type args struct {
		rawEventString string
	}
	tests := []struct {
		name string
		args args
		want model.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseEvent(tt.args.rawEventString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToTimestamps(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)
	type args struct {
		dateString string
	}
	tests := []struct {
		name    string
		args    args
		want    []time.Time
		wantErr bool
	}{
		{
			name: "parses single date with year",
			args: args{
				dateString: "Saturday, November 2, 2024",
			},
			want: []time.Time{time.Date(2024, 11, 2, 0, 0, 0, 0, loc)},
		},
		{
			name: "parses single date without year",
			args: args{
				dateString: "Saturday, November 2",
			},
			want: []time.Time{time.Date(2024, 11, 2, 0, 0, 0, 0, loc)},
		},
		{
			name: "parses date range using 'and' phrasing",
			args: args{
				dateString: "Sunday, November 3 and Monday, November 4, 2024",
			},
			want: []time.Time{
				time.Date(2024, 11, 3, 0, 0, 0, 0, loc),
				time.Date(2024, 11, 4, 0, 0, 0, 0, loc),
			},
		},
		{
			name: "parses date range using 'through' phrasing",
			args: args{
				dateString: "Sunday, November 3 through Monday, November 4, 2024",
			},
			want: []time.Time{
				time.Date(2024, 11, 3, 0, 0, 0, 0, loc),
				time.Date(2024, 11, 4, 0, 0, 0, 0, loc),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToTimestamps(tt.args.dateString)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToTimestamps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.want, got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("convertToTimestamps() got = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_parseSingleDate(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)
	type args struct {
		dateStr           string
		layout            string
		layoutWithoutYear string
		year              int
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "parses single date with year",
			args: args{
				dateStr:           "Saturday, November 2, 2024",
				layout:            "Monday, January 2, 2006",
				layoutWithoutYear: "Monday, January 2",
				year:              2024,
			},
			want: time.Date(2024, 11, 2, 0, 0, 0, 0, loc),
		},
		{
			name: "parses single date without year",
			args: args{
				dateStr:           "Saturday, November 2",
				layout:            "Monday, January 2, 2006",
				layoutWithoutYear: "Monday, January 2",
				year:              2024,
			},
			want: time.Date(2024, 11, 2, 0, 0, 0, 0, loc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSingleDate(tt.args.dateStr, tt.args.layout, tt.args.layoutWithoutYear, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSingleDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSingleDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformDescription(t *testing.T) {
	type args struct {
		description string
		color       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "transforms description with one color and internal quotes",
			args: args{
				description: `in recognition of the Alzheimer Foundations’ annual "Light the World Teal" Campaign`,
				color:       "teal",
			},
			want: `Tonight City Hall will be teal in recognition of the Alzheimer Foundations’ annual "Light the World Teal" Campaign`,
		},
		{
			name: "transforms description with one color and internal quotes in format `“`",
			args: args{
				description: `in recognition of the Alzheimer Foundations’ annual “Light the World Teal“ Campaign`,
				color:       "teal",
			},
			want: `Tonight City Hall will be teal in recognition of the Alzheimer Foundations’ annual "Light the World Teal" Campaign`,
		},
		{
			name: "transforms description with one color and no joiner",
			args: args{
				description: `SFDPH "Living Proof" campaign`,
				color:       "Blue",
			},
			want: `Tonight City Hall will be blue in recognition of SFDPH "Living Proof" campaign`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformDescription(tt.args.description, tt.args.color); got != tt.want {
				t.Errorf("transformDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformColors(t *testing.T) {
	type args struct {
		colors string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   bool
		errString string
	}{
		{
			name: "empty string returns error",
			args: args{
				colors: "",
			},
			wantErr:   true,
			errString: "invalid input: color must not be empty",
		},
		{
			name: "string with only spaces returns error",
			args: args{
				colors: "     ",
			},
			wantErr:   true,
			errString: "invalid input: color must not be empty",
		},
		{
			name: "one color in the format 'red' is a no-op",
			args: args{
				colors: "teal",
			},
			want: "teal",
		},
		{
			name: "one color in the format 'shades of amber' is a no-op",
			args: args{
				colors: "shades of amber",
			},
			want: "shades of amber",
		},
		{
			name: "transforms two colors in the format 'red/white'",
			args: args{
				colors: "orange/gold",
			},
			want: "orange and gold",
		},
		{
			name: "transforms three colors in the format 'red/white/blue'",
			args: args{
				colors: "red/white/blue",
			},
			want: "red, white, and blue",
		},
		{
			name: "transforms six colors in the format 'red/white/blue/cyan/magenta/yellow'",
			args: args{
				colors: "red/white/blue/cyan/magenta/yellow",
			},
			want: "red, white, blue, cyan, magenta, and yellow",
		},
		{
			name: "lowercases one color in the format 'Red'",
			args: args{
				colors: "Blue",
			},
			want: "blue",
		},
		{
			name: "lowercases one color in the format 'Shades of Amber'",
			args: args{
				colors: "Shades of Amber",
			},
			want: "shades of amber",
		},
		{
			name: "lowercases multiple colors in the format 'Red/White/Blue'",
			args: args{
				colors: "Red/White/Blue",
			},
			want: "red, white, and blue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformColors(tt.args.colors)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeEventsToFile() error = %v, wantErr %v", err, tt.wantErr)
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
