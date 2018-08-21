import React from 'react';
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from '../HomePage'
import Transaction from '../Transaction'
import Dashboard from '../Dashboard'

import ErrorPage from '../ErrorPage'
import { AuthenticateService } from '../../services/AuthenticateService'

class IndexPage extends React.Component {
	constructor() {
        super()
        this.state = {
			isAuth: AuthenticateService.isAuthenticate(),
		}
    }
	render() {
		const { isAuth } = this.state;
		return (
			<div>
				{	
					!isAuth &&
					<div id="index-page">
						<Switch>
							{/* <Route path='/transaction' component={Transaction} /> */}
							<Route exact path='/' component={HomePage} />
							<Route render={() => (<Redirect to="/login"/>)}/>
							{/* <Route component={ErrorPage} /> */}
						</Switch>
					</div>
				}
				{
					isAuth &&
					<React.Fragment>
						<Switch>
							<Route path='/dashboard' component={Dashboard} />
							<Route path='/transaction' component={Transaction} />
							<Route exact path='/login' render={() => (<Redirect to="/"/>)}/>
							<Route exact path='/register' render={() => (<Redirect to="/"/>)}/>
							<Route exact path='/' component={HomePage} />
							<Route component={ErrorPage} />
						</Switch>
					</React.Fragment>
				}
			</div>
			
		);
	 }
}
export default IndexPage;