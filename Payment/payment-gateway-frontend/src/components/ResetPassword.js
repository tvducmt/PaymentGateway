import React, {Component} from 'react';
import { sendRequest } from '../services/Http.service'
import swal from 'sweetalert2'

class ResetPassword extends Component {
    constructor() {
		super();
		this.state = {
            newPw: '',
            confirmNewPw: '',
            formErrors: {newPw: '', confirmNewPw: ''},
            sending: false
        }
    }
    onChangePassword = (e) => {
		this.setState({ newPw: e.target.value })
		this.validateField('password', e.target.value)
    }
    onChangeConfirmPassword = (e) => {
		this.setState({ confirmNewPw: e.target.value })
		this.validateField('confirmPw', e.target.value)
    }
    validateField(fieldName, value) {
		let fieldValidationErrors = this.state.formErrors;
        let newPwValid = this.state.newPw;
        let confirmNewPwValid = this.state.confirmNewPw;
	  
		switch(fieldName) {
			case 'password':
				newPwValid = value.length >= 6;
				fieldValidationErrors.newPw = newPwValid ? '': 'Mật khẩu ít nhất 6 ký tự';
                break;
            case 'confirmPw':
                confirmNewPwValid = (value === newPwValid);
				fieldValidationErrors.confirmNewPw = confirmNewPwValid ? '': 'Mật khẩu xác nhận không khớp';
				break;
			default:
				break;
		}
    }
    onSubmit = (e) => {
        e.preventDefault();
        
        const { newPw, confirmNewPw, formErrors } = this.state;
        if (formErrors.newPw === '' && formErrors.confirmNewPw === '') {
            let data = {
                new_pw: newPw,
                confirm_new_pw: confirmNewPw,
                token: this.props.match.params.token
            }
            this.setState({sending: true})
            sendRequest('post', '/user/reset-password', data).then((res) => {
                if (!res.isError) {
                    if( typeof(res.data.error) !== "undefined") {
                        swal('Error!', `${res.data.error}`, 'error')
                    } else if ( typeof(res.data.message) !== "undefined" ) {
                        swal('Success!', `${res.data.message}`, 'success')
                        let reactEl = this
                        setTimeout(function(){ reactEl.props.history.push("/login") }, 1000);
                    }
                } else {
                    swal('Error!', `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                }
                this.setState({sending: false})
            })
        }
    }
    render() {
        const {
            newPw,
            confirmNewPw,
            formErrors
        } = this.state
        let errNewPw, errConfirmNewPw;
		if (formErrors.newPw) {
			errNewPw = <label id="email-error" className="error">{formErrors.newPw}</label>;
        }
        if (formErrors.confirmNewPw) {
			errConfirmNewPw = <label id="email-error" className="error">{formErrors.confirmNewPw}</label>;
		}  
        return (
			<div id="page-login">
				<div className="container">
					<div className="card card-container">
						<img id="profile-img" alt="" className="profile-img-card" src="//ssl.gstatic.com/accounts/ui/avatar_2x.png" />
						{/* {this.renderErrorLogin()} */}
						<form onSubmit={this.onSubmit} id="form-signin" className="form-signin">
							<span id="reauth-email" className="reauth-email"></span>
							<div className="form-group">
								<input onChange={this.onChangePassword} value={newPw}  type="password" className="form-control" placeholder="Nhập mật khẩu mới" required />
								{errNewPw}
							</div>
							<div className="form-group">
								<input onChange={this.onChangeConfirmPassword} value={confirmNewPw} type="password" className="form-control" placeholder="Xác nhận mật khẩu mới" required />
								{errConfirmNewPw}
							</div>
							<button className={this.state.sending? 'btn btn-lg btn-primary btn-block btn-signin sending' : 'btn btn-lg btn-primary btn-block btn-signin'} type="submit">
								<div className="loader"></div>
								Thay doi</button>
						</form>
						<div><a href="/login" className="register">Đăng nhập</a></div>
					</div>
				</div>
			</div>
       );
    }
 }
 export default ResetPassword;
 