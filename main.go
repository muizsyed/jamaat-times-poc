package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

var masjids []*Masjid
var masjidDay map[string]*[]Day = map[string]*[]Day{}

func main() {

	loadData()

	engine := html.New("./views", ".html")
	engine.Reload(true)

	appRevision := os.Getenv("VCS_REVISION")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {

		day := time.Now().Day()

		masjidDays := make([]MasjidDay, len(masjids))

		for i := 0; i < len(masjids); i++ {
			currentMasjid := masjids[i]
			days := *masjidDay[currentMasjid.Slug]
			masjidDays[i] = MasjidDay{
				Masjid: *currentMasjid,
				Day:    days[day],
			}
		}

		return c.Render("index", fiber.Map{
			"Title":      "Hello, World!",
			"Revision":   appRevision,
			"Masjids":    masjids,
			"MasjidDays": masjidDays,
		}, "layouts/main")
	})

	app.Listen(":8081")
}

func loadData() {
	path, _ := os.Getwd()
	fmt.Fprintln(os.Stdout, "Path is:", path, ", and month is", strings.ToLower(time.Now().Month().String()), int(time.Now().Month()))
	monthFile := fmt.Sprintf("%d.json", int(time.Now().Month()))
	fileInfos, _ := os.ReadDir(filepath.Join(path, "data"))

	for _, masjidDirectory := range fileInfos {
		fmt.Printf("Masjid: %s\n", masjidDirectory.Name())
		metadataPath := filepath.Join(path, "data", masjidDirectory.Name(), "metadata.json")
		timetablePath := filepath.Join(path, "data", masjidDirectory.Name(), monthFile)

		masjid := loadMasjidMetadata(metadataPath)
		//masjidda, _ := json.Marshal(masjid)
		//fmt.Println(string(masjidda))

		days := loadMasjidTimetable(timetablePath)

		masjids = append(masjids, masjid)
		masjidDay[masjid.Slug] = days

	}

}

func loadMasjidMetadata(path string) *Masjid {
	data, err := os.ReadFile(path)

	fmt.Fprintln(os.Stdout, "Loading metadata", path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	masjid := &Masjid{}

	json.Unmarshal(data, masjid)

	return masjid
}

func loadMasjidTimetable(path string) *[]Day {
	data, err := os.ReadFile(path)

	fmt.Fprintln(os.Stdout, "Loading timetable", path)

	if err != nil {
		fmt.Println(err)
		//return nil
	}

	days := make([]Day, daysIn(time.Now().Month(), time.Now().Year()))
	err = json.Unmarshal(data, &days)

	if err != nil {
		fmt.Println(err)
	}

	return &days
}

type Masjid struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	City     string
	Postcode string
}

type MasjidDay struct {
	Masjid Masjid
	Day    Day
}

type Prayer struct {
	Start  string `json:"start"`
	Jamaat string `json:"jamaat"`
}

type Month struct {
	Name string
	Days []Day
}

type Day struct {
	Date    string `json:"date"`
	Fajr    Prayer `json:"fajr"`
	Zuhr    Prayer `json:"zuhr"`
	Asr     Prayer `json:"asr"`
	Maghrib Prayer `json:"maghrib"`
	Esha    Prayer `json:"esha"`
}
