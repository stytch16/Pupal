package app

import (
	"log"
	"net/http"

	ctx "golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"

	gorCtx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	firebase "github.com/wuman/firebase-server-sdk-go"
)

func init() {

	// Initialize Firebase SDK
	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "app/firebase/serviceAccountCredentials.json",
	})

	// Create gorilla router
	r := mux.NewRouter()

	// Any requests to the default page will be served with the static index page in the templates folder.
	r.Handle("/", http.FileServer(http.Dir("./templates/")))

	// Server set up to serve static assets from the /static/{file} route.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Common middleware
	common := negroni.New(negroni.HandlerFunc(ValidateToken), negroni.HandlerFunc(GetUserFromCache))
	admin := negroni.New(negroni.HandlerFunc(AdminToken))

	// Project handlers
	projectRouter := mux.NewRouter().PathPrefix("/projects").Subrouter()
	projectRouter.HandleFunc("/{id}", ProjectGetHandler).Methods("GET")
	projectRouter.HandleFunc("/{id}/comment", ProjectCommentHandler).Methods("POST")
	projectRouter.HandleFunc("/{id}/subscribe", ProjectSubsHandler).Methods("GET")

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
	domainRouter.HandleFunc("/userlist", DomainUserListHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}", DomainGetHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}/member", DomainGetMemberHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}/projectlist", DomainProjectListHandler).Methods("GET")
	domainRouter.HandleFunc("/{id}/join", DomainJoinHandler).Methods("POST")

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
	gorCtx.Set(r, "UID", uid) // context.Get(r, "UID").(string) to retrieve it
	next(w, r)
}

// GetUserFromCache gets the pupal user from memcache. If pupal user does not exist,
// then it gets the pupal user from datastore. If the pupal user does not exist there,
// then the user must be registered to datastore by a POST /user/registerPupalUser
func GetUserFromCache(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	c := appengine.NewContext(r)
	uid := gorCtx.Get(r, "UID").(string)

	var pu PupalUser
	if _, err := memcache.Gob.Get(c, uid, &pu); err == memcache.ErrCacheMiss {
		// Grab from datastore
		pu.Key = datastore.NewKey(c, "PupalUser", uid, 0, datastore.NewKey(c, "Domain", "~pupal", 0, nil))
		if err := datastore.Get(c, pu.Key, &pu); err == datastore.ErrNoSuchEntity {
			// User is entirely new and must be added to datastore by a POST /user/registerPupalUser request
			// For now we initialize pupal user with empty fields inside datastore.
			if _, err := datastore.Put(c, pu.Key, &pu); err != nil {
				NewError(w, 500, "Failed to put new user inside the datastore", err, "GetUserFromCache")
				return
			}
		} else if err != nil {
			NewError(w, 500, "Failed to get the pupal user from the datastore", err, "GetUserFromCache")
			return
		} else {
			// Set memcache with pupal user from datastore
			if err := SetCache(c, uid, pu); err != nil {
				NewError(w, 500, "Failed to set pupal user in memcache", err, "GetUserFromCache")
				return
			}
		}
	} else if err != nil {
		NewError(w, 500, "Failed to get the pupal user from memcache", err, "GetUserFromCache")
		return
	}
	gorCtx.Set(r, "PupalUser", pu) // context.Get(r, "PupalUser").(PupalUser) to retrieve it
	next(w, r)
}

// SetCache sets an item inside memcache with key and obj
func SetCache(c ctx.Context, key string, obj interface{}) error {
	item := &memcache.Item{
		Key:    key,
		Object: obj,
	}
	return memcache.Gob.Set(c, item)
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

// NewError supports error handling for each request for console and log.
func NewError(w http.ResponseWriter, code int, msg string, err error, handler string) {
	w.WriteHeader(code)
	msg = msg + ": " + err.Error()
	w.Write([]byte(msg))
	log.Println(handler + ":" + msg)
}
