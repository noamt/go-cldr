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
		{"Maldives", args{region: regionFromLocale("dv-MV")}, time.Friday},
		{"United Arab Emirates", args{region: regionFromLocale("ar-AE")}, time.Saturday},
		{"US", args{region: regionFromLocale("en-US")}, time.Sunday},
		{"Netherlands", args{region: regionFromLocale("nl-NL")}, time.Monday},
		{"Generic", args{region: language.Region{}}, time.Monday},
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
