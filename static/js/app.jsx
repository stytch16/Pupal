class App extends React.Component {
	constructor() {
		super();
		this.state = {
			loggedIn: null
		};
	}
	componentWillMount() {
		this.configFirebase();
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
			});
		} else {
			console.log("Error logging user out (500): User state was lost")
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
					<button type="button" onClick={()=>this.props.onGoogleClick()} className="btn btn-danger loginBtn loginBtn--google">
						Google Login
					</button>
					<button type="button" onClick={()=>this.props.onFacebookClick()} className="btn btn-primary loginBtn loginBtn--facebook">
						Facebook Login
					</button>
				</div>
			</div>
		);
	}
}

function List(props) {
	return (
		<div className="list-group" key="domain_listing">
		{
			props.domains.map((item) => 
				<button type="button" className="list-group-item" key={item} onClick={() => props.onDomainClick(item)}>
					{item}
				</button>)
		}
		</div>
	);
}

class Home extends React.Component {
	constructor() {
		super();
		this.state = {
			option: null, // home menu option
			domain: null, // enter domain name
			belong: null, // user's domain name
			initialDomains: [], // list of all domains
			domains: [] // list of updated domains based on user input
		};
		this.filterList = this.filterList.bind(this);
	}
	componentWillMount() {
		const setDomains = (res) => {
			res = JSON.parse(res).map((item) => item.name);
			this.setState({initialDomains: res, domains: res});
		};
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/list",
				type: "GET",
				beforeSend: function(xhr){xhr.setRequestHeader(
					'Authorization', token);},
				success: (res) => setDomains(res)
			});
			/*
			$.ajax({
				url:"/users/registerDomain",
				type: "POST",
				beforeSend: function(xhr){xhr.setRequestHeader(
					'Authorization', token);},
				success: (res) => regUser(res)
			});
			*/
		});
	}
	componentDidMount() {
		console.log("Mounted!");
		const showDomainModal = () => {
			this.setState({belong:false});
			$('#join_domain_modal').modal('toggle'); };

		var user = firebase.auth().currentUser;
		console.log("Set firebase user");
		firebase.database().ref('users/' + user.uid).set({
			name: user.displayName,
			email: user.email,
			photo: user.photoURL
		});

		console.log("Checking user domain field");
		firebase.database().ref('users/' + user.uid).once('value').then(function(snapshot) {
			if (snapshot.val().domain !== undefined) {
				console.log(snapshot.val().name, " already belongs to a domain = [", snapshot.val().domain,"]");

				// User has a domain
			} else {
				firebase.auth().currentUser.getToken(true).then(function(token) {
					$.ajax({
						url:"/users/registerPupalUser",
						type: "POST",
						contentType: "application/json",
						data: JSON.stringify({name: user.displayName, email: user.email,photo: user.photoURL}),
						beforeSend: function(xhr){xhr.setRequestHeader(
							'Authorization', token);},
						success: () => showDomainModal() 
					});
				});
			}
		});
	}
	handleDomainClick(domain) {
		this.setState({option: 1, domain: domain});
	}
	filterList(event) {
		let updatedList = this.state.initialDomains;
		updatedList = updatedList.filter(function(item) {
			return item.toLowerCase().search(event.target.value.toLowerCase()) !== -1;
		});
		this.setState({domains: updatedList});
	}
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="container">
				<div id="user_info">
					<img src={user.photoURL} id="user-photo" className="img-fluid"></img>
					<h5>{user.email}</h5>
				</div>
				<div className="col-xs-12 text-center">
					<div id="join_domain_modal" className="modal fade" role="dialog">
						<div className="modal-dialog">
							<div className="modal-content">
								<div className="modal-header">
									<button type="button" className="close" data-dismiss="modal">&times;</button>
									<h4 className="modal-title">Search and join your Pupal domain!</h4>
								</div>
								<div className="modal-body">
									<p>
										Your domain can be your school, university, group and/or organization.<br />Pupal associates your projects to your domain(s) while allowing you to subscribe to<br />people and projects from other domains for notifications and updates.
									</p>
								</div>
								<div className="modal-footer">
									<button type="button" className="btn btn-default" data-dismiss="modal">Close</button>
								</div>
							</div>
						</div>
					</div>
					<h2 className="title">
						Welcome to Pupal, {user.displayName} !
					</h2>
					<div className="filtered-list md-form">
						<input type="text" className="form-control " placeholder="Enter domain" onChange={this.filterList} />
						<List domains={this.state.domains} onDomainClick={(domain)=>this.handleDomainClick(domain)} />
					</div>
					<button onClick={()=>this.setState({option:2})} className="hostButton btn btn-default">
						Host a Project
					</button>
					<button onClick={()=>this.setState({option:3})} className="profileButton btn btn-default">
						Go to Profile
					</button>
					<button onClick={()=>this.props.onLogoutClick()} className="logoutBtn btn btn-default">
						Logout
					</button>

					<Domain option={this.state.option} name={this.state.domain} />
					<HostProject option={this.state.option} />
					<GoToProfile option={this.state.option} />
				</div>
			</div>
		);
	}
}

class Domain extends React.Component {

	render() {
		if (this.props.option != 1) {
			return null;
		}
		return (
			<div className="domain_page">
				<h3>Show {this.props.name}'s page !</h3>
			</div>
		)
	}
}

class BrowseProjects extends React.Component {
	render() {
		if (this.props.option != 1) {
			return null;
		} 
		return (
			<div className="browse_projects_page">
				<h1>Browse projects !</h1>
			</div>
		);
		
	}
}

class BrowseUsers extends React.Component {
	render() {
		if (this.props.option != 2) {
			return null;
		} 
		return (
			<div className="browse_users_page">
				<h1>Browse Users !</h1>
			</div>
		);
		
	}
}

class HostProject extends React.Component {
	render() {
		if (this.props.option != 2) {
			return null;
		} 
		return (
			<div className="host_project_page">
				<h3>Host projects !</h3>
			</div>
		);
		
	}
}


class GoToProfile extends React.Component {
	render() {
		if (this.props.option != 3) {
			return null;
		} 
		return (
			<div className="profile_page">
				<h3>Go to Profile !</h3>
			</div>
		);
		
	}
}

ReactDOM.render(
	<App />, 
	document.getElementById('app')
);
