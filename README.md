```
                      _   _                           _                                _      
                     (_) (_)                         | |                              | |     
  __ _   ___    ___   _   _   ______    __ _   _ __  | |_   ______  __      __   ___  | |__   
 / _` | / __|  / __| | | | | |______|  / _` | | '__| | __| |______| \ \ /\ / /  / _ \ | '_ \  
| (_| | \__ \ | (__  | | | |          | (_| | | |    \ |_            \ V  V /  |  __/ | |_) | 
 \__,_| |___/  \___| |_| |_|           \__,_| |_|     \__|            \_/\_/    \___| |_.__/  
                                                                                              
                                                                                              
By Allen Lee (ylee) & Othmane Afiali (oafiali)
```

# ascii-art-web
A web version of the ascii-art with simple GUI for easy usability.
The website is hosted on the local machine using port 8080.
Only `localhost:8080/` and `localhost:8080/ascii-art`are valid URL.

## Usage
Prerequisites: Golang, git and a browser application.

Download the project using `git clone https://01.gritlab.ax/git/ylee/ascii-art-web.git`. Navigate to the project folder in the terminal and start the server by entering `go run main.go`. Key in a valid URL, `localhost:8080/` or `localhost:8080/ascii-art`. And you will see the home page (refer to image below). You will see:
- a top bar with a link to the `homepage` and `about`
- an ascii art banner
- a short description on how to use the site
- a text field for user to input text
- radio buttons to select banner style
- a button to generate the ascii art. 
- result output

![alt text](assets/images/image.png)

Simply enter the your text, select the banner style and click the button to generate the ASCII art banner.
Leading and trailing newlines will be remove (user can enter the newlines themselves).
Only ASCII characters (from 32 to 127) and newline is considered acceptable input.

Non-printable ASCII characters will result in an error.


## Implementation
Below is a flowchart that maps out how the program works.

```mermaid
flowchart TB
    subgraph client [CLIENT Browser Application]
        pressBtn((Press 'Generate ASCII Art'))
        invalidURL["404 not found"]
        inValidRq["405 method not allowed"]
        rFormErr["400 bad request"]
        missing["500 internal server error"]
        index["display index.html
        200 OK"]
    end

    subgraph server [SERVER]
    start{"go run main.go"}
        subgraph main ["main()"]
            checkReq("checkRequired()
            If missing, log.Fatal")
            static("define static 
            FileServer route")
            register("http.HandleFunc(/, homeHandler)")
            listener("http.ListenAndServe(8080)")
        end
        homeFn("homeHandler()")
        template("getTemplate()")
        subgraph handlePost ["handlePost()"]
            rForm("getFormInputs()")
            clean("clean inputs for genArt()")
            genArt("GenArt()")
        end

        checkReq--> static -->register--> listener--> homeFn;
        homeFn-->
        |Valid GET request|template-->
        index--User input-->
        pressBtn-->
        |Post request to
        localHost:8080/ascii-art| listener;
        homeFn-->
        |POST request|handlePost;
        rForm-->clean-->genArt-->|art/ error generated|template;
    end
    
    start --> main;
    homeFn--Invalid URL
    redirect to /error-->invalidURL;
    homeFn--Invalid Request-->inValidRq;
    rForm--Uable to grab data-->rFormErr;
    genArt--Missing files-->missing;
```

1. In main, checkRequired() function checks if all necessary files (such as templates, static resources, etc.) exist in the project folder.
2. `/static` is defined as static server route to server static files effectively to clients.
3. The homeHandler() function is registered to the pattern `/`, handling incoming requests to the root (in this situation, it is effectively handling all incoming requests other than `/static`).
4. http.ListenAndServe("8080") starts listening on local port 8080 and use the [DefaultServeMux] to handle requests.
5. The client (browser application) can send HTTP requests to the server.  If it is valid, the server will reply with a HTTP response (HTML page) and the 200 status code. Else, an error page with 404 status code.
5. If it is a valid `GET` request, getTemplate() is called. getTemplate() uses the html.Template.ParseFiles() and Execute() to generate the response file (HTML).
6. The response (home page) is display by the client and the user can input the desired parameters. When the `Generate ASCII Art` button is hit, a POST request is sent to the server with the user input.
7. If it is a `POST` request, handlePost() is called and in it:
    - we grab the user inputs using getFormInputs() (http.Request.ParseForm())
    - we generate the output with GenArt()
    - and getTemplate() to generate the response with the output

Note:
- net/http help us starts the server, listen and handles the requests.
- html/template help us to generate the desired html response page without spending too much time on the html files. We have to include some of html.Template snytax in the html files and parse it to generate the desired result.
- When encountering an error with a HTTP request or response, errorHandler() is used to generate a custom error page. If the `error.html` is not found, a simple error message is display instead.
