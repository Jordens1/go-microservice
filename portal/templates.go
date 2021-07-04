package portal

import "html/template"

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error

	rootTemplate, err = template.ParseFiles(
		"../../protal/students.html",
		"../../protal/student.html")
	if err != nil {
		return err
	}

	return nil
}
