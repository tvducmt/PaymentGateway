import React, {Component} from 'react';
import { sendRequest } from '../services/Http.service'

class ForgotPassword extends Component {
    constructor() {
		super();
		this.state = {
			email: '',
			formErrors: {email: ''},
			errForgotPw: false,
			sending: false,
			sendsuccess: false
		}
	}
	onChangeEmail = (e) => {
		this.setState({ email: e.target.value })
		this.validateField('email', e.target.value)
	}
	validateField(fieldName, value) {
		let fieldValidationErrors = this.state.formErrors;
		let emailValid = this.state.email;
	  
		switch(fieldName) {
			case 'email':
				emailValid = value.match(/^([\w.%+-]+)@([\w-]+\.)+([\w]{2,})$/i);
				fieldValidationErrors.email = emailValid ? '' : 'Địa chi email không hợp lệ';
				break;
			default:
				break;
		}
	}
	renderErrorLogin() {
		if (this.state.errForgotPw) {
			return (
				<div className="alert alert-danger">
					<button type="button" className="close" data-dismiss="alert" aria-hidden="true">&times;</button>
					<strong>Error!</strong> {this.state.msg}
				</div>); 
		} else {
			return;
		}
	}
	onSubmit = (e) => {
		e.preventDefault();
		const { email, formErrors } = this.state;
		if (formErrors.email === '') {
			let data = { email };
			this.setState({sending: true});
			sendRequest('post', '/user/forgot-password', data).then((res) => {
				console.log('onSubmit: ', res)
				if (!res.isError) {
					if (res.data.error === "") {
						this.setState({errForgotPw: false});
						this.setState({msg: ''});
						this.setState({sendsuccess: true});
						// console.log(this.props.history)
					} else {
						this.setState({errForgotPw: true})
						this.setState({msg: res.data.error})
					}
					this.setState({sending: false});
				} else {
					this.setState({sending: false});
					this.setState({errForgotPw: true});
					this.setState({msg: "Something went wrong!"})
				}
			})
		}
	}
    render() {
		const {
            email,
			formErrors,
			sendsuccess
		} = this.state
		let errEmail;
		if (formErrors.email) {
			errEmail = <label id="email-error" className="error">{formErrors.email}</label>;
		}   
		return (
			<div id="page-login">
				<div className="container">
					<div className="card card-container">
						<img id="profile-img" alt="" className="profile-img-card" src="//ssl.gstatic.com/accounts/ui/avatar_2x.png" />
						{	
							!sendsuccess &&
							<React.Fragment>
								{this.renderErrorLogin()}
								<form onSubmit={this.onSubmit} id="form-signin" className="form-signin">
									<span id="reauth-email" className="reauth-email"></span>
									<div className="form-group"><h4>Quên mật khẩu?</h4><h5>Nhập địa chỉ email của bạn để khôi phục lại mật khẩu</h5></div>
									<div className="form-group">
										<input onChange={this.onChangeEmail} value={email}  type="email" name="email" id="email" className="form-control" placeholder="Email address" required />
										{errEmail}
									</div>
									<button className={this.state.sending? 'btn btn-lg btn-primary btn-block btn-signin sending' : 'btn btn-lg btn-primary btn-block btn-signin'} type="submit">
										<div className="loader"></div>
										Gửi</button>
								</form>
							</React.Fragment>
						}
						{
							sendsuccess &&
							<div className="forgot-password-result">
								<h3>Kiểm tra email của bạn</h3>
								<h5>Email {email} của bạn sẽ nhận được hướng dẫn làm thế nào để khôi phục mật khẩu.</h5>
							</div>
						}
						<a href="/login" className="register">Đăng nhập</a>
					</div>
				</div>
			</div>
       );
    }
 }
 export default ForgotPassword;
 