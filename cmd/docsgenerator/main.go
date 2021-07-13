package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
)

func main() {

	fmt.Println("Hello, World! I generate docs")

	surname := "Иванов"
	name := "Иван"
	patronymic := "Иванович"

	createPDF(surname, name, patronymic)

	fmt.Println("PDF saved successfully")
}

func createPDF(surname, name, patronymic string) *gofpdf.Fpdf {

	os.Remove("passport.pdf")

	pwd, err1 := os.Getwd()
	if err1 != nil {
		fmt.Println(err1)
	}

	pdf := gofpdf.New("L", "mm", "A4", pwd+"/font")

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")

	pdf.AddPage()
	pdf.Image("passport.png", 0, 0, 298, 0, false, "", 0, "")
	pdf.SetFont("Helvetica", "", 18)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	// add surname to pdf
	pdf.SetXY(15, 22)
	pdf.Cell(40, 10, tr(surname))

	// add name to pdf
	pdf.SetXY(15, 34)
	pdf.Cell(40, 10, tr(name))

	// add patronymic to pdf
	pdf.SetXY(15, 47)
	pdf.Cell(40, 10, tr(patronymic))

	pdf.SetFont("Helvetica", "", 13)
	pdf.SetXY(45, 212)
	pdf.Cell(20,-23, tr("Не явлется официальным документом! Разработано исключительно в образовательных целях"))
	err := pdf.OutputFileAndClose("passport.pdf")
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
	}
	return pdf
}
