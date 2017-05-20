// begin app component
class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			loggedIn: null
		};
	}
	componentWillMount() {
		this.setLoginListener();
	}
	setLoginListener() {
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
		if (this.state.loggedIn) {
			firebase.auth().signOut().then(function() {
				console.log("User has signed out");
			}, function(error) {
				console.log("Error logging user out (" + error.code + "): " + error.message);
			});
		} else {
			console.log("Error logging user out (500): User state was lost");
		}
	}
	render() {
		if (this.state.loggedIn===true) {
			return (
				<div className="container">
					<Navbar onLogoutClick={()=>this.handleLogoutClick()}/>
					{this.props.children || <Home/>}
				</div>
			)
		} else if (this.state.loggedIn===false) {
			return (
				<div className="container">
					<Login onGoogleClick={()=>this.handleGoogleBtnClick()} 
						onFacebookClick={()=>this.handleFacebookBtnClick()} />
				</div>
			)
		} else {
			return (
				<div className="container">
					Loading ...
				</div>
			)
		}
	}
}
// end app component

// begin login function component
class Login extends React.Component{
	componentDidMount() {
		$('#login-modal').modal({
			backdrop: 'static',
			keyboard: false
		}, 'toggle')
	}
	render () {
		return (
			<div id="login-modal" className="modal fade" role="dialog">
				<div className="modal-dialog">
					<div className="loginmodal-content text-center">
						<div className="modal-header">
							<img id="pupal-login-image" src="/static/images/favicon-32x32.png" alt="Pupal Login"></img>
							<h3>Pupal</h3>
						</div>
						<div className="modal-body">
							<div className="btn-group">
								<a className="btn btn-danger disabled">
									<i className="fa fa-google"></i></a>
								<a className="btn btn-danger login-prompt-btn"
									data-dismiss="modal"
									onClick={()=>this.props.onGoogleClick()}>
									Sign in with Google</a>
							</div>
							<br/><br/>
							<div className="btn-group">
								<a className="btn btn-primary disabled">
									<i className="fa fa-facebook"></i></a>
								<a className="btn btn-primary login-prompt-btn"
									data-dismiss="modal"
									onClick={()=>this.props.onFacebookClick()}>
									Sign in with Facebook</a>
							</div>
							<br/><br/>
							<p>Your Pupal account will be set up from your Google or Facebook account.</p>
						</div>
					</div>
				</div>
			</div>
		)
	}
}
// end login function component

// begin navbar function component
function Navbar(props) {
	var user = firebase.auth().currentUser
	return (
		<nav id="header-nav" className="navbar navbar-default navbar-fixed-top" role="navigation">
			<div className="container-fluid">
				<div className="navbar-header">
					<button className="navbar-toggle collapsed" type="button" data-toggle="collapse" 
						data-target="#header-navbar-content" aria-expanded="false" aria-controls="navbar">
						<span className="sr-only">Toggle navigation</span>
						<span className="icon-bar"></span>
						<span className="icon-bar"></span>
						<span className="icon-bar"></span>
					</button>
					<a className="navbar-brand"><Link to="/">Pupal</Link></a>
				</div>
				<div id="header-navbar-content" className="navbar-collapse collapse">
					<ul className="nav navbar-nav navbar-right">
						<li className="dropdown">
							<a className="dropdown-toggle" data-toggle="dropdown" role="button" 
								aria-haspopup="true" aria-expanded="false">
								Domains<span className="caret"></span></a>
							<ul className="dropdown-menu">
								<UserDomains />
							</ul>
						</li>
						<li><Link to="/profile">Profile</Link></li>
						<li><a onClick={()=>props.onLogoutClick()}>Log out</a></li>
						<img id="navbar-user-pic" src={user.photoURL} alt={user.displayName}></img>
					</ul>
				</div>
			</div>
		</nav>
	)
}
// end navbar function component

