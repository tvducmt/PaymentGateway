import React, {Component} from 'react';
import { AuthenticateService } from '../services/AuthenticateService'
import { sendRequest } from '../services/Http.service'
import swal from 'sweetalert2'

class LoginOTP extends Component {
    constructor() {
        super();
        // console.log(this.props.location.state.data)
        this.state = {
            isAuth: AuthenticateService.isAuthenticate(),
            codeOTP: '',
            formErrors: '',
            sending: false,
            sendsuccess: false,
            data: null
        }    
    }
    componentWillMount() {
        if(!this.state.isAuth && typeof(this.props.location.state) !== 'undefined') {
            this.setState({ data: this.props.location.state.data })
        } else {
            this.props.history.push("/")
        }
    }
	onChangeOTP = (e) => {
		this.setState({ codeOTP: e.target.value })
		this.validateField('codeOTP', e.target.value)
    }
    validateField(fieldName, value) {
		let fieldValidationErrors = this.state.formErrors;
		let codeOTPValid = this.state.codeOTP;
	  
		switch(fieldName) {
			case 'codeOTP':
                codeOTPValid = value.match(/^([\d]{6,6})$/i);
                fieldValidationErrors = codeOTPValid ? '' : 'Mã OTP không hợp lệ';
                this.setState({formErrors: fieldValidationErrors})
				break;
			default:
				break;
		}
    }
    onSubmit = (e) => {
        e.preventDefault();
        if(this.state.data != null) {
            let data = {
                user_id: this.state.data.result.id,
                otp_pin: this.state.codeOTP,
                accesstoken_encrypted: this.state.data.token
            }
            this.setState({sending: true});
            sendRequest('put', '/user/otp', data).then((res) => {
                // console.log('onSubmit OTP: ', res)
				if (!res.isError) {
                    if(typeof(res.data.error) !== 'undefined') {
                        swal('Error!', `${res.data.error}`, 'error')
                        this.setState({sending: false});
                    }
                    else if(res.data.status === 200) {
                        AuthenticateService.setAuthenticateUser(res.data.token, res.data.result)
                    }
                } else {
                    swal('Error!', `Có lỗi xảy ra, vui lòng thử lại`, 'error')
					this.setState({sending: false});
				}
            });
        }
    }
    componentDidMount() {
    }
    render() {
        const { codeOTP, formErrors } = this.state
        let errOTP;
		if (formErrors) {
			errOTP = <label id="otp-error" className="error">{formErrors}</label>;
		}  
        return (
			<div id="page-login">
				<div className="container">
					<div className="card card-container">
						<img id="profile-img" alt="" className="profile-img-card" src="//ssl.gstatic.com/accounts/ui/avatar_2x.png" />

                        <form onSubmit={this.onSubmit} id="form-signin" className="form-signin">
                            <span id="reauth-email" className="reauth-email"></span>
                            <div className="form-group">
                                <h4>Nhập mã OTP!</h4>
                                {/* <h5>Kiểm tra gmail của bạn</h5> */}
                            </div>
                            <div className="form-group">
                                <input onChange={this.onChangeOTP} value={codeOTP}  type="number" className="form-control" placeholder="Nhập OTP" required />
                                {errOTP}
                            </div>
                            <button className={this.state.sending? 'btn btn-lg btn-primary btn-block btn-signin sending' : 'btn btn-lg btn-primary btn-block btn-signin'} type="submit">
                                <div className="loader"></div>
                                Gửi</button>
                        </form>
					</div>
				</div>
			</div>
       );
    }
 }
 export default LoginOTP;
 