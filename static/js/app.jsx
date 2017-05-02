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
				this.setState({loggedIn: true});
			} else {
				this.setState({loggedIn: false});
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

// Login page
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

// List of button display
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

function ModalDialog(props) {
	return (
		<div id={props.id} className="modal fade" role="dialog">
			<div className="modal-dialog">
				<div className="modal-content">
					<div className="modal-header">
						<button type="button" className="close" data-dismiss="modal">&times;</button>
						<h4 className="modal-title">{props.title}</h4>
					</div>
					<div className="modal-body">
						<p>
							{props.body}
						</p>
					</div>
					<div className="modal-footer">
						<button type="button" className="btn btn-default" data-dismiss="modal">Close</button>
					</div>
				</div>
			</div>
		</div>
	);
}

// Home page
class Home extends React.Component {
	constructor() {
		super();
		this.state = {
			option: null, // home menu option
			belong: null, // user's domain name
			domain: null, // enter domain
			initialDomains: [], // list of all domains
			domains: [] // list of updated domains based on user input
		};
		this.filterList = this.filterList.bind(this);
	}
	componentWillMount() {
		const setDomains = (res) => {
			res = res.map(item => item.name);
			this.setState({initialDomains: res});
		};
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/list",
				type: "GET",
				beforeSend: function(xhr){xhr.setRequestHeader(
					'Authorization', token);},
				success: (res) => setDomains(res)
			});
		});
	}
	componentDidMount() {
		console.log("Mounted!");
		const showDomainModal = () => {
			this.setState({belong:false});
			$('#join_domain_modal').modal('toggle'); };

		var user = firebase.auth().currentUser;
		console.log("Set firebase user");

		firebase.database().ref('users/' + user.uid).once('value').then(function(snapshot) {
			if (snapshot.val() === null) {
				console.log("User is new");
				firebase.database().ref('users/'+user.uid).set({
					name: user.displayName,
					email: user.email,
					photo: user.photoURL,
					domain: false,
					messages: []
				});
				user.getToken(true).then(function(token) {
					$.ajax({
						url:"/users/registerPupalUser",
						type: "POST",
						contentType: "application/json",
						data: JSON.stringify(
							{
							name: user.displayName, 
							email: user.email,
							photo: user.photoURL
							}),
						beforeSend: function(xhr) {
							xhr.setRequestHeader('Authorization', token);
							},
						success: () => showDomainModal() 
					});
				});
			} else if (snapshot.val().domain === false) {
				console.log("User has no domain but is already a Pupal user");
				showDomainModal()
			}
		});
	}
	handleDomainClick(domain) {
		console.log("Mounting Domain")
		this.setState({option:1, domain:domain});
	}
	filterList(event) {
		if (event.target.value.length !== 0) {
			let updatedList = this.state.initialDomains;
			updatedList = updatedList.filter(function(item) {
				return item.toLowerCase().search(event.target.value.toLowerCase()) !== -1;
			});
			this.setState({domains: updatedList});
		} else {
			this.setState({domains: []});
		}
	}
	render() {
		var user = firebase.auth().currentUser;
		if (this.state.option === 1) {
			return <Domain name={this.state.domain} onLogoutClick={()=>this.props.onLogoutClick()} />
		}
		return (
			<div className="container">
				<nav className="navbar navbar-default">
					<div className="container-fluid">
						<div className="navbar-header">
							<a className="navbar-brand" href="#">Pupal</a>
						</div>
						<div className="navbar-collapse collapse">
							<ul className="nav navbar-nav navbar-right">
								<li className="active"><a href="/">Home</a></li>
								<li><a href="#">About</a></li>
								<li><a href="#">Contact</a></li>
								<img className="user-pic" src={user.photoURL} alt="User"></img>
							</ul>
						</div>
					</div>
				</nav>
					
				<div className="content">
					<ModalDialog id="join_domain_modal" title="Looks like you need to join a Pupal domain!" body="Your domain can be your school, university, group and/or organization.<br /><br />Pupal associates your projects to your domain(s) while allowing you to subscribe to<br />people and projects from other domains for notifications and updates."/>
					<h2 className="title col-xs-8">
						Welcome to Pupal, {user.displayName} !
					</h2>
					<div className="filtered-list md-form col-xs-4">
						<input type="text" className="form-control " placeholder="Enter domain" onChange={this.filterList} />
						<List domains={this.state.domains} onDomainClick={(domain)=>this.handleDomainClick(domain)} />
					</div>
					
					<div className="domain-container">
					</div>

					<div className="user-options text-center">
						<button onClick={()=>this.setState({option:2})} className="hostButton btn btn-default">
							Host a Project
						</button>
						<button onClick={()=>this.setState({option:3})} className="profileButton btn btn-default">
							Go to Profile
						</button>
						<button onClick={()=>this.props.onLogoutClick()} className="logoutBtn btn btn-default">
							Logout
						</button>
					</div>
					
					<HostProject option={this.state.option} />
					<GoToProfile option={this.state.option} />
				</div>
			</div>
		);
	}
}

