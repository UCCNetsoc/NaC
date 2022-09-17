package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"
)

const port = 80

const form = `<!doctype html>
<html lang=en>
<head>
<meta charset=UTF-8>
<meta name=viewport content='width=device-width,initial-scale=1'>
<link rel=icon type=image/png href=https://raw.githubusercontent.com/UCCNetsoc/cloud/master/ui/public/favicon.ico>
<title>FoundryVTT Installer</title>
<style>
    @import url('https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;1,100&display=swap');

    body {
        padding: 0;
        margin: 0;
        background-color: #212121;
        font-size: 14px;
        font-weight: 400;
        font-family: Roboto, sans-serif!important;
        color: #fff
    }

    #header {
        background-color: #2196f3;
        text-align: center;
        height: 64px;
    }

    #header img {
        height: 36px;
        margin: 0 auto;
        padding: 16px
    }

    #content {
        text-align: center;
        max-width: 440px;
        margin: 0 auto
    }

    #content-body {
        padding-top: 16px;
        border-top: 1px solid #333;
        color: gray;
        border-bottom: 1px solid #333
    }

    p {
        max-width: 380px;
        margin: 0 auto;
        padding: 16px;
    }

    h1 {
        font-size: 24px;
        font-weight: 200;
        padding: 8px 0
    }

    label {
        position: relative;
        opacity: 1;
        margin: 0 auto;
        font-size: 1.2em;
        padding: 8px
    }

    input {
        margin: 8px auto;
        height: 16px;
        background-color: #00000040;
        border: none;
        padding: 8px;
        color: #fff;
        width: calc(100% - 16px);
        font-size: 16px
    }

    form {
        width: 100%;
        margin: 0 auto;
        text-align: left
    }

    a {
        color: #fff
    }

    #submit {
        width: 100%;
        background-color: #2196f3;
        font-weight: 500;
        margin-top: 20px;
        padding: 20px 0;
        text-transform: uppercase;
        font-family: Roboto, sans-serif!important;
        border-radius: 5px;
        line-height: 0;
        letter-spacing: .09em
    }

</style>
</head>
<body>
<div id=header><img src=https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg></div>
<div id=content>
    <div>
        <h1>FoundryVTT Installer</h1>
    </div>
    <div id=content-body>
        <form id=form action=/config>
            <label for=license>FoundryVTT License:</label> 
            <p>To view a purchased license</p>
            <ul>
                <li>Go to <a href="https://foundryvtt.com/" target="_blank">foundryvtt.com</a></li>
                <li>Login and go to your profile</li>
                <li>Select <strong>Purchased Licenses</strong> from the left column</li>
            </ul>
            <input id=license name=license>
            <label for=password>Admin Password:</label> 
            <input id=password name=password type="password">
        </form>
    </div>
    <input type=submit value='Launch Foundry VTT Instance' id=submit form=form></div>
    <div style="
      margin: 0 auto;
      text-align: center;
      padding-top: 16px;
      font-size: 11px;
      color: gray;
    ">
      <a target="_blank" href="https://discord.netsoc.co">Discord</a>
      <span> • </span>
      <a target="_blank" href="https://wiki.netsoc.co/en/services/terms-of-service">Terms of Service</a>
    </div>
</body>
`

const response = `<!doctype html>
<html lang=en>
<head>
<meta charset=UTF-8>
<meta name=viewport content='width=device-width,initial-scale=1'>
<link rel=icon type=image/png href=https://raw.githubusercontent.com/UCCNetsoc/cloud/master/ui/public/favicon.ico>
<title>FoundryVTT Installer</title>
<style>
    @import url('https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100;0,300;0,400;1,100&display=swap');

    body {
        padding: 0;
        margin: 0;
        background-color: #212121;
        font-size: 14px;
        font-weight: 400;
        font-family: Roboto, sans-serif !important;
        color: #fff
    }

    #header {
        background-color: #2196f3;
        text-align: center;
        height: 64px;
    }

    #header img {
        height: 36px;
        margin: 0 auto;
        padding: 16px
    }

    #content {
        text-align: center;
        max-width: 440px;
        margin: 0 auto
    }

    #content-body {
        padding-top: 16px;
        border-top: 1px solid #333;
        color: gray;
        border-bottom: 1px solid #333
    }

    p {
        max-width: 380px;
        margin: 0 auto;
        padding: 16px;
    }

    h1 {
        font-size: 24px;
        font-weight: 200;
        padding: 8px 0
    }

    h2 {
        font-family: Roboto, sans-serif !important;
    }

    label {
        position: relative;
        opacity: 1;
        margin: 0 auto;
        font-size: 1.2em;
        padding: 8px
    }

    input {
        margin: 8px auto;
        height: 16px;
        background-color: #00000040;
        border: none;
        padding: 8px;
        color: #fff;
        width: calc(100% - 16px);
        font-size: 16px
    }

    form {
        width: 100%;
        margin: 0 auto;
        text-align: left
    }

    a {
        color: #fff
    }

    #submit {
        width: 100%;
        background-color: #2196f3;
        font-weight: 500;
        margin-top: 20px;
        padding: 20px 0;
        text-transform: uppercase;
        font-family: Roboto, sans-serif!important;
        border-radius: 5px;
        line-height: 0;
        letter-spacing: .09em
    }

    .not-installed {
        color: red !important;
    }

    .installed {
        color: green !important;
    }

    .not-installed-visibility {
        display: none;
    }

    .installed-visibility {
        display: inherit;
    }

</style>
</head>
<body>
<div id=header><img src=https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg></div>
<div id=content>
    <div>
        <h1>FoundryVTT Installer</h1>
    </div>
    <div id="content-body">
        <p>Your FoundryVTT install should come up in a minute <a href="/" target="_blank">here</a>
        <p>Contact Netsoc SysAdmins on our <a href='https://discord.netsoc.co'>Discord</a> if you encounter any issues
    </div>
    <div style="
      margin: 0 auto;
      text-align: center;
      padding-top: 16px;
      font-size: 11px;
      color: gray;
    ">
      <a target="_blank" href="https://discord.netsoc.co">Discord</a>
      <span> • </span>
      <a target="_blank" href="https://wiki.netsoc.co/en/services/terms-of-service">Terms of Service</a>
    </div>
</body>`

var license string
var password string

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startHTTPServer(httpServerExitDone)

	done := make(chan bool)
	configureFunc := configure(done)

	http.HandleFunc("/", showForm)
	http.HandleFunc("/config", configureFunc)

	<-done

	// Wait for server to shutdown and wait additional 10s for good measure
	_ = srv.Shutdown(context.TODO())
	httpServerExitDone.Wait()
	time.Sleep(10 * time.Second)

	// Commands to execute after web server shutdown
	exec.Command("systemctl", "disable", "config_tool.service").Run()
	docker_cmd := exec.Command("docker-compose", "up", "--detach")
    docker_cmd.Dir = "/root/foundry"
    docker_cmd.Run()
}

func startHTTPServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port)}

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

/// Web Handler Funcs

func showForm(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("form").Parse(form)
	t.Execute(w, nil)
}

const secrets_json_template = `{
    "foundry_admin_key": "{{.Password}}",
    "foundry_license_key":  "{{.License}}"
}`

func configure(done chan bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		license = r.FormValue("license")
        password = r.FormValue("password")

        file, _ := os.Create("/root/foundry/secrets.json")
        defer file.Close()

        tmpl, _ := template.New("secrets").Parse(secrets_json_template)
        tmpl.Execute(file , map[string]interface{}{"License": license, "Password": password})

        w.Write([]byte(response))
		done <- true
	}
}
