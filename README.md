# Pupal

Uses

* Go Google App Engine Standard for Pupal backend API -> main.go
* Datastore (NoSQL)
* React+Babel (JSX compiler) for client side -> static/js/app.jsx
* Firebase Auth

Requires

 * Install the Go App Engine SDK. Follow ["Before You Begin"](https://cloud.google.com/appengine/docs/standard/go/quickstart).
 * go get github.com/gorilla/context
 * go get github.com/gorilla/mux
 * go get github.com/urfave/negroni
 * go get github.com/wuman/firebase-server-sdk-go

Run app (Note it will set a listener for any changes to your code and auto-reload.)

* $ dev_appserver.py app.yaml
* or $ dev_appserver.py .
* To clear datastore, $ dev_appserver.py --clear_datastore.py .

### Docs:

* Go GAE Standard : https://cloud.google.com/appengine/docs/standard/go/how-to
* Datastore: https://cloud.google.com/datastore/docs/how-to
* Firebase Auth : https://firebase.google.com/docs/auth/web/manage-users
* React : https://facebook.github.io/react/docs/react-component.html