// begin userdomains component
class UserDomains extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			domains: null
		};
	}
	componentDidMount() {
		const setDomains = (res) => { this.setState({domains: res}) }
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/userlist",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setDomains(res)
			});
		});
	}
	render() {
		return (
			<div className="user-domain-listing-content">
				{(this.state.domains !== null && this.state.domains.length > 0) ? (
					<List domains={this.state.domains}/>
				) : (
					<a type="button" className="list-group-item disabled">You have no domains.</a>
				)}
			</div>
		)
	}
}
// end user domains component

// begin home component
class Home extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			initialDomains: [], // list of all domains
			filteredDomains: [] // list of updated domains based on user input
		};

		this.filterList = this.filterList.bind(this);
	}
	componentDidMount() {
		var user = firebase.auth().currentUser;

		const setDomains = (res) => { this.setState({initialDomains: res}); }
		// Get JSON array of id:name pairs of domain listing
		user.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/list",
				type: "GET",
				beforeSend: function(xhr){xhr.setRequestHeader(
					'Authorization', token);},
				success: (res) => setDomains(res)
			});
		});

		const welcomeNewUser = () => { $('#welcome-pupal-modal').modal('toggle'); }
		// Read user's record on firebase db. 
		// If nonexistent, register pupal user on GAE datastore.
		// If exist but no associated domain, show modal dialog to user to join a domain.
		firebase.database().ref('users/' + user.uid).once('value').then(function(snapshot) {
			if (snapshot.val() === null) {
				console.log("User is new");
				firebase.database().ref('users/'+user.uid).set({
					domain: false,
					email: user.email,
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
							xhr.setRequestHeader('Authorization', token); },
						success: () => welcomeNewUser()
					});
				});
			} 
		});
	}
	filterList(event) {
		if (event.target.value.length !== 0) {
			let updatedList = this.state.initialDomains;
			updatedList = updatedList.filter(function(item) {
				return item.name.toLowerCase().search(event.target.value.toLowerCase()) !== -1;
			});
			this.setState({filteredDomains: updatedList});
		} else {
			this.setState({filteredDomains: []});
		}
	}
	render() {
		var user = firebase.auth().currentUser;
		return (
			<div className="content home-content">
				<WelcomePupal />
				<div className="recent-activity-content col-xs-8">
					<h3>Recent Activity</h3>
				</div>
				<div className="filtered-list md-form col-xs-4">
					<input type="text" className="form-control" 
						placeholder="Search a domain" 
						onChange={this.filterList} />
					<List domains={this.state.filteredDomains} />
				</div>
			</div>
		);
	}
}
// end home component

// begin list function component
function List(props) {
	function fetchDomLink(id) {
		return "/dom/" + id + "?view=Info"
	}
	return (
		<div className="list-group" key="domain_listing">
		{
			props.domains.map((item) => 
				<a type="button" className="list-group-item" key={item.id}>
					<Link to={fetchDomLink(item.id)}>{item.name}</Link>
				</a>)
		}
		</div>
	);
}
// end list function component

// begin welcomepupal component
function WelcomePupal() {
	return (
		<div id="welcome-pupal-modal" className="modal fade" role="dialog">
			<div className="modal-dialog">
				<div className="modal-content text-center">
					<div className="modal-header">
						<button type="button" className="close" data-dismiss="modal">&times;</button>
						<h4 className="modal-title">
							Welcome to Pupal! To get started, join a domain!</h4>
					</div>
					<div className="modal-body">
						<p>
							Your domain can be your school, group and/or organization.<br/>
							Search for your domain on the right of the page and join!</p>
					</div>
					<div className="modal-footer">
						<button type="button" className="btn btn-default" data-dismiss="modal">
							Close</button>
					</div>
				</div>
			</div>
		</div>
	);
}
// end welcomepupal component

