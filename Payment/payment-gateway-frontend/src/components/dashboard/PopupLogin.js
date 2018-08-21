import React, {Component} from 'react';

class PopupLogin extends Component {
    constructor() {
        super()
        this.state = {
            email: '',
			password: '',
			formErrors: {email: '', password: ''},
        }
    }
    onChangePassword = (e) => {
		this.setState({ password: e.target.value })
		this.validateField('password', e.target.value)
    }
    onSubmit = (e) => {
		e.preventDefault();
		const { email, password, formErrors } = this.state;
		if (formErrors.email === '' && formErrors.password === '') {
			let data = { email, password };
			// console.log(data);
			this.setState({sending: true});
			sendRequest('post', '/login', data).then((res) => {
				console.log('onSubmit: ', res)
				if (!res.isError) {
					if (res.data.status === "OK") {
						this.setState({errLogin: false});
						this.setState({msg: ''});
						// console.log(this.props.history)
						AuthenticateService.setAuthenticateUser(res.data.token, res.data.result)
						// this.props.history.push("/")
					} else {
						this.setState({errLogin: true})
						this.setState({msg: res.data.error})
					}
					this.setState({sending: false});
				} else {
					this.setState({sending: false});
					this.setState({errLogin: true});
					this.setState({msg: "Something went wrong!"})
				}
			})
		}
	}
    render() {
       return (
            <div className="check-login">
                <form onSubmit={this.onSubmit} id="form-signin" className="form-signin">
                    <span id="reauth-email" className="reauth-email"></span>
                    {/* <div className="form-group">
                        <input onChange={this.onChangeEmail} value={email}  type="email" name="email" id="email" className="form-control" placeholder="Email address" required />
                        {errEmail}
                    </div> */}
                    <div className="form-group">
                        <input onChange={this.onChangePassword} value={password} type="password" name="password" id="password" className="form-control" placeholder="Password" required />
                        {errPw}
                    </div>
                    <button className={this.state.sending? 'btn btn-lg btn-primary btn-block btn-signin sending' : 'btn btn-lg btn-primary btn-block btn-signin'} type="submit">
                        <div className="loader"></div>
                        Xác minh</button>
                </form>
            </div>
       );
    }
 }
 export default PopupLogin;
 