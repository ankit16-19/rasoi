package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Coupon structure
type Coupon struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Userid       string        `bson:"userid" json:"userid"`
	Gender       string        `bson:"gender" json:"gender"`
	UserName     string        `bson:"username" json:"username"`
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
	Name   string
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

// DAO  :
type DAO struct{}

// Db :
var Db *mgo.Database

// Connect :
func (c *DAO) Connect() {
	session, err := mgo.Dial("172.16.1.213")
	if err != nil {
		log.Fatal(err)
	}
	Db = session.DB("rasoi")
}

func main() {
	var d = DAO{}
	d.Connect()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)
	var filename string
	// NOTE: studentcouponss
	{

		orQuery := []bson.M{}
		days := []string{"mon", "tue", "wed", "thr", "fri", "sat", "sun"}
		times := []string{"breakfast", "lunch", "dinner"}
		isMessUP := false
		dd, _ := time.Parse("2006-01-02", "2019-03-31")
		for _, day := range days {
			for _, t := range times {
				andQuery := []bson.M{}
				andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "ismessup": isMessUP})
				andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isSelected": true})

				orQuery = append(orQuery, bson.M{"$and": andQuery})
			}
		}
		var coupons []Coupon

		err2 := Db.C("coupons").Find(bson.M{"$query": bson.M{"weekstartdate": bson.M{"$gte": dd}, "$or": orQuery}, "$orderby": bson.M{"userid": 1}}).All(&coupons)
		if err2 != nil {
			fmt.Print("Error in getting coupons ", err2)
		}
		var messname string
		if isMessUP {
			messname = "Mess-Up"
			filename = "studentCouponsMessUp.pdf"
		} else {
			messname = "Mess-Down"
			filename = "studentCouponsMessDown.pdf"
		}
		for i := 0; i < len(coupons); i++ {
			if i%2 == 0 {
				PrintCouponToPDF(coupons[i], pdf, messname, true, isMessUP)
			} else {
				PrintCouponToPDF(coupons[i], pdf, messname, false, isMessUP)
			}
		}

	}

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Print("err", err)
	}

	// NOTE: StudentCouponInfo
	{
		StudentCouponInfoToPdfFunc(true)
		StudentCouponInfoToPdfFunc(false)
	}

	// NOTE: Day wise
	{
		PrintTotalCountToPDFFUNC(true)
		PrintTotalCountToPDFFUNC(false)
	}

}

// PrintTotalCountToPDFFUNC :
func PrintTotalCountToPDFFUNC(mess bool) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)
	var filename string

	var tcc totalCouponCount
	daysCapital := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	foodForDay := []string{"BVeg", "BNVeg", "LVeg", "LNVeg", "DVeg", "DNVeg"}
	days := []string{"mon", "tue", "wed", "thr", "fri", "sat", "sun"}
	times := []string{"breakfast", "lunch", "dinner"}
	ismessup := mess
	dd, _ := time.Parse("2006-01-02", "2019-03-31")
	var messname string
	if ismessup {
		messname = "Mess-Up"
		filename = "studentCouponTotalMessUp.pdf"
	} else {
		messname = "Mess-Down"
		filename = "studentCouponTotalMessDown.pdf"
	}
	for i, day := range days {
		for j, t := range times {
			andQuery := []bson.M{}
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isVeg": true})
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isSelected": true})
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "ismessup": ismessup})

			countVeg, err := Db.C("coupons").Find(bson.M{"weekstartdate": bson.M{"$gte": dd}, "$and": andQuery}).Count()
			if err != nil {
				fmt.Print("Got error in getting coupon count")
			}
			reflect.ValueOf(&tcc).Elem().FieldByName(daysCapital[i]).FieldByName(foodForDay[j*2]).SetInt(int64(countVeg))
			andQuery = andQuery[:0]
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isVeg": false})
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isSelected": true})
			andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "ismessup": ismessup})

			countNVeg, err2 := Db.C("coupons").Find(bson.M{"weekstartdate": bson.M{"$gte": dd}, "$and": andQuery}).Count()
			if err2 != nil {
				fmt.Print("Got error in getting coupon count")
			}
			reflect.ValueOf(&tcc).Elem().FieldByName(daysCapital[i]).FieldByName(foodForDay[j*2+1]).SetInt(int64(countNVeg))

		}
	}
	PrintTotalCountToPDF(&tcc, pdf, messname)

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Print("err", err)
	}
}

