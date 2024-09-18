package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type CRMLead struct {
	EID string
	HRO	string
}
type WRAdminLead struct {
	EID string
	GEID string
	HRO string
	Perc string
	Store string
	FName string
	LName string
	Street string
	City string
	State string
	Zip	string
	Phone string
	HomePhone string
	Email string
	DoB string
	EmergencyContact string
	Status string
	Code string
	Agent string
	BatchDate string
	Notes string
	OrganizationStartDate string
	JobTitle string
	JobTitleCode string
	InternetEmail string
	InternetEmailPassword string
	ProgramID string
	SVTURL string
}

var CRMLeads []CRMLead
var WRAdminLeads []WRAdminLead
func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    file1, _, err := r.FormFile("file1")
    if err != nil {
        http.Error(w, "Invalid file1", http.StatusBadRequest)
        return
    }
    defer file1.Close()
    file2, _, err := r.FormFile("file2")
    if err != nil {
        http.Error(w, "Invalid file2", http.StatusBadRequest)
        return
    }
    defer file2.Close()

    err = readCSV(file1)
    if err != nil {
        http.Error(w, "Error reading file1 CSV", http.StatusInternalServerError)
        return
    }

	err = readCSVforWR(file2)
    if err != nil {
        http.Error(w, "Error reading file2 CSV", http.StatusInternalServerError)
        return
    }
	w.Header().Set("Contenct-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=ActiveLeadsNotInCRM.csv")
	writer := csv.NewWriter(w)
	header := []string{"Employee ID", "GEID", "HR/Owner", "%", "Store", "First Name", "Last Name", "Street", "City", "State", "Zip Code", "Phone", "Home Phone", "Email", "DoB", "Emergency Contact", "Status", "Code", "Agent", "Batch Date", "NOTES", "Organization Start Date", "Job Title", "Job Title Code", "Internal Email", "Internal Email Password", "Program ID", "CVT URL"}
	writer.Write(header)

	for _, lead := range WRAdminLeads {
        writer.Write([]string{lead.EID, lead.GEID, lead.HRO, lead.Perc, lead.Store, lead.FName, lead.LName, lead.Street, lead.City, lead.State, lead.Zip, lead.Phone, lead.HomePhone, lead.Email, lead.DoB, lead.EmergencyContact, lead.Status, lead.Code, lead.Agent, lead.BatchDate, lead.Notes, lead.OrganizationStartDate, lead.JobTitle, lead.JobTitleCode, lead.InternetEmail, lead.InternetEmailPassword, lead.ProgramID, lead.SVTURL})
    }
	writer.Flush()

	CRMLeads = nil
    WRAdminLeads = nil

    // Redirect to the homepage
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func readCSV(file io.Reader) (error) {

	csvReader := csv.NewReader(file)

	// Skip header
	_, err := csvReader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

        // Parse the record and create a Data struct
        d := CRMLead{
                EID:      strings.TrimSpace(record[1]),
                HRO:      strings.TrimSpace(record[2]),
			}
		

		CRMLeads = append(CRMLeads, d)
	}

	return nil
}


func readCSVforWR(file io.Reader) (error) { 
	csvReader := csv.NewReader(file)

	// Skip header
	_, err := csvReader.Read()
	if err != nil {
		return err
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		found := 0
		foundOwner := 0
		for _, val := range CRMLeads {
			found = 0
			if strings.TrimSpace(record[0]) == strings.TrimSpace(val.EID) {
				found = 1
			}
			if found == 1 {
				break
			}
		}

		if found == 0 {
			for _, val := range CRMLeads {
				foundOwner = 0
				if strings.TrimSpace(val.HRO) == strings.TrimSpace(record[2]) {
				foundOwner = 1
				}
				if foundOwner == 1 {
					break
				}
			}
		}

		if found == 0 && foundOwner == 1 {
			WRAdminLeads = append(WRAdminLeads, WRAdminLead{
                EID:      record[0],
                GEID:	record[1],
				HRO:	record[2],
				Perc:	record[3],
				Store:	record[4],
				FName: record[5],
				LName: record[6],
				Street:	record[7],
				City:	record[8],
				State: record[9],
				Zip: record[10],
				Phone:	record[11],
				HomePhone: 	record[12],
				Email:	record[13],
				DoB: record[14],
				EmergencyContact: record[15],
				Status:	record[16],
				Code: record[17],
				Agent: record[18],
				BatchDate: record[19],
				Notes: record[20],
				OrganizationStartDate: record[21],
				JobTitle: record[22],
				JobTitleCode: record[23],
				InternetEmail: record[24],
				InternetEmailPassword: record[25],
				ProgramID: record[26],
				SVTURL: record[27],
			}) 
		}

	}
	return nil
}
