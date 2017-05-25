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
			}, function(error) {
				alert("Error authenticating thru Google (" + error.code + "): " + error.message);
				return;
			});
	}
	handleFacebookBtnClick() {
		var provider = new firebase.auth.FacebookAuthProvider();
		firebase.auth().signInWithPopup(provider).then(function(result) {
			console.log(result.user.displayName + " has signed in using FB.");
			}, function(error) {
				alert("Error authenticating thru Facebook (" + error.code + "): " + error.message);
				return;
			});
	}
	handleLogoutClick() {
		if (this.state.loggedIn) {
			firebase.auth().signOut().then(function() {
			}, function(error) {
				alert("Error logging user out (" + error.code + "): " + error.message);
			});
		} else {
			alert("Error logging user out (500): User state was lost");
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
							<h1>Pupal</h1>
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
			domains: []
		};
	}
	componentDidMount() {
		const updateUserDomains = (domains) => { this.setState({domains: domains}) };
		var ref = firebase.database().ref("users/" + firebase.auth().currentUser.uid + "/domains")
		if (ref) {
			let domains = this.state.domains
			ref.on('child_added', function(snapshot) {
				domains.push(snapshot.val())
				updateUserDomains(domains)
			})
		}
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
			filteredDomains: [], // list of updated domains based on user input
			messages: []
		};
		this.filterList = this.filterList.bind(this);
	}
	componentDidMount() {
		this.registerSignIn(firebase.auth().currentUser)
		this.getDomainList()
		this.setMessagesListener()
	}
	setMessagesListener() {
		const updateMessages = (messages) => { this.setState({messages: messages}) };
		var ref = firebase.database().ref("users/" + firebase.auth().currentUser.uid + "/messages")
		if (ref) {
			let messages = this.state.messages
			ref.on('child_added', function(data) {
				messages.unshift(data.val())
				updateMessages(messages)
			});
		}
	}
	getDomainList() {
		const setDomains = (res) => { this.setState({initialDomains: res}); }
		// Get JSON arry of id:name pairs of domain listing
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
	registerSignIn(user) {
		const welcomeNewUser = () => { $('#welcome-pupal-modal').modal('toggle'); }
		firebase.database().ref('users/' + user.uid).once('value').then(function(snapshot) {
			// If nonexistent, register pupal user on Firebase DB and GAE datastore.
			if (snapshot.val() === null) {
				firebase.database().ref('users/'+user.uid).set({
					name: user.displayName,
					email: user.email,
					photo: user.photoURL,
					summary: "No summary provided"
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
	handleMsgClick(projId, domId, authorUid) {
		if (projId !== "") {
			hashHistory.push("/dom/"+domId+"?view=Projects&proj="+projId)
		} else {
			hashHistory.push("/user/"+authorUid)
		}
	}
	render() {
		return (
			<div className="content home-content">
				<WelcomePupal />
				<div className="domain-search-content">
					<h1 className="text-center">Find domains here.</h1>
					<div className="filtered-list md-form">
						<input type="text" className="form-control" 
							placeholder="Search a domain" 
							onChange={this.filterList} />
						<List domains={this.state.filteredDomains} />
					</div>
				</div>
				<DomainUpdates onMsgClick={(p, d, a)=>this.handleMsgClick(p, d, a)} messages={this.state.messages}/>
			</div>
		);
	}
}
// end home component

// begin list function component
function List(props) {
	function fetchDomLink(id) {
		return "/dom/" + id + "?view=Projects"
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
							Welcome to Pupal! <br/></h4>
					</div>
					<div className="modal-body">
						<p>
							<h4><i>Join a domain!</i></h4><br/>
							Your domain can be your school, group and/or organization.<br/>
							Search for your domain on the right of the page and join!
							<br/><br/>
							<h4><i>Visit your profile!</i></h4><br/>
							Submit a summary about you that other people can read.<br/>
							Remember to add tags that describe your skills and preferences.<br />
						</p>
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

// begin domainupdates component
function DomainUpdates(props) {
	if (props.messages !== null && props.messages.length > 0) {
		return (
			<div className="domain-activity-content">
			{
				props.messages.map((message) => 
					<div className="domain-message-entry" key={message.body}>
						<a onClick={()=>props.onMsgClick(message.projId, message.domId, message.authorUid)}>
							{message.body}
						</a>
						<hr></hr>
					</div>
				)
			}
			</div>
		);
	} else {
		return (
			<div className="domain-activity-content">
				<p><i>No activity in your domains.</i></p>
				<br/>
			</div>
		);
	}
}

// end domainupdates component
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
		const setMember = () => { this.setState({member: true}) }
		var user = firebase.auth().currentUser;
		firebase.database().ref("users/" + user.uid + "/domains").once('value', function(snapshot) {
			snapshot.forEach(function(childSnapshot) {
				if (childSnapshot.val().id === id) {
					setMember()
				}
			})
		})
	}
	handleJoin(id) {
		var user = firebase.auth().currentUser;
		const setJoinState = (res) => { 
			this.setState({member: true})
			// Append to user's domain list in firebase
			var newDomainRef = firebase.database().ref('users/' + user.uid + '/domains').push()
			newDomainRef.set({
				id: res.id,
				name: res.name
			});
		};
		// Send a request for user to join the domain
		user.getToken(true).then(function(token) {
			$.ajax({
				url: "/domain/"+id+"/join",
				type: "POST",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setJoinState(res)
			});
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
	if (props.view === "Projects") {
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
						<li className={this.isActive('Projects')}>
							<Link to={this.fetchViewLink(this.props.id, "Projects")}>Projects</Link></li>
						<li className={this.isActive('Users')}>
							<Link to={this.fetchViewLink(this.props.id, "Users")}>Users</Link></li>
					</ul>
					<ul className="nav navbar-nav navbar-right">
						{!this.props.member && <li><a onClick={()=>this.props.onJoinClick(this.props.id)}>
							<i className="fa fa-tags" aria-hidden="true"></i>
							Request to join</a></li>}
						{this.props.member && <li><Link to={this.fetchViewLink(this.props.id, "Host")}>
							<i className="fa fa-paper-plane" aria-hidden="true"></i>
							Host a project</Link></li>}
					</ul>
				</div>
			</div>
		</nav>
		)
	}
}
// end domain navbar function component

// begin projects component
class Projects extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			projects: [],
			proj: null,
			likeClick: false
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
	sendUpdates(uids, msg, projId, domId, author) {
		for (let uid of uids) {
			var newMsgRef = firebase.database().ref('users/' + uid + '/messages').push()
			newMsgRef.set({
				body: msg,
				projId: projId,
				domId: domId,
				authorUid: author
			});
		}
	}
	handleLikeClick(id) {
		const updateLike = (res) => { 
			this.setState({likeClick: true}) 
			var newLikeRef = firebase.database().ref('users/' + firebase.auth().currentUser.uid + '/likes').push()
			newLikeRef.set({
				projId: id
			});
			this.sendUpdates(res.collab_uids, res.msg, "", "", firebase.auth().currentUser.uid)
		}
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "/projects/"+id+"/like",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => updateLike(res)
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
					likeClick={this.state.likeClick}
					onProjLikeClick={()=>this.handleLikeClick(this.state.proj.id)}
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
										<i>{proj.likes} like(s)</i></p>
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
							{
							!props.proj.is_collaborator && 
								<button type="button" className="btn btn-default btn-circle btn-lg"
									onClick={()=>props.onProjLikeClick()}>
									{
									(!props.proj.has_liked && !props.likeClick) ? 
										<i className="fa fa-thumbs-o-up" aria-hidden="true"></i> 
										: <i className="fa fa-thumbs-up" aria-hidden="true"></i> 
									}
								</button>
							}
							{
							props.proj.is_collaborator && 
								<button type="button" className="btn btn-primary btn-circle btn-lg">
									<i className="fa fa-pencil-square-o" aria-hidden="true"></i>
								</button>
							}
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
							<div className="collaborator-content">
								{
								props.proj.collaborators.map((collab) =>
									<div className="collab-entry" key={collab.uid}>
									<a onClick={()=>props.onUserClick(collab.uid)} 
										data-dismiss="modal" >
										<img className="proj-collab-image img-fluid" 
											src={collab.photo} alt={collab.name}></img>
										<div className="collab-info-contact">
											<h4>{collab.name}</h4>
										</div>
									</a>
									<br/><br/>
									</div>)
								}
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
			<h2>Coming soon!</h2>
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
	sendUpdates(uids, msg, projId, domId, author) {
		for (let uid of uids) {
			var newMsgRef = firebase.database().ref('users/' + uid + '/messages').push()
			newMsgRef.set({
				body: msg,
				projId: projId,
				domId: domId,
				authorUid: author
			});
		}
	}
	handleSubmitClick(id, titl, desc, ts, web) {
		const setProject = (res) => {
			$('#failure-alert').hide();
			$('#success-alert').show();
			this.setState({projId: res.proj_id})
			this.sendUpdates(res.updated_uids, res.msg, res.proj_id, id, firebase.auth().currentUser.uid)
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
				<hr></hr>
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
				<hr></hr>
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
				<hr></hr>
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
			user: null,
			domains: [],
			ownership: false,
		};
	}
	componentWillMount() {
		this.setOwnership()
	}
	setOwnership() {
		if (this.props.params.id === firebase.auth().currentUser.uid) {
			this.setState({ownership: true})
		}
	}
	componentDidMount() {
		this.getUser(this.props.params.id)
	}
	getUser(id) {
		const updateUser = (user) => { 
			this.setState({user: user}) 
		}
		const updateDomains = (domains) => { this.setState({domains: domains}) }

		firebase.database().ref("users/"+this.props.params.id).once('value', function(snapshot) {
			updateUser(snapshot.val())
		});

		let domains = this.state.domains
		firebase.database().ref("users/"+this.props.params.id+"/domains").once('value', function(snapshot) {
			snapshot.forEach(function(childSnapshot) {
				domains.push(childSnapshot.val().name)
			})
			updateDomains(domains)
		});
	}
	handleSend(msg) {
		var newMsgRef = firebase.database().ref('users/' + this.props.params.id + '/messages').push()
		msg = firebase.auth().currentUser.displayName + ": " + msg
		newMsgRef.set({
			body: msg,
			projId: "",
			domId: "",
			authorUid: firebase.auth().currentUser.uid
		});
		$('#message-sent-popover').popover('toggle')
	}
	sendUpdates(uids, msg, projId, domId, author) {
		for (let uid of uids) {
			var newMsgRef = firebase.database().ref('users/' + uid + '/messages').push()
			newMsgRef.set({
				body: msg,
				projId: projId,
				domId: domId,
				authorUid: author
			});
		}
	}
	handleProjClick(id, uid) {
		const setUpdates = (res) => {
			$('#new-collab-popover').popover('toggle')
			this.sendUpdates(res.uids, res.msg, res.proj_id, res.dom_id, firebase.auth().currentUser.uid)
		}
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "projects/"+id+"/newCollab",
				type: "POST",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				data: {uid: uid},
				success: (res) => setUpdates(res)
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
						<hr></hr>
						<h4 id="user-profile-summary">{this.state.user.summary}</h4>
						<hr></hr>
						{
						(this.state.domains.length > 0) && 
							<div className="user-domain-content">
								<h4>Domains joined</h4>
								<div className="user-domains">
								{
									this.state.domains.map((dom) =>
									<div className="user-dom" key={dom}>
										<p><i className="fa fa-university" aria-hidden="true"></i>{dom}</p>
									</div>
									)
								}
								</div>
							</div>
						}
					</div>
					{ 
					!this.state.ownership && 
						<UserMessageForm onSend={(body)=>{this.handleSend(body)}} />
					}
					{
					!this.state.ownership && 
						<CollaborationInvite onProjClick={(id)=>this.handleProjClick(id, this.props.params.id)} />
					}
				</div>
			)
		}
		return (
			<div className="content user-content">
				Loading...
			</div>
		)
	}
}
// end user component

// begin messageform component
class UserMessageForm extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			body: "",
			valid: false
		};
		this.handleMsgChange = this.handleMsgChange.bind(this)
	}
	handleMsgChange(event) {
		this.setState({body: event.target.value});
		if (event.target.value.length > 0) {
			this.setState({valid: true})
		} else {
			this.setState({valid: false})
		}
	}
	checkMessageClick() {
		if (this.state.valid) {
			this.props.onSend(this.state.body)
		}
	}
	render() {
		return (
			<div className="message-form-content">
				<textarea value={this.state.body} 
					onChange={this.handleMsgChange} 
					className="form-control" id="user-message-form-textarea" 
					rows="5" 
					placeholder="Enter your message here.">
				</textarea>
				<div id="submit-buttons">
					<a onClick={()=>this.checkMessageClick()} type="button" tabIndex="0" className="btn btn-success"
						id="message-sent-popover" data-toggle="popover" data-trigger="focus" data-content="Message sent!">
						<i className="fa fa-check" aria-hidden="true"></i>Send message</a>
				</div>
			</div>
		);
	}
}
// end messageform component

// begin collaborationinvite component
class CollaborationInvite extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			projects: [],
			projId: null
		};
		this.handleProjId = this.handleProjId.bind(this)
	}
	componentDidMount() {
		const setProjects = (res) => { this.setState({projects: res})}
		firebase.auth().currentUser.getToken(true).then(function(token) {
			$.ajax({
				url: "users/projects",
				type: "GET",
				beforeSend: function(xhr){
					xhr.setRequestHeader('Authorization', token);
				},
				success: (res) => setProjects(res)
			});
		});
	}
	handleProjId(event) {
		this.setState({projId: event.target.value})
	}
	render() {
		if (this.state.projects.length > 0) {
			return (
			<div className="collaboration-invite-content">
				<hr></hr>
				<h3>Add as collaborator</h3>
				<select value={this.state.projId} 
					onChange={this.handleProjId}>
					<option value=""></option>
					{
					this.state.projects.map((proj) => 
						<option key={proj.id} value={proj.id}>{proj.name}</option>
					)
					}
				</select>
				{this.state.projId !== null && this.state.projId.length > 0 && 
				<a onClick={()=>this.props.onProjClick(this.state.projId)} type="button" tabIndex="0" className="btn btn-success"
					id="new-collab-popover" data-toggle="popover" data-trigger="focus" data-content="Added new collaboration!">
					<i className="fa fa-handshake-o" aria-hidden="true"></i>Confirm</a>}
			</div>
			);
		}
		return null
	}
}
// end collaborationinvite component


// begin profile component
function Profile()  {
	return (
		<div className="content profile-content">
			Coming soon!
		</div>
	);
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

// Query options -> dom/:id?view=Projects dom/:id?view=Projects&proj=<id> dom/:id?view=Users dom/:id?view=Host
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