// StudentCouponInfoToPdfFunc :
func StudentCouponInfoToPdfFunc(mess bool) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)
	var filename string

	{
		var coupons []Coupon
		orQuery := []bson.M{}
		days := []string{"mon", "tue", "wed", "thr", "fri", "sat", "sun"}
		times := []string{"breakfast", "lunch", "dinner"}
		isMessUP := mess
		dd, _ := time.Parse("2006-01-02", "2019-03-31")
		for _, day := range days {
			for _, t := range times {
				andQuery := []bson.M{}
				andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "ismessup": isMessUP})
				andQuery = append(andQuery, bson.M{"coupon" + "." + day + "." + t + "." + "isSelected": true})

				orQuery = append(orQuery, bson.M{"$and": andQuery})
			}
		}

		err2 := Db.C("coupons").Find(bson.M{"$query": bson.M{"weekstartdate": bson.M{"$gte": dd}, "$or": orQuery}, "$orderby": bson.M{"gender": 1, "userid": 1}}).All(&coupons)
		if err2 != nil {
			fmt.Print("Error in getting coupons ", err2)
		}
		var messname string
		if isMessUP {
			messname = "Mess-Up"
			filename = "studentCouponInfoMessUp.pdf"
		} else {
			messname = "Mess-Down"
			filename = "studentCouponInfoMessDown.pdf"
		}
		for i := 0; i < len(coupons); i++ {
			var sci studentCouponInfo
			sci.UserID = coupons[i].Userid
			sci.Gender = coupons[i].Gender
			sci.Name = coupons[i].UserName
			if isMessUP {
				sci.Total = coupons[i].Amount1
			} else {
				sci.Total = coupons[i].Amount2
			}
			PrintStudentCouponInfoToPdf(i, sci, pdf, messname)
		}
	}
	err := pdf.OutputFileAndClose(filename)
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
		printCell(p, l, h, mess)
		p.Ln(-1)
		headers := []string{"Sr. No", "Student Id", "Name", "Gender", "Total"}
		for _, header := range headers {
			printCell(p, l, h, header)
		}
		p.Ln(-1)
	}

	stotal := strconv.Itoa(s.Total)
	sinfo := []string{s.UserID, s.Name, s.Gender, stotal, ""}
	printCell(p, l, h, strconv.Itoa(srl+1))
	for _, si := range sinfo {
		printCell(p, l, h, si)
	}
	p.Ln(-1)

}

// PrintTotalCountToPDF :
func PrintTotalCountToPDF(t *totalCouponCount, p *gofpdf.Fpdf, mess string) {
	var l, h float64
	l = 25
	h = 10
	p.CellFormat(0, h, mess, "0", 1, "M", false, 0, "")
	headers := []string{"Date", "B.Fast Veg", "B.Fast NVeg", "Lunch Veg", "Lunch NVeg", "Dinner Veg", "Dinner NVeg"}
	for _, header := range headers {
		printCell(p, l, h, header)
	}
	p.Ln(-1)

	dates := WholeWeekDates(time.Now().AddDate(0, 0, 7))

	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	fields := []string{"BVeg", "BNVeg", "LVeg", "LNVeg", "DVeg", "DNVeg"}
	for i, day := range days {
		printCell(p, l, h, GetDateFromTime(dates[i]))
		for _, field := range fields {
			count := reflect.ValueOf(t).Elem().FieldByName(day).FieldByName(field).Int()
			// int64 to int
			printCell(p, l, h, strconv.Itoa(int(count)))
		}
		p.Ln(-1)
	}
}

// PrintCouponToPDF :
func PrintCouponToPDF(c Coupon, p *gofpdf.Fpdf, mess string, isCouponLeft bool, ismessup bool) {

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

	header1 := []string{mess, c.UserName, c.Userid}

	for _, h := range header1 {
		printCell(p, 32, h1, h)
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
	dates := WholeWeekDates(time.Now().AddDate(0, 0, 7))
	days := []string{"Mon", "Tue", "Wed", "Thr", "Fri", "Sat", "Sun"}
	times := []string{"Breakfast", "Lunch", "Dinner"}
	var ismessup2 bool
	for i, day := range days {
		printCell(p, l1, h1, GetDateFromTime(dates[i]))
		for _, time := range times {
			if ismessup {
				ismessup2 = reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsMessup").Bool()
			} else {
				ismessup2 = !reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsMessup").Bool()
			}

			booked := reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsSelected").Bool()
			isVeg := reflect.ValueOf(c.Coupon).FieldByName(day).FieldByName(time).FieldByName("IsVeg").Bool()

			if booked && ismessup2 {
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
	if ismessup {
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
