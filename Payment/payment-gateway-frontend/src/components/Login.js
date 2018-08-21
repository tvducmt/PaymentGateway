import React from 'react';
import { sendRequest } from '../services/Http.service'
import { AuthenticateService } from '../services/AuthenticateService'

class Login extends React.Component {
	constructor() {
		super();
		this.state = {
			email: '',
			password: '',
			formErrors: {email: '', password: ''},
			errLogin: false,
			msg: '',
			sending: false
		}
	}
	
	onChangeEmail = (e) => {
		this.setState({ email: e.target.value })
		this.validateField('email', e.target.value)
	}
	onChangePassword = (e) => {
		this.setState({ password: e.target.value })
		this.validateField('password', e.target.value)
	}
	validateField(fieldName, value) {
		let fieldValidationErrors = this.state.formErrors;
		let emailValid = this.state.email;
		let passwordValid = this.state.password;
	  
		switch(fieldName) {
			case 'email':
				emailValid = value.match(/^([\w.%+-]+)@([\w-]+\.)+([\w]{2,})$/i);
				fieldValidationErrors.email = emailValid ? '' : 'Địa chi email không hợp lệ';
				break;
			case 'password':
				passwordValid = value.length >= 6;
				fieldValidationErrors.password = passwordValid ? '': 'Mật khẩu ít nhất 6 ký tự';
				break;
			default:
				break;
		}
	}
	onSubmit = (e) => {
		e.preventDefault();
		const { email, password, formErrors } = this.state;
		if (formErrors.email === '' && formErrors.password === '') {
			let data = { email, password };
			// console.log(data);
			this.setState({sending: true});
			sendRequest('post', '/login', data).then((res) => {
				// console.log('onSubmit: ', res)
				if (!res.isError) {
					if (res.data.status === "OK") {
						this.setState({errLogin: false});
						this.setState({msg: ''});	
						if(!res.data.otpenable) {
							AuthenticateService.setAuthenticateUser(res.data.token, res.data.result)
						} else {
							this.props.history.push({
								pathname: '/otp-verify',
								state: { data: res.data }
							})
						}
						// this.props.history.push("/")
					} else {
						this.setState({errLogin: true})
						this.setState({sending: false});
						this.setState({msg: res.data.error})
					}
				} else {
					this.setState({sending: false});
					this.setState({errLogin: true});
					this.setState({msg: "Something went wrong!"})
				}
			})
		}
	}
	renderErrorLogin() {
		if (this.state.errLogin) {
			return (
				<div className="alert alert-danger">
					<button type="button" className="close" data-dismiss="alert" aria-hidden="true">&times;</button>
					<strong>Error!</strong> {this.state.msg}
				</div>); 
		} else {
			return;
		}
	}
	render() {
		const {
            email,
			password,
			formErrors
		} = this.state
		let errEmail, errPw;
		if (formErrors.email) {
			errEmail = <label id="email-error" className="error">{formErrors.email}</label>;
		}
		if (formErrors.password) {
			errPw = <label id="password-error" className="error">{formErrors.password}</label>;
		}
		return (
			<div id="page-login">
				<div className="container">
					<div className="card card-container">
						<img id="profile-img" alt="" className="profile-img-card" src="//ssl.gstatic.com/accounts/ui/avatar_2x.png" />
						<p id="profile-name" className="profile-name-card"></p>
						{this.renderErrorLogin()}
						<form onSubmit={this.onSubmit} id="form-signin" className="form-signin">
							<span id="reauth-email" className="reauth-email"></span>
							<div className="form-group">
								<input onChange={this.onChangeEmail} value={email}  type="email" name="email" id="email" className="form-control" placeholder="Email address" required />
								{errEmail}
							</div>
							<div className="form-group">
								<input onChange={this.onChangePassword} value={password} type="password" name="password" id="password" className="form-control" placeholder="Password" required />
								{errPw}
							</div>
							<button className={this.state.sending? 'btn btn-lg btn-primary btn-block btn-signin sending' : 'btn btn-lg btn-primary btn-block btn-signin'} type="submit">
								<div className="loader"></div>
								Đăng nhập</button>
						</form>
						<div><a href="/register" className="register">Đăng ký một tài khoản!</a></div>
						<div><a href="/forgot-password" className="register">Quên mật khẩu</a></div>
					</div>
				</div>
			</div>
		);
	 }
}
export default Login;