// begin domain component
class Domain extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			domain: null,
			member: null
		};
	}
	componentDidMount() {
		this.getDomainInfo(this.props.params.id);
		this.handleMembership(this.props.params.id);
	}
	getDomainInfo(id) {
		const setStates = (res) => {
			this.setState({domain: res});
		};
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+id,
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setStates(res)
			});
		});
	}
	handleMembership(id) {
		const setMember = (res) => {
			if (res === "true") {
				this.setState({member: true})
			} else {
				this.setState({member:false})
			}
		}
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+id+"/member",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setMember(res)
			});
		});

	}
	handleJoin(id) {
		const setJoinState = () => { 
			this.setState({member:true})
		};
		var user = firebase.auth().currentUser;
		user.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+id+"/join",
				type: "POST",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: () => setJoinState()
			});
		});
		firebase.database().ref('users/' + user.uid).update({
			domain: true
		});
	}
	render() {
		if (this.state.domain === null) {
			return (
				<div className="domain-content">
					Loading...
				</div>
			)
		}
		var user = firebase.auth().currentUser;
		return (
			<div className="content domain-content">
				<div className="domain-header">
					<br/>
					<h1>
						<img src={this.state.domain.photo_url} className="domain-image img-fluid"></img>
						<br/>
						{this.state.domain.name}
					</h1>
					<h5>{this.state.domain.description}</h5>
					<br/><br/>
				</div>
				<div className="domain-navbar">
					<DomainNavbar id={this.props.params.id} view={this.props.location.query.view} member={this.state.member}
						onJoinClick={(id)=>this.handleJoin(id)} />
				</div>
				<div className="domain-view-content">
					<DomainView id={this.props.params.id} view={this.props.location.query.view} proj={this.props.location.query.proj}/>
				</div>
			</div>
		)
	}
}
// end domain component

// begin domainview component
function DomainView(props) {
	if (props.view === "Info") {
		return <Info id={props.id} />
	} else if (props.view === "Projects") {
		return <Projects id={props.id} proj={props.proj} />
	} else if (props.view === "Users") {
		return <Users id={props.id} />
	} else if (props.view === "Host") {
		return <Host id={props.id} />
	} else {
		return null
	}
}
// end domainview component

// begin domainnavbar component
class DomainNavbar extends React.Component{
	fetchViewLink(id, view) {
		return "/dom/"+id+"?view="+view
	}
	isActive(view) {
		return (view === this.props.view) ? "active" : ""
	}
	render() {
		return (
		<nav className="navbar navbar-default">
			<div className="container-fluid">
				<div className="navbar-header">
					<button type="button" className="navbar-toggle collapsed" 
						data-toggle="collapse" data-target="#domain-navbar-content" 
						aria-expanded="false" aria-controls="navbar">
						<span className="sr-only">Toggle navigation</span>
						<span className="icon-bar"></span>
						<span className="icon-bar"></span>
						<span className="icon-bar"></span>
					</button>
				</div>
				<div id="domain-navbar-content" className="navbar-collapse collapse">
					<ul className="nav navbar-nav">
						<li className={this.isActive('Info')}>
							<Link to={this.fetchViewLink(this.props.id, "Info")}>Info</Link></li>
						<li className={this.isActive('Projects')}>
							<Link to={this.fetchViewLink(this.props.id, "Projects")}>Projects</Link></li>
						<li className={this.isActive('Users')}>
							<Link to={this.fetchViewLink(this.props.id, "Users")}>Users</Link></li>
						<li className={this.isActive('Comments')}><a href="#">Comments</a></li>
					</ul>
					<ul className="nav navbar-nav navbar-right">
						<li className="dropdown">
							<a className="dropdown-toggle" data-toggle="dropdown" role="button" 
								aria-haspopup="true" aria-expanded="false">
								Actions<i className="fa fa-bars" aria-hidden="true"></i></a>
							<ul className="dropdown-menu">
								{!this.props.member && <li><a onClick={()=>this.props.onJoinClick(this.props.id)}>
									<i className="fa fa-tags" aria-hidden="true"></i>
									Request to join</a></li>}
								<li role="separator" className="divider"></li>
								{this.props.member && <li><Link to={this.fetchViewLink(this.props.id, "Host")}>
									<i className="fa fa-paper-plane" aria-hidden="true"></i>
									Host a project</Link></li>}
							</ul>
						</li>
					</ul>
				</div>
			</div>
		</nav>
		)
	}
}
// end domain navbar function component

