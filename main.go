package main

import (
	"encoding/json"
	"fmt"
	//"golang.org/x/oauth2"
	"github.com/gosimple/oauth2"
	"io/ioutil"
	"net/http"
)

var service = oauth2.Service(
	"ac3048ae7984fc888620",
	"c2329a72c1e0e76a2ea8e8fb63662eba1f9136aa",
	"https://github.com/login/oauth/authorize",
	"https://github.com/login/oauth/access_token",
)
var repoDatas = []repoInfo{}

type userInfo struct {
	Name     string `json:"Name"`
	Login    string `json:"Login"`
	ReposURL string `json:"repos_url"`
}

type repoInfo struct {
	CloneURL string `json:"clone_url"`
	Name	string `json:"name"`
}

func loadLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got Request")

	// Set custom redirect
	service.RedirectURL = "http://127.0.0.1:80/github_callback"

	// Get authorization url.
	authUrl := service.GetAuthorizeURL("")
	fmt.Println("AuthURL is ", authUrl)
	http.Redirect(w, r, authUrl, 301)
	// Send user to authUrl and get code
}

func github_callback(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Callback")
	//fmt.Fprintf(w, "Got Callback")
	codes := r.URL.Query()["code"]
	code := codes[0]
	// Get access token.
	//fmt.Println("Code is ", code)
	token, err := service.GetAccessToken(code)
	if err != nil {
		fmt.Println("Get access token error: ", err)
		return
	}
	//fmt.Println("Token ", token)
	// Prepare resource request
	apiBaseURL := "https://api.github.com/"
	//fmt.Println("making request with token ", token)
	github := oauth2.Request(apiBaseURL, token.AccessToken)
	github.AccessTokenInHeader = true
	github.AccessTokenInHeaderScheme = "token"
	//github.AccessTokenInURL = true

	// Make the request.
	// Provide API end point (http://developer.github.com/v3/users/#get-the-authenticated-user)
	apiEndPoint := "user"
	githubUserData, err := github.Get(apiEndPoint)
	if err != nil {
		fmt.Println("Get: ", err)
		return
	}

	body, err := ioutil.ReadAll(githubUserData.Body)
	if err != nil {
		fmt.Println("Error reading response")
		return
	}
	//fmt.Fprintf(w, string(body))
	//fmt.Println("####", string(body))
	response := userInfo{}
	json.Unmarshal(body, &response)
	//fmt.Println("*** ", response.Login)
	endpoint := "users/" + response.Login + "/repos"

	repoData, err := github.Get(endpoint)
	if err != nil {
		fmt.Println("Error reading response")
		return
	}

	body, err = ioutil.ReadAll(repoData.Body)
	if err != nil {
		fmt.Println("Error reading response")
		return
	}
	//fmt.Println(string(body))

	json.Unmarshal(body, &repoDatas)

	//fmt.Println(repoDatas[0].CloneURL, repoDatas[0].Name)
	//cmd := exec.Command("git", "clone", repoDatas[0].CloneURL)
	//err = cmd.Run()
	//if err != nil {
	//	// something went wrong
	//}
	//fmt.Println("data: ", string(body))
	http.Redirect(w,r,"/display_repos", 301)
	defer githubUserData.Body.Close()
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "This is github cloning project")
}

func main() {
	fmt.Println("Starting Server")
	http.HandleFunc("/", DisplayConnectGithub)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/connect_github", loadLogin)
	http.HandleFunc("/github_callback", github_callback)
	http.HandleFunc("/display_repos", DisplayRadioButtons)
	http.HandleFunc("/selected", UserSelected)
	errs := make(chan error)
	go func() {
		errs <- http.ListenAndServe(":80", nil)
	}()
	fmt.Println("Exiting with Status ", <-errs)
}
