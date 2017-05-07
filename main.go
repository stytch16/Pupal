package app

import (
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	firebase "github.com/wuman/firebase-server-sdk-go"
)

func init() {

	// Initialize Firebase SDK
	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "app/firebase/serviceAccountCredentials.json",
	})

	r := mux.NewRouter()

	// Any requests to the default page will be served with the static index page in the templates folder.
	r.Handle("/", http.FileServer(http.Dir("./templates/")))

	// Server set up to serve static assets from the /static/{file} route.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Common middleware

	common := negroni.New(negroni.HandlerFunc(ValidateToken))
	admin := negroni.New(negroni.HandlerFunc(AdminToken))

	// Project handlers
	projectRouter := mux.NewRouter().PathPrefix("/projects").Subrouter()
	projectRouter.HandleFunc("/{id}", ProjectGetHandler).Methods("GET")
	projectRouter.HandleFunc("/{id}/comment", ProjectCommentHandler).Methods("POST")
	projectRouter.HandleFunc("/{id}/subscribe", ProjectSubsHandler).Methods("POST")

	// Host project handlers
	projectRouter.HandleFunc("/{domain}/host", ProjectHostPostHandler).Methods("POST")

	// Project routes middleware
	r.PathPrefix("/projects").Handler(common.With(negroni.Wrap(projectRouter)))

	// User handlers
	userRouter := mux.NewRouter().PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/registerPupalUser", UserRegisterPupalHandler).Methods("POST")
	userRouter.HandleFunc("/registerDomain", UserRegisterDomainHandler).Methods("POST")
	userRouter.HandleFunc("/delete", UserDeleteHandler).Methods("GET")
	userRouter.HandleFunc("/{id}", UserGetHandler).Methods("GET")
	userRouter.HandleFunc("/{id}/message", UserMsgHandler).Methods("POST")

	// User routes middleware
	r.PathPrefix("/users").Handler(common.With(negroni.Wrap(userRouter)))

	// Domain handlers
	domainRouter := mux.NewRouter().PathPrefix("/domain").Subrouter()
	domainRouter.HandleFunc("/list", DomainListHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}", DomainGetHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}/join", DomainJoinHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}/subscribe", DomainSubsHandler).Methods("GET")

	// Domain routes middleware
	r.PathPrefix("/domain").Handler(common.With(negroni.Wrap(domainRouter)))

	// Profile handlers
	profileRouter := mux.NewRouter().PathPrefix("/profile").Subrouter()
	profileRouter.HandleFunc("/", ProfileHandler).Methods("GET")
	profileRouter.HandleFunc("/", ProfilePostHandler).Methods("POST")
	profileRouter.HandleFunc("/projects", ProfileProjectsHandler).Methods("GET")
	profileRouter.HandleFunc("/projects/{id}", ProfileProjectsPostHandler).Methods("PUT")
	profileRouter.HandleFunc("/projects/{id}", ProfileProjectsDeleteHandler).Methods("DELETE")
	profileRouter.HandleFunc("/subscriptions", ProfileSubsHandler).Methods("GET")
	profileRouter.HandleFunc("/subscriptions/{id}", ProfileSubsDeleteHandler).Methods("DELETE")

	// Profile routes middleware
	r.PathPrefix("/profile").Handler(common.With(negroni.Wrap(profileRouter)))

	// Matchmaking handlers
	matchMakerRouter := mux.NewRouter().PathPrefix("/match").Subrouter()
	matchMakerRouter.HandleFunc("/project/{id}", MatchProjectHandler).Methods("GET")
	matchMakerRouter.HandleFunc("/user/{id}", MatchUserHandler).Methods("GET")
	matchMakerRouter.HandleFunc("/qa", MatchQAHandler).Methods("GET")
	matchMakerRouter.HandleFunc("/qa/{id}/{resp}", MatchQAPostHandler).Methods("POST")
	matchMakerRouter.HandleFunc("/qa/{id}", MatchQAPutHandler).Methods("PUT")

	// Matchmaking routes middleware
	r.PathPrefix("/match").Handler(common.With(negroni.Wrap(matchMakerRouter)))

	// Admin routes
	adminRouter := mux.NewRouter().PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/domain/add", AdminAddDomainHandler).Methods("POST")
	adminRouter.HandleFunc("/pupalusers", AdminGetUsersHandler).Methods("GET")
	r.PathPrefix("/admin").Handler(admin.With(negroni.Wrap(adminRouter)))

	http.Handle("/", r)
}

// ValidateToken validates the user's firebase token in the Authorization header field
// of the request. After authenticating, it extracts the user's uid.
// See 'negroni' middleware function for function signature.
func ValidateToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	c := appengine.NewContext(r)
	auth, _ := firebase.GetAuth()
	log.Printf("ValidateToken: %s\n", r.Header.Get("Authorization"))
	token, err := auth.VerifyIDTokenWithTransport(
		r.Header.Get("Authorization"), urlfetch.Client(c).Transport)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("ValidateToken: Failed to validate a token\n%v\n", err)
		return
	}
	uid, _ := token.UID()
	context.Set(r, "UID", uid) // context.Get(r, "UID).(string) to retrieve it
	next(w, r)
}

// AdminToken validates an admin token in the Authorization header field of the request.
func AdminToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Header.Get("Authorization") != AdminUID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid admin token."))
		return
	}
	next(w, r)
}