// begin info component
class Info extends React.Component {
	render() {
		return (
			<div className="info-content">
				<h2>Display domain of id={this.props.id} info here!</h2>
			</div>
		)
	}
}
// end info component

// begin projects component
class Projects extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			projects: [],
			proj: null
		};
	}
	componentDidMount() {
		this.getProjects(this.props.id) // Get projects of domain id in URL
		this.getProject(this.props.proj) // Get project of project id in URL if exists
	}
	getProjects(id) {
		const setProjects = (res) => { this.setState({ projects: res }) }
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+id+"/projectlist",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setProjects(res)
			});
			
		});
	}
	getProject(id) {
		const setProjModal = (res) => { this.setState({proj: res}, ()=>{$("#proj-modal").modal("toggle")}) }
		if (id !== undefined && id !== null) {
			firebase.auth().currentUser.getToken(true).then(function(token) {
				$.ajax({
					url: "/projects/"+id,
					type: "GET",
					beforeSend: function(xhr){
						xhr.setRequestHeader('Authorization', token);
					},
					success: (res) => setProjModal(res)
				});
			});
		}
	}
	handleProjSubClick(id) {
		const resetProj = () => { this.getProject(this.props.id) }
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/projects/"+id+"/subscribe",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: () => resetProj()
			});
		});
	}
	handleUserClick(id) {
		// Redirect to the user page
		hashHistory.push("/user/"+id)
	}
	render() {
		return (
			<div className="projects-content">
				<ProjModal proj={this.state.proj} 
					onProjSubClick={()=>this.handleProjSubClick(this.state.proj.id)}
					onUserClick={(id)=>this.handleUserClick(id)}/>
				<div className="project-group" key="projects-listing">
				{
					this.state.projects.map((proj) => 
					<div className="proj-entry" key={proj.id}>
						<a onClick={()=>this.getProject(proj.id)}>
						<div className="card proj-entry-content">
							<div className="card-block">
								<h3 className="card-title">{proj.title}</h3>
								<div className="card-text">
									<p className="num-subscribes-text">
										<i>{proj.num_subscribes} subscriber(s)</i></p>
									<br />
									<p className="desc-text">{proj.description}</p>
									{proj.tags !== null && proj.tags.length > 0 && proj.tags.map((tag) => 
									<div className="proj-entry-tag" key={tag}>
									<p><i className="fa fa-tag" aria-hidden="true"></i>
										{tag}</p></div>)
									}
									<br />
								</div>
							</div>
						</div>
						</a>
					</div>
					)
				}
				</div>
			</div>
		);
	}
}
// end projects component

// begin projmodal component
function ProjModal(props) {
	if (props.proj !== null) {
		return (
			<div id="proj-modal" className="modal fade" tabIndex="-1" role="dialog" 
				aria-labelledby="proj-modal-label" aria-hidden="true">
				<div className="modal-dialog modal-lg modal-notify modal-info" role="document">
					<div className="modal-content">
						<div className="modal-header text-center">
							<button type="button" className="close" data-dismiss="modal">&times;
							</button>
							<h1 className="modal-title">{props.proj.title}</h1>
							<button type="button" 
								className="btn btn-default btn-circle btn-lg"
								onClick={()=>this.props.onProjSubClick()}>
							<i className="fa fa-star-o" aria-hidden="true"></i>
							</button>
						</div>
						<div className="modal-body">
							<div className="desc-header-info">
								<br />
								<h5>{props.proj.created_at}<br /></h5>
								<h4>Team size: {props.proj.team_size}</h4>
								<br/>
							</div>
							<div className="desc-info">
								<h4>{props.proj.description}</h4>
								<br/>
							</div>
							<div className="author-info">
								<a onClick={()=>props.onUserClick(props.proj.author.pupal_id)} 
									data-dismiss="modal" >
									<img className="proj-author-image img-fluid" 
										src={props.proj.author.photo} alt={props.proj.author.name}></img>
									<div className="author-info-contact">
										<h4>{props.proj.author.name}</h4>
									</div>
								</a>
								<br/><br/>
							</div>
						</div>
						<div className="modal-footer">
							{(props.proj.website.localeCompare("NA") !== 0) && <a href={props.proj.website} className="proj-website">{props.proj.website}</a>}
						</div>
					</div>
				</div>
			</div>
		);
	}
	return null
}
// end projmodal component

