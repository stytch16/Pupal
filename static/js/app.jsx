class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			loggedIn: null,
			option: 0, // menu option: 0->Home, 1->Domain
			domain: null // domain entered
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
				this.setState({loggedIn: false, option: 0});
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
			console.log("Error logging user out (500): User state was lost");
		}
	}
	handleDomainClick(domain) {
		console.log("Entering " + domain);
		this.setState({option: 1, domain: domain});
	}
	handleHostProjectClick(domain) {
		console.log("Hosting a project at " + domain);
		this.setState({option: 3, domain: domain});
	}
	// App render
	render() {
		console.log("Option = ", this.state.option);
		var user = firebase.auth().currentUser;
		if (this.state.loggedIn) {
			return (
				<div className="container">
					<div className="navbar navbar-default" role="navigation">
						<div className="container-fluid">
							<div className="navbar-header">
								<button className="navbar-toggle" type="button" data-toggle="collapse" data-target=".navbar-collapse"><span className="sr-only">Toggle navigation</span><span className="icon-bar"></span><span className="icon-bar"></span><span className="icon-bar"></span></button><a rel="nofollow" rel="noreferrer"className="navbar-brand" href="#">Pupal</a>
							</div>
							<div className="navbar-collapse collapse">
								<ul className="nav navbar-nav navbar-right">
									<li><a onClick={()=>this.setState({option:0})}>Home</a></li>
									<li className="dropdown"><a rel="nofollow" rel="noreferrer" className="dropdown-toggle" href="#" data-toggle="dropdown">Domains<b className="caret"></b></a>
										<ul className="dropdown-menu">
										</ul>
									</li>
									<li><a onClick={()=>this.setState({option:2})}>Profile</a></li>
									<li><a onClick={()=>this.handleLogoutClick()}>Log out</a></li>
									<img className="user-pic" src={user.photoURL} alt={user.displayName}></img>
								</ul>
							</div>
						</div>
					</div>
					{this.state.option === 0 && <Home onDomainClick={(domain)=>this.handleDomainClick(domain)} />}
					{this.state.option === 1 && <Domain name={this.state.domain} onHostProjectClick={(domain)=>this.handleHostProjectClick(domain)}/>}
					{this.state.option === 2 && <Profile /> }
					{this.state.option === 3 && <HostProject domain={this.state.domain} onBackToDomainClick={(domain)=>this.handleDomainClick(domain)} onProjectPageClick={(id)=>handleProjectPageClick(id)} /> }
				</div>
			);
		} else {
			return (<Login onGoogleClick={()=>this.handleGoogleBtnClick()} onFacebookClick={()=>this.handleFacebookBtnClick()} />)
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

// List render
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
	constructor(props) {
		super(props);
		this.state = {
			hOption: null, // home menu option
			belong: null, // user's domain name
			initialDomains: [], // list of all domains
			domains: [] // list of updated domains based on user input
		};
		this.filterList = this.filterList.bind(this);
	}

	componentDidMount() {
		var user = firebase.auth().currentUser;

		const setDomains = (res) => {
			res = res.map(item => item.name);
			this.setState({initialDomains: res});
		};

		user.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/list",
				type: "GET",
				beforeSend: function(xhr){xhr.setRequestHeader(
					'Authorization', token);},
				success: (res) => setDomains(res)
			});
		});

		const showDomainModal = () => {
			this.setState({belong:false});
			$('#join_domain_modal').modal('toggle'); };

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
	// Home render
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="content">
				<ModalDialog id="join_domain_modal" title="Looks like you need to join a Pupal domain!" body="Your domain can be your school, university, group and/or organization. Pupal associates your projects to your domain(s) while allowing you to subscribe to people and projects from other domains for notifications and updates."/>
				<h2 className="title col-xs-8">
					Welcome to Pupal, {user.displayName} !
				</h2>
				<div className="filtered-list md-form col-xs-4">
					<input type="text" className="form-control" placeholder="Search a domain" onChange={this.filterList} />
					<List domains={this.state.domains} onDomainClick={(domain)=>this.props.onDomainClick(domain)} />
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
	constructor(props) {
		super(props);
		this.state = {
			desc: null,
			photo: null,
			comments: [],
			members: [],
			subscribers: [],
			isMember: null,
			isSubscriber: null
		};
	}
	componentDidMount() {
		this.getDomainInfo(this.props.name);
	}
	getDomainInfo(name) {
		const setStates = (res) => {
			console.log("Setting info of domain, ", name)
			this.setState({
				desc: res.description,
				photo: res.photo_url,
				comments: res.comments,
				members: res.members,
				subscribers: res.subscribers,
				isMember: res.is_member,
				isSubscriber: res.is_subscriber
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
				success: () => {
					console.log("user has joined " + name);
					$('#joinPopover').popover('show');
				}
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
				success: () => { 
					console.log("user has subscribed to " + name);
					$('#subscribePopover').popover('show');
				}
			});
		});
	}
	// Domain render
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="content col-xs-12">
				<div className="domain-header text-center">
					<div className="domain-img">
						<img src={this.state.photo} className="little-round-image img-fluid"></img>
					</div>
					<div className="domain-page">
						<h2>{this.props.name}</h2>
					</div>
					<div className="domain-desc">
						<h5>{this.state.desc}</h5>
					</div>
					{ (!this.state.isMember) ? 
					(<a tabIndex="0" id="joinPopover" className="btn btn-default" role="button" onClick={()=>this.joinDomain(this.props.name)} data-toggle="popover" data-trigger="focus" data-placement="left" data-animation="true" data-content="Success. Host a project!">Join</a>)
					: (<button onClick={()=>this.props.onHostProjectClick(this.props.name)} className="hostButton btn btn-default">Host a Project</button>)}
					{ (!this.state.isSubscriber) ?
					(<a tabIndex="0" id="subscribePopover" className="btn btn-default" role="button" onClick={()=>this.subscribeDomain(this.props.name)} data-toggle="popover" data-trigger="focus" data-placement="bottom" data-animation="true" data-content="Success!">Subscribe</a>)
					: (null) }
				</div>
				<DisplayMemberPhotoPanel members={this.state.members} />
				<DisplaySubscriberPhotoPanel subscribers={this.state.subscribers} />
			</div>
		)
	}
}

