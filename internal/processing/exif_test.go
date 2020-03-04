package processing

import (
	"reflect"
	"testing"
	"time"
)

func Test_timeFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{"Signal files", "signal-2018-12-24-145345-1.jpg", "2018-12-24 14:53:45"},
		{"WhatsApp files", "IMG-20150724-WA0000.jpeg", "2015-07-24 00:00:00"},
		{"HDR suffix", "20130815_012230-HDR.jpg", "2013-08-15 01:22:30"},
		{"HDR suffix copy", "20130331_132717-HDR(1).jpg", "2013-03-31 13:27:17"},
		{"Many suffxes", "20130518_203736-HDR-EFFECTS-EFFECTS.jpg", "2013-05-18 20:37:36"},
		{"Resized with suffix", "Resized_20181224_145239_9203.jpeg", "2018-12-24 14:52:39"},
		{"Wordy", "Resized_Screenshot_20200220-111528_Challenge_Yo.jpg", "2020-02-20 11:15:28"},
		{"PANO prefix", "PANO_20120626_172613.jpg", "2012-06-26 17:26:13"},
		{"PANO suffix", "20130816_160034-PANO.jpg", "2013-08-16 16:00:34"},
		{"Prefix and suffix", "IMG_20131026_232959-SMILE.jpg", "2013-10-26 23:29:59"},
		{"EditN suffix", "IMG_20110807_135212_edit0.jpg", "2011-08-07 13:52:12"},
		{"Screenshot", "Screenshot_20200216-040915.jpg", "2020-02-16 04:09:15"},
		{"Raw date", "2015-02-07.jpg", "2015-02-07 00:00:00"},
		{"Raw date time", "2013-02-02 12.31.57.jpg", "2013-02-02 12:31:57"},
		{"Underscored date time", "20130811_012704.JPG", "2013-08-11 01:27:04"},
		{"Unix timestamp", "screenshot-1319130046901.png", "2011-10-20 17:00:46.901"},
		{"Numbers", "11298.jpg", ""},
		{"Text", "foo-bar.jpg", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var want *time.Time
			if len(tt.want) > 0 {
				t, _ := time.Parse("2006-01-02 15:04:05.99999", tt.want)
				want = &t
			}
			if got := timeFromFilename(tt.fileName); !reflect.DeepEqual(got, want) {
				t.Errorf("timeFromFilename() = %v, want %v", got, want)
			}
		})
	}
}
