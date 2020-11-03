package templates

import "html/template"

// LicenseTemplate : The licensing data HTML template
var LicenseTemplate = template.Must(template.ParseFiles("templates/license.html"))

// OfficerTemplate : The officer data HTML template
var OfficerTemplate = template.Must(template.ParseFiles("templates/officer.html"))
