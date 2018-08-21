import React, { Component } from 'react';
import { Route, Switch, Redirect } from 'react-router-dom'
// import { StripeProvider } from 'react-stripe-elements';
import { AuthenticateService } from './services/AuthenticateService'
// import logo from './logo.svg';
import './App.css';

import IndexPage from './components/indexPage/IndexPage'
import Login from './components/Login'
import Register from './components/Register'
import ForgotPassword from './components/ForgotPassword'
import ResetPassword from './components/ResetPassword'
import ErrorPage from './components/ErrorPage';
import LoginOTP from './components/LoginOTP';

class App extends Component {
	constructor() {
		super();
		this.state = {
			isAuth: AuthenticateService.isAuthenticate(),
		}
	}
	render() {
		const { isAuth } = this.state
		return (
			<React.Fragment>
			{/* <StripeProvider apiKey="pk_test_3KZy6O2l5hseQy9hglRwswYT"> */}
			<div>
			{
				!isAuth &&
				
					<div className="app">
						<Switch>
							<Route exact path='/login' component={Login} />
							<Route exact path='/register' component={Register} />
							<Route exact path='/forgot-password' component={ForgotPassword} />
							<Route exact path='/otp-verify' component={LoginOTP} />
							<Route exact path='/reset-password/:token' component={ResetPassword} />
							<Route exact path='/' component={IndexPage} />
							<Route exact path='/404' component={ErrorPage} />
							<Route render={() => (<Redirect to="/login"/>)}/>
						</Switch>
					</div>
				// </StripeProvider>
			}
			{
				isAuth &&
				// <StripeProvider apiKey="pk_test_3KZy6O2l5hseQy9hglRwswYT">
					<div className="app">
						<Route exact path='/404' component={ErrorPage} />
						<Route path='/' component={IndexPage} />
					</div>
				
			}
			</div>
			{/* </StripeProvider> */}
			</React.Fragment>
		
		);
	}
}

export default App;
