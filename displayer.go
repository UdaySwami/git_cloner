package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
)

type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type PageVariables struct {
	PageTitle        string
	PageRadioButtons []RadioButton
	Answer           string
}

func DisplayConnectGithub(w http.ResponseWriter, r *http.Request) {
	Title := "Select Below button to connect github"
	MyRadioButtons := []RadioButton{
		RadioButton{"ConnectGithub", "connect_github", false, false, "Connect Github"},
	}
	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}
	t, err := template.ParseFiles("home.html") //parse the html file homepage.html
	if err != nil {                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func DisplayRepos(w http.ResponseWriter, r *http.Request) {
	// Display some radio buttons to the user

	Title := "Which Repo to Clone?"
	MyRadioButtons := []RadioButton{}
	for _, v := range repoDatas {
		fmt.Printf(v.Name)
		rb := RadioButton{
			Name:       "RepoName",
			Value:      v.Name,
			IsDisabled: false,
			IsChecked:  false,
			Text:       v.Name,
		}
		MyRadioButtons = append(MyRadioButtons, rb)
	}

	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}

	t, err := template.ParseFiles("select.html") //parse the html file homepage.html
	if err != nil {                              // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func UserSelected(w http.ResponseWriter, r *http.Request) {
	log.Println("UserSelected Repo")
	r.ParseForm()

	selectedRepoName := r.Form.Get("RepoName")

	Title := "Your selected Repo Name"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    selectedRepoName,
	}

	var cloneURL string
	for _, v := range repoDatas {
		if selectedRepoName == v.Name {
			cloneURL = v.CloneURL
		}
	}

	cmd := exec.Command("git", "clone", cloneURL)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error Cloning Repo:", err)
	}
	// generate page by passing page variables into template
	t, err := template.ParseFiles("select.html") //parse the html file homepage.html
	if err != nil {                              // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
