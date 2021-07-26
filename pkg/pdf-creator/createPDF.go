package pdf_creator

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
)

type PatientPersonalData struct {
	Surname     string
	Name        string
	Patronymic  string
	Birthday    string
	Gender      string
	HomeAddress string
	FirstDate   string
	SecondDate  string
	Vaccine     string
	PdfName     string
}

func CreatePDF(newPatient *PatientPersonalData) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	pdf := gofpdf.New("L", "mm", "A4", pwd+"/pkg/pdf-creator/font")

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")

	pdf.AddPage()
	pdf.Image("pkg/pdf-creator/passport.png", 0, 0, 298, 0, false, "", 0, "")
	pdf.SetFont("Helvetica", "", 18)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	// add surname to pdf
	pdf.SetXY(15, 22)
	pdf.Cell(40, 10, tr(newPatient.Surname))

	// add name to pdf
	pdf.SetXY(15, 34)
	pdf.Cell(40, 10, tr(newPatient.Name))

	// add patronymic to pdf
	pdf.SetXY(15, 47)
	pdf.Cell(40, 10, tr(newPatient.Patronymic))

	// add birthday to pdf
	pdf.SetXY(55, 65)
	pdf.Cell(40, 10, tr(newPatient.Birthday))

	// add gender to pdf
	pdf.SetXY(110, 65)
	pdf.Cell(40, 10, tr(newPatient.Gender))

	// add address to pdf
	pdf.SetXY(65, 78)
	pdf.Cell(40, 10, tr(newPatient.HomeAddress))

	// add firstDate to pdf
	pdf.SetXY(70, 117)
	pdf.Cell(40, 10, tr(newPatient.FirstDate))
	pdf.SetFont("Helvetica", "", 16)
	pdf.SetXY(160, 48)
	pdf.Cell(40, 10, tr(newPatient.FirstDate))

	// add secondDate to pdf
	pdf.SetXY(160, 124)
	pdf.Cell(40, 10, tr(newPatient.SecondDate))

	// add firstVaccine to pdf
	pdf.SetXY(193, 48)
	pdf.Cell(40, 10, tr(newPatient.Vaccine))

	// add secondVaccine to pdf
	pdf.SetXY(193, 124)
	pdf.Cell(40, 10, tr(newPatient.Vaccine))

	pdf.SetFont("Helvetica", "", 13)
	pdf.SetXY(45, 212)
	pdf.Cell(20, -23, tr("Не явлется официальным документом! Разработано исключительно в образовательных целях"))

	err = pdf.OutputFileAndClose(newPatient.PdfName)
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
	} else {
		fmt.Println("Генерация прошла успешно Файл лежит в контейнере!")
	}
}
