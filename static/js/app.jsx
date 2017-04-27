class App extends React.Component {
	constructor() {
		super();
		this.state = {
			loggedIn: false
		};
	}
	componentWillMount() {
		this.configFirebase();
		this.setupAjax();
		this.setLoginListener();
	}
	configFirebase() {
		var config = {
		    apiKey: "AIzaSyASuOuLQZVYAJnRILcQlE9ImeSTZHezcYk",
		    authDomain: "pupal-164400.firebaseapp.com",
		    databaseURL: "https://pupal-164400.firebaseio.com",
		    projectId: "pupal-164400",
		    storageBucket: "pupal-164400.appspot.com",
		    messagingSenderId: "96889471646"
		};
		firebase.initializeApp(config);
		console.log("Firebase initialized");
	}
	setupAjax() {
		$.ajaxSetup({
			'beforeSend': function(xhr) {
				var user = firebase.auth().currentUser;
				if (user) {
					user.getToken(true).then(function(token) {
						xhr.setRequestHeader('Authorization', 'Bearer ' + token);
					}, function(error) {
						console.log("Error sending AJAX request (" + error.code + "): " + error.message);
					});
				}
			}
		});
		console.log("AJAX Requests set for auth");
	}
	setLoginListener() {
		console.log("Setting up login listener");
		firebase.auth().onAuthStateChanged((user) => {
			if (user) {
				this.setState({loggedIn: true}, ()=>console.log("set Login to ", this.state.loggedIn));
			} else {
				this.setState({loggedIn: false}, ()=>console.log("set Login to ", this.state.loggedIn));
			}
		});
	}
	handleGoogleBtnClick() {
		var provider = new firebase.auth.GoogleAuthProvider();
		firebase.auth().signInWithPopup(provider).then(function(result) {
			console.log(result.user.displayName + " has signed in using Google.");
			}, function(error) {
				console.log("Error authenticating thru Google (" + error.code + "): " + error.message);
				return;
			});
	}
	handleFacebookBtnClick() {
		var provider = new firebase.auth.FacebookAuthProvider();
		firebase.auth().signInWithPopup(provider).then(function(result) {
			console.log(result.user.displayName + " has signed in using FB.");
			}, function(error) {
				console.log("Error authenticating thru FB (" + error.code + "): " + error.message);
				return;
			});
	}
	handleLogoutClick() {
		var user = firebase.auth().currentUser;
		if (user) {
			firebase.auth().signOut().then(function() {
				console.log("User has signed out");
			}, function(error) {
				console.log("Error logging user out (" + error.code + "): " + error.message);
				//alert("Error logging user out (" + error.code + "): " + error.message);
			});
		} else {
			console.log("Error logging user out (500): User state was lost")
			//alert("Error logging user out (500): User state was lost");
		}
	}
	render() {
		if (this.state.loggedIn) {
			return (<Home onLogoutClick={()=>this.handleLogoutClick()} />);	
		} else {
			return (<Login onGoogleClick={()=>this.handleGoogleBtnClick()} onFacebookClick={()=>this.handleFacebookBtnClick()} />);
		}
	}
}

class Login extends React.Component {
	render() {
		return (
			<div className="container">
				<div className="col-xs-12 text-center">
					<h1 className="title">
						Pupal
					</h1>
					<button onClick={()=>this.props.onGoogleClick()} className="loginBtn loginBtn--google">
						Login with Google
					</button>
					<button onClick={()=>this.props.onFacebookClick()} className="loginBtn loginBtn--facebook">
						Login with Facebook
					</button>
				</div>
			</div>
		);
	}
}

class Home extends React.Component {
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="container">
				<div className="col-xs-12 text-center">
					<h1 className="title">
						Welcome to Pupal, {user.displayName} !
					</h1>
					<button onClick={()=>this.props.onLogoutClick()} className="logoutBtn">
						Logout
					</button>
				</div>
			</div>
		);
	}
}

ReactDOM.render(
	<App />, 
	document.getElementById('app')
);
