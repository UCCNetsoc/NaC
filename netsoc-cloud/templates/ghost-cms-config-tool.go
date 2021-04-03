package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

const port = 80

const form = `<!doctype html>
<html lang=en>
<head>
<meta charset=UTF-8>
<meta name=viewport content='width=device-width,initial-scale=1'>
<link rel=icon type=image/png href=https://raw.githubusercontent.com/UCCNetsoc/cloud/master/ui/public/favicon.ico>
<title>GhostCMS Installer</title>
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
<script>
    window.onload = () => {
        document.getElementById('host').placeholder = window.location.host, document.getElementById('host').value = window.location.host
    }
</script>
</head>
<body>
<div id=header><img src=https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg></div>
<div id=content>
    <div>
        <h1>GhostCMS installer</h1>
    </div>
    <div id=content-body>
        <form id=form action=/config><label for=host>virtual host:</label> <input id=host name=host></form>
        <div id=info>
            <p>
                <span style="color:white">Ensure that the entered host is correctly configured as a Virtual Host in the Netsoc Cloud control panel</span>
                <br/><br/>
                Contact Netsoc SysAdmins via the UCC Netsoc Discord if you are unsure, having issues or require assistance!
            </p>
        </div>
    </div>
    <input type=submit value='Set host' id=submit form=form></div>
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
<title>GhostCMS Installer</title>
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
<script>
    window.onload = () => {
        let attempts = 0

        const intervalHdl = setInterval(async () => {
            document.getElementById("install").innerText = "Installing GhostCMS for {{.Host}}" + '.'.repeat(attempts%3)
            document.getElementById("install").className = "not-installed"
            document.getElementById("post-install-instructions").className = "not-installed-visibility"
            attempts++;

            try {
                const res = await fetch("https://{{.Host}}");

                if (res.status == 200) {
                    document.getElementById("install").innerText = "Installation successful!";
                    document.getElementById("install").className = "installed"
                    document.getElementById("post-install-instructions").className = "installed-visibility"

                    clearInterval(intervalHdl);
                }
            }
            catch (e) {

            }

            if (attempts > 90) {
                clearInterval(intervalHdl);
                document.getElementById("install").innerText = "Installation failed, contact SysAdmins";
                document.getElementById("install").className = "not-installed"
            }
        }, 1000)
    }
</script>
</head>
<body>
<div id=header><img src=https://raw.githubusercontent.com/UCCNetsoc/wiki/master/assets/logo-horizontal.svg></div>
<div id=content>
    <div>
        <h1>GhostCMS installer</h1>
    </div>
    <div id="content-body">
        <h2 id="install" class="not-installed">Installing GhostCMS for {{.Host}}</h2>
        <p id="post-install-instructions" class="not-installed-visibility" style="color: white">
            The administration page is available here<br/>
            <a href='https://{{.Host}}/ghost' target='_blank'>https://{{.Host}}/ghost</a><br/>
            Do <b>NOT</b> forget to bookmark it, it is not linked on your site
        </p>

        <p>Contact Netsoc SysAdmins on our <a href='discord.netsoc.co'>Discord</a> if you encounter any issues
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

var host string

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
	exec.Command("docker", "run", "-d", "--name", "ghost", "-v",
		"/ghost:/var/lib/ghost/content", "-e", fmt.Sprintf("url=https://%s", host),
		"-p", "80:2368", "--restart", "always", "ghost").Run()
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

func configure(done chan bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		host = r.FormValue("host")
		re := regexp.MustCompile(`https?:\/\/`)
		host = re.ReplaceAllString(host, "")
		host = strings.TrimSuffix(host, "/")

		t, _ := template.New("response").Parse(response)
		t.Execute(w, map[string]interface{}{"Host": host})
		done <- true
	}
}