// Display little photos
function DisplayPhotoPanel(props) {
	return (
		<div className="photo_array">
		{
			props.users.map((item) =>
				<img key={item.photo} className="super-little-rcorner-image img-fluid" src={item.photo}></img>)
		}
		</div>
	)
}

function DisplayMemberPhotoPanel(props) {
	return (
		<div className="members-panel panel panel-default col-xs-6">
			<div className="panel-body">
				{(props.members.length !== 0) ? (
					<div className="photo-display-panel">
						<h5><i>{props.members.length} member(s).</i></h5>
						<DisplayPhotoPanel users={props.members} />
					</div>
				) : (
					<h5><i>No members yet. Join or refer someone to join!</i></h5>
				)}
			</div>
		</div>
	);
}

function DisplaySubscriberPhotoPanel(props) {
	return (
		<div className="subscribers-panel panel panel-default col-xs-6">
			<div className="panel-body">
				{(props.subscribers.length !== 0) ? (
					<div className="photo-display-panel">
						<h5><i>{props.subscribers.length} subscriber(s).</i></h5>
						<DisplayPhotoPanel users={props.subscribers} />
					</div>
				) : (
					<h5><i>No subscribers yet. Be the first to subscribe!</i></h5>
				)}
			</div>
		</div>
	);
}

// Domain page
class Domain extends React.Component {
	constructor() {
		super();
		this.state = {
			desc: null,
			photo: null,
			comments: [],
			members: [],
			subscribers: []
		};
	}
	componentWillMount() {
		this.getDomainInfo(this.props.name);
	}
	getDomainInfo(name) {
		console.log("Getting info of domain, ", name)
		const setStates = (res) => {
			this.setState({
				desc: res.description,
				photo: res.photo_url,
				comments: res.comments,
				members: res.members,
				subscribers: res.subscribers
				});
		};
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+name,
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setStates(res)
			});
		});
	}
	joinDomain(name) {
		var user = firebase.auth().currentUser;
		console.log("Joining", name);
		user.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+name+"/join",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: () => console.log("user has joined " + name)
			});
		});
		firebase.database().ref('users/' + user.uid).update({
			domain: true
		});

	}
	subscribeDomain(name) {
		console.log("Subscribing to", name);
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+name+"/subscribe",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: () => console.log("user has subscribed to " + name)
			});
		});
	}
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="container col-xs-12">
				<div className="content text-center">
					<div className="domain-img">
						<img src={this.state.photo} className="little-round-image img-fluid"></img>
					</div>
					<div className="domain-page">
						<h2>{this.props.name}</h2>
					</div>
					<div className="domain-desc">
						<h5>{this.state.desc}</h5>
					</div>
					<button onClick={()=>this.joinDomain(this.props.name)} className="joinBtn btn btn-default">
						Join
					</button>
					<button onClick={()=>this.subscribeDomain(this.props.name)} className="subscribeBtn btn btn-default">
						Subscribe
					</button>
					<button onClick={()=>this.props.onLogoutClick()} className="logoutBtn btn btn-default">
						Logout
					</button>
				</div>

				<DisplayMemberPhotoPanel members={this.state.members} />
				<DisplaySubscriberPhotoPanel subscribers={this.state.subscribers} />
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