// Host Project page
class HostProject extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			title: '', titleFeedback: 'AlphaGo', titleValid: false,
			description: '', descFeedback: 'An AI computer program to play the board game of Go using a Monte Carlo tree search algorithm #machine-learning', descValid: false,
			teamSize: '0', teamSizeFeedback: 'Select how big the team will be', teamSizeValid: false,
			website: '', websiteFeedback: 'https://deepmind.com/research/alphago/  or \'NA\' if you do not have one now', websiteValid: false
		}
		this.handleTitleChange = this.handleTitleChange.bind(this);
		this.handleDescChange = this.handleDescChange.bind(this);
		this.handleTeamSizeChange = this.handleTeamSizeChange.bind(this);
		this.handleWebsiteChange = this.handleWebsiteChange.bind(this);
		this.url_regex = new RegExp(/https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi);
	}
	componentDidMount() {
			$('#success-alert').hide();
			$('#failure-alert').hide();
	}
	handleTitleChange(event) {
		this.setState({title: event.target.value});
		if (event.target.value.length < 3) {
			this.setState({titleFeedback: 'Too short!', titleValid: false});
		} else if (event.target.value.length > 50) {
			this.setState({titleFeedback: 'Too long!', titleValid: false});
		} else {
			this.setState({titleFeedback: 'Great title!', titleValid: true});
		}
	}
	handleDescChange(event) {
		this.setState({description: event.target.value});
		if (event.target.value.length < 30) {
			this.setState({descFeedback: 'Provide more info please!', descValid: false});
		} else if (event.target.value.length > 1000) {
			this.setState({descFeedback: 'Whoa! Try cutting down.', descValid: false});
		} else {
			this.setState({descFeedback: 'Awesome!', descValid: true});
		}
	}
	handleTeamSizeChange(event) {
		this.setState({teamSize: event.target.value});
		this.setState({teamSizeFeedback: 'Got it!', teamSizeValid: true});
	}
	handleWebsiteChange(event) {
		this.setState({website: event.target.value});
		if (event.target.value.match(this.url_regex)) {
			this.setState({websiteFeedback: 'Website looks good!', websiteValid: true});
		} else if (event.target.value.localeCompare('NA') == 0) {
			this.setState({websiteFeedback: 'Got it! Hopefully you get a website soon!', websiteValid: true});
		} else {
			this.setState({websiteFeedback: 'Website does not look correct. Enter \'NA\' if you do not have one now.', websiteValid: false});
		}
	}
	handleSubmitClick() {
		if (this.state.titleValid && this.state.descValid && this.state.teamSizeValid && this.state.websiteValid) {
			$('#failure-alert').hide();
			$('#success-alert').show();
		} else {
			$('#failure-alert').show();
		}
	}
	handleResetClick() {
		this.setState({
			title: '', titleFeedback: 'AlphaGo', titleValid: false,
			description: '', descFeedback: 'an AI computer program to play the board game of Go using a Monte Carlo tree search algorithm #machine-learning', descValid: false,
			teamSize: '0', teamSizeValid: false,
			website: '', websiteFeedback: 'https://deepmind.com/research/alphago/', websiteValid: false
		});
		$('#failure-alert').hide();
		$('#success-alert').hide();
	}
	render() {
		return (
			<div className="content col-xs-8">
				<h3 className="title">You are hosting a project at {this.props.domain}! </h3>
					<div id="success-alert" className="alert alert-success" role="alert">
						<strong>Well done!</strong> You have successfully submitted a project to {this.props.domain}
					</div>
					<div id="failure-alert" className="alert alert-danger" role="alert">
						<strong>Oh snap!</strong> Change a few things, check your helper text, and try submitting again.
					</div>

					<form className="col">
						<div className="form-group">
							<label htmlFor="titleInput"><h3>Title</h3></label>
							<input type="text" 
								value={this.state.title} 
								onChange={this.handleTitleChange} 
								className="form-control" id="titleInput" 
								aria-describedby="titleHelp" 
								placeholder="Got a good name?"></input>
							<p id="titleHelp" className="form-text text-muted">
								{this.state.titleFeedback}</p>
						</div>
						<div className="form-group">
							<label htmlFor="descriptionInput"><h3>Description</h3></label>
							<textarea value={this.state.description} 
								onChange={this.handleDescChange} 
								className="form-control" id="descriptionInput" 
								rows="5" 
								aria-describedby="descriptionHelp" 
								placeholder="Describe your project and include your hashtags"></textarea>
							<p id="descriptionHelp" className="form-text text-muted">
								{this.state.descFeedback}</p>
						</div>
						<div className="form-group">
							<label htmlFor="teamSizeInput"><h3>Size of project team</h3></label>
							<select value={this.state.teamSize} 
								onChange={this.handleTeamSizeChange} 
								className="form-control" id="teamSizeInput">
								<option value="1">1</option>
								<option value="1-3">1-3</option>
								<option value="3-5">3-5</option>
								<option value="5-10">5-10</option>
								<option value="10+">10+</option>
							</select>
							<p id="descriptionHelp" className="form-text text-muted">
								{this.state.teamSizeFeedback}</p>
						</div>
						<div className="form-group">
							<label htmlFor="websiteInput"><h3>Got a website?</h3></label>
							<input type="text" 
								value={this.state.website} 
								onChange={this.handleWebsiteChange} 
								className="form-control" id="websiteInput" 
								placeholder="Got a website for this project?"></input>
							<p id="websiteHelp" className="form-text text-muted">
								{this.state.websiteFeedback}</p>
						</div>

						<button onClick={()=>this.handleSubmitClick()} 
							type="submit" className="btn btn-success">Host my project! Mmhmm!</button>
						<button onClick={()=>this.handleResetClick()} 
							type="submit" className="btn btn-warning">Ugh! Start over.</button>
						<button onClick={()=>this.props.onBackToDomainClick(this.props.domain)} 
							type="submit" className="btn btn-danger">Rut ro. Cancel.</button>

					</form>
			</div>
		);
	}
}


class Profile extends React.Component {
	render() {
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
