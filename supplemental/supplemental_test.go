package supplemental_test

import (
	"github.com/noamt/go-cldr/supplemental"
	"golang.org/x/text/language"
	"testing"
	"time"
)

func TestFirstDays_ByRegion(t *testing.T) {
	type args struct {
		region language.Region
	}
	tests := []struct {
		name string
		args args
		want time.Weekday
	}{
		{"NL", args{region: regionFromLocale("nl-NL")}, time.Monday},
		{"US", args{region: regionFromLocale("en-US")}, time.Sunday},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := supplemental.FirstDay.ByRegion(tt.args.region); got != tt.want {
				t.Errorf("ByRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func regionFromLocale(locale string) language.Region {
	tag, _ := language.Parse(locale)
	region, _ := tag.Region()
	return region
}