// begin users component
class Users extends React.Component {
	render() {
		return (
		<div className="users-content">
			<h2>Display domain of id={this.props.id} users here!</h2>
		</div>
		);
	}
}
// end users component

// begin host component
class Host extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			title: '', 
			titleFeedback: '*Example: AlphaGo', 
			titleValid: false,
			description: '', 
			descFeedback: '*Example: an AI computer program to play the board game of Go using a Monte Carlo tree search algorithm #machine-learning', 
			descValid: false,
			teamSize: '1-3', 
			teamSizeFeedback: '*Select how big the team will be.', 
			website: '', 
			websiteFeedback: '*Example: https://deepmind.com/research/alphago/ OR enter \'NA\' if you do not have one now', 
			websiteValid: false,

			projId: null
		}
		this.handleTitleChange = this.handleTitleChange.bind(this);
		this.handleDescChange = this.handleDescChange.bind(this);
		this.handleTeamSizeChange = this.handleTeamSizeChange.bind(this);
		this.handleWebsiteChange = this.handleWebsiteChange.bind(this);
		this.url_regex = new RegExp(/https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)?/gi);
	}
	componentDidMount() {
		$('#failure-alert').hide();
		$('#success-alert').hide();
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
		this.setState({teamSize: event.target.value, teamSizeFeedback: 'Got it!'});
	}
	handleWebsiteChange(event) {
		this.setState({website: event.target.value});
		if (event.target.value.match(this.url_regex)) {
			this.setState({websiteFeedback: 'Website looks good!', websiteValid: true});
		} else if (event.target.value.localeCompare('NA') == 0) {
			this.setState({websiteFeedback: 'Okay! Hopefully you get a website soon!', websiteValid: true});
		} else {
			this.setState({websiteFeedback: 'Website does not look correct. Enter \'NA\' if you do not have one now.', websiteValid: false});
		}
	}
	handleSubmitClick(id, titl, desc, ts, web) {
		const setProject = (res) => {
			$('#failure-alert').hide();
			$('#success-alert').show();
			this.setState({projId: res})
		};
		if (this.state.titleValid && this.state.descValid && this.state.websiteValid) {
			firebase.auth().currentUser.getToken(true).then(function(token) {
				$.ajax({
					url: "/projects/"+id+"/host",
					type: "POST",
					contentType: "application/json",
					beforeSend: function(xhr) {
						xhr.setRequestHeader('Authorization', token);
					},
					data: JSON.stringify(
						{
						title: titl,
						description: desc,
						team_size: ts,
						website: web
						}),
					success: (res) => setProject(res)
				});
			});
		} else {
			$('#failure-alert').show();
		}
	}
	handleResetClick() {
		$('#failure-alert').hide();
		this.setState({
			title: '', 
			titleFeedback: '*Example: AlphaGo', 
			titleValid: false,
			description: '', 
			descFeedback: '*Example: an AI computer program to play the board game of Go using a Monte Carlo tree search algorithm #machine-learning', 
			descValid: false,
			teamSize: '1-3', 
			teamSizeFeedback: '*Select how big the team will be.', 
			website: '', 
			websiteFeedback: '*Example: https://deepmind.com/research/alphago/ OR enter \'NA\' if you do not have one now', 
			websiteValid: false
		})
	}
	fetchNewProjLink(id) {
		return "/dom/"+id+"?view=Projects&proj="+this.state.projId
	}
	render() {
		return (
		<div className="host-content">
			<form className="col">
				<div className="form-group">
					<label htmlFor="titleInput"><h3>Title</h3></label>
					<input type="text" 
						value={this.state.title} 
						onChange={this.handleTitleChange} 
						className="form-control" id="titleInput" 
						aria-describedby="titleHelp" 
						placeholder="Got a good name for your project?"></input>
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
						placeholder="Describe your project and include any single-word hashtags to attach (max 5)!"></textarea>
					<p id="descriptionHelp" className="form-text text-muted">
						{this.state.descFeedback}</p>
				</div>
				<div className="form-group">
					<label htmlFor="teamSizeInput"><h3>Size of project team</h3></label>
					<select value={this.state.teamSize} 
						onChange={this.handleTeamSizeChange} 
						className="form-control" id="teamSizeInput">
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
				<div id="host-submit-buttons">
					<a onClick={()=>this.handleSubmitClick(
						this.props.id, this.state.title, this.state.description, 
						this.state.teamSize, this.state.website)} 
						type="button" className="btn btn-success">
						<i className="fa fa-check" aria-hidden="true"></i>Host my project!</a>
					<a onClick={()=>this.handleResetClick()} 
						type="button" className="btn btn-warning">
						<i className="fa fa-exclamation-triangle" aria-hidden="true"></i>Start over</a>
					<a type="button" className="btn btn-danger"><Link to={this.fetchNewProjLink(this.props.id)}>
						<i className="fa fa-times" aria-hidden="true"></i>Cancel</Link></a>
				</div>
			</form>
			<div id="failure-alert" className="alert alert-danger" role="alert">
				<h4><strong>Oh snap!  </strong>
					Look at your helper text and try submitting again.</h4>
			</div>
			<div id="success-alert" className="alert alert-success" role="alert">
				<h4><strong>Excellent!  </strong>
					<Link to={this.fetchNewProjLink(this.props.id)}>
						Click here to check out your project at the project page!
					</Link></h4>
			</div>
		</div>
		);
	}
}
// end host component

