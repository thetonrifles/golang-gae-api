# Go Endpoints for simulating API usage
This application allows to validate requests with a specific API KEY, that is returned when asking it with a specific endpoint.

Follow these steps for configuring your local environment:

1. Install [Go](https://golang.org/dl/)
2. Install [Google App Engine SDK](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)
3. Clone repository
4. Install libraries

    ```sh
    $ goapp get google.golang.org/appengine
    $ goapp get github.com/gorilla/mux
    ```

5. Run Application

    ```sh
    goapp serve <PATH-TO-ROOT-FOLDER>
    ```
