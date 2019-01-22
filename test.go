package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

// Coupon structure
type Coupon struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Userid       string        `bson:"userid" json:"userid"`
	Gender       string        `bson:"gender" json:"gender"`
	Amount1      int           `bson:"amount1" json:"amount1"`
	Amount2      int           `bson:"amount2" json:"mount2"`
	Total        int           `bson:"Total" json:"Total"`
	WeekStartDay time.Time     `bson:"weekstartdate" json:"weekstartdate"`
	Coupon       CouponForWeek `bson:"coupon" json:"coupon"`
}

// CouponForWeek : Coupon structure for whole week
type CouponForWeek struct {
	Mon CouponForDay `bson:"mon" json:"mon"`
	Tue CouponForDay `bson:"tue" json:"tue"`
	Wed CouponForDay `bson:"wed" json:"wed"`
	Thr CouponForDay `bson:"thr" json:"thr"`
	Fri CouponForDay `bson:"fri" json:"fri"`
	Sat CouponForDay `bson:"sat" json:"sat"`
	Sun CouponForDay `bson:"sun" json:"sun"`
}

// CouponForDay : Coupon structure for a day
type CouponForDay struct {
	Breakfast FoodType `bson:"breakfast" json:"breakfast"`
	Lunch     FoodType `bson:"lunch" json:"lunch"`
	Dinner    FoodType `bson:"dinner" json:"dinner"`
}

// FoodType : FoodStructure for a time
type FoodType struct {
	IsSelected bool `bson:"isSelected" json:"isSelected"`
	IsVeg      bool `bson:"isVeg" json:"isVeg"`
	IsMessup   bool `bson:"ismessup" json:"isMessUp"`
}

type studentCouponInfo struct {
	UserID string
	Gender string
	Total  int
}

type totalCouponCount struct {
	Mon totalFoodForDay
	Tue totalFoodForDay
	Wed totalFoodForDay
	Thr totalFoodForDay
	Fri totalFoodForDay
	Sat totalFoodForDay
	Sun totalFoodForDay
}

type totalFoodForDay struct {
	BVeg  int
	BNVeg int
	LVeg  int
	LNVeg int
	DVeg  int
	DNVeg int
}

func main() {
	// c := [10]Coupon{}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)

	// for i := 0; i < 10; i++ {
	// 	if i%2 == 0 {
	// 		PrintCouponToPDF(c[i], pdf, "Mess-1", true)
	// 	} else {
	// 		PrintCouponToPDF(c[i], pdf, "Mess-1", false)

	// 	}
	// }

	// studentinfo := [10]studentCouponInfo{}
	// for i, si := range studentinfo {
	// 	PrintStudentCouponInfoToPdf(i, si, pdf, "Mess-1")
	// }

	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Print("err", err)
	}
}

// PrintStudentCouponInfoToPdf :
func PrintStudentCouponInfoToPdf(srl int, s studentCouponInfo, p *gofpdf.Fpdf, mess string) {
	var l, h float64
	l = 39
	h = 10
	if srl == 0 {
		headers := []string{"Sr. No", "Student Id", "Gender", "Total", "Signature"}
		for _, header := range headers {
			printCell(p, l, h, header)
		}
		p.Ln(-1)
	}

	stotal := strconv.Itoa(s.Total)
	sinfo := []string{s.UserID, s.Gender, stotal, ""}
	printCell(p, l, h, strconv.Itoa(srl))
	for _, si := range sinfo {
		printCell(p, l, h, si)
	}
	p.Ln(-1)

}

// PrintTotalCountToPDF :
func PrintTotalCountToPDF(t totalCouponCount, p *gofpdf.Fpdf, mess string) {
	var l, h float64
	l = 25
	h = 10
	p.CellFormat(0, h, mess, "0", 1, "M", false, 0, "")
	headers := []string{"Date", "B.Fast Veg", "B.Fast NVeg", "Lunch Veg", "Lunch NVeg", "Dinner Veg", "Dinner NVeg"}
	for _, header := range headers {
		printCell(p, l, h, header)
	}
	p.Ln(-1)

	dates := WholeWeekDates(time.Now())

	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	fields := []string{"BVeg", "BNVeg", "LVeg", "LNVeg", "DVeg", "DNVeg"}
	for i, day := range days {
		printCell(p, l, h, GetDateFromTime(dates[i]))
		for _, field := range fields {
			count := reflect.ValueOf(t).FieldByName(day).FieldByName(field).Int()
			// int64 to int
			printCell(p, l, h, strconv.Itoa(int(count)))
		}
		p.Ln(-1)
	}
}

// PrintCouponToPDF :
func PrintCouponToPDF(c Coupon, p *gofpdf.Fpdf, mess string, isCouponLeft bool) {

	var l1, l2, h1, shift float64
	l1 = 24
	l2 = 48
	h1 = 6
	shift = 107

	if !isCouponLeft {
		p.CellFormat(1, h1, " ", "LR", 0, "M", false, 0, "")
		p.SetXY(p.GetX(), p.GetY()-54)
	}
	// Header1
	header1 := []string{mess, c.Userid}
	for _, h := range header1 {
		printCell(p, l2, h1, h)
	}
	p.Ln(-1)
	if !isCouponLeft {
		p.SetX(shift)
	}

	// Header2
	header2 := []string{"Date", "B.Fast", "Lunch", "Dinner"}
	for _, h := range header2 {
		printCell(p, l1, h1, h)
	}
	p.Ln(-1)
	if !isCouponLeft {
		p.SetX(shift)
	}

	// Body
	dates := WholeWeekDates(time.Now())
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	times := []string{"Breakfast", "Lunch", "Dinner"}

	for i, day := range days {
		printCell(p, l1, h1, GetDateFromTime(dates[i]))
		for _, time := range times {

			booked := reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsSelected").Bool()
			isVeg := reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsVeg").Bool()

			if booked {
				if isVeg {
					printCell(p, l1, h1, "VEG")
				} else {
					printCell(p, l1, h1, "NVEG")
				}
			} else {
				printCell(p, l1, h1, "***")
			}

		}
		p.Ln(-1)
		if !isCouponLeft {
			p.SetX(shift)
		}
	}

	// Footer
	var footer [2]string
	amt := "Amount to be paid. RS "
	if mess == "Mess-1" {
		footer[0] = amt + strconv.Itoa(c.Amount1)
	} else {
		footer[0] = amt + strconv.Itoa(c.Amount2)
	}
	footer[1] = c.Gender
	for _, f := range footer {
		printCell(p, l2, h1, f)
	}
	if !isCouponLeft {
		p.Ln(-1)
		p.Ln(-1)
	}

}

func printCell(p *gofpdf.Fpdf, l float64, h float64, txt string) {
	fullBorder := "1"
	middleAlign := "M"
	emptyLink := ""
	contSameLine := 0
	noColor := false
	noLink := 0
	p.CellFormat(l, h, txt, fullBorder, contSameLine, middleAlign, noColor, noLink, emptyLink)
}

// GetDateFromTime :
func GetDateFromTime(t time.Time) string {
	return t.Format(time.RFC3339)[:10]
}

// FirstDayofWeek :
func FirstDayofWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

// WholeWeekDates :
func WholeWeekDates(t time.Time) []time.Time {
	var array []time.Time
	t = FirstDayofWeek(t)
	// // Add 7 days to get next week monday date
	// t = t.AddDate(0, 0, 7)
	array = append(array, t)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
		array = append(array, t)
	}
	return array
}