// begin user component
class User extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			user: null
		};
	}
	componentDidMount() {
		this.getPupalUser(this.props.params.id)
	}
	getPupalUser(id) {
		const setUser = (res) => {
			this.setState({user: res})
		}
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/users/"+id,
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setUser(res)
			});
		});
	}
	render() {
		if (this.state.user) {
			return (
				<div className="content user-content">
					<div className="header">
						<img id="user-profile-pic" src={this.state.user.photo} alt={this.state.user.name}></img>
						<h1 id="user-profile-name">{this.state.user.name}</h1>
					</div>
					<div className="body col-xs-8">
						<br /><br />
						<h4 id="user-profile-summary">{this.state.user.summary}</h4>
						<br /><br />
					</div>
				</div>
			)
		} else {
			return (
				<div className="content user-content">
					Loading...
				</div>
			)
		}
	}
}
// end user component

// begin profile component
class Profile extends React.Component {
	render() {
		return (
			<div className="content profile_content">
				<h2>Display profile here</h2>
			</div>
		);
	}
}
// end profile component

// Firebase config
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

// React router components
var { Router,
      Route,
      hashHistory,
      Link } = ReactRouter;

// React router config
// Use this.props.params.id to access URL id parameter.
// If query string /foo?bar=baz, use this.props.location.query.bar to get value of bar -> baz

// Query options -> dom/:id?view=Info dom/:id?view=Projects dom/:id?view=Projects&proj=<id> dom/:id?view=Users dom/:id?view=Host
ReactDOM.render((
	<Router history={hashHistory}>
		<Route path="/" component={App}>
			<Route path="dom/:id" component={Domain}/>
			<Route path="user/:id" component={User}/>
			<Route path="profile" component={Profile}/>
		</Route>
	</Router>),
	document.getElementById('app')
);
