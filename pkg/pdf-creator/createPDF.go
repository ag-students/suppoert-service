package pdf_creator

import (
	"fmt"
	"os"
	"github.com/jung-kurt/gofpdf"
)

func CreatePDF(surname, name, patronymic, pdf_name string)  {
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
	pdf.Cell(40, 10, tr(surname))

	// add name to pdf
	pdf.SetXY(15, 34)
	pdf.Cell(40, 10, tr(name))

	// add patronymic to pdf
	pdf.SetXY(15, 47)
	pdf.Cell(40, 10, tr(patronymic))

	pdf.SetFont("Helvetica", "", 13)
	pdf.SetXY(45, 212)
	pdf.Cell(20, -23, tr("Не явлется официальным документом! Разработано исключительно в образовательных целях"))
	
	err = pdf.OutputFileAndClose(pdf_name)
	if err != nil {
		fmt.Println("⚠️  Could not save PDF:", err)
	} else {
		fmt.Println("Генерация прошла успешно Файл лежит в контейнере!")	
		// docker cp <container_name>:/support-service/passport.pdf ./
		// чтобы достать файл из контейнера
	}
}
