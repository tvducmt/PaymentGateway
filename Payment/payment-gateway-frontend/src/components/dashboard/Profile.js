import React, {Component} from 'react';
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { sendRequest } from '../../services/Http.service'
import swal from 'sweetalert2'

class Profile extends Component {
    constructor() {
        super()
        this.state = {
            userInfo: {
                userid: '',
                email: '',
                phone: '',
                fullname: '',
                otpenable: false,
                otpsecretkey: ''
            },
            saveUserInfo: {
                userid: '',
                email: '',
                phone: '',
                fullname: '',
                otpenable: false,
                otpsecretkey: ''
            },
            isAuth: AuthenticateService.isAuthenticate(),
            isLoadInfo: false,
            isLoadChangePw: false,
            userPW: {
                oldPw: '',
                newPw: '',
                confirmNewPw: ''
            },
            isLoad: true,
            isLoadOTP: false,
            srcImageQRCode: ''
        }
    }
    onChangeName = (e) => {
        this.setState({ userInfo: { ...this.state.userInfo, fullname: e.target.value } })
    }
    onChangePhone = (e) => {
		this.setState({ userInfo: { ...this.state.userInfo, phone: e.target.value } })
    }
    onChangeOldPw = (e) => {
        this.setState({ userPW: { ...this.state.userPW, oldPw: e.target.value } })
    }
    onChangeNewPw = (e) => {
		this.setState({ userPW: { ...this.state.userPW, newPw: e.target.value } })
    }
    onChangeConfirmPw = (e) => {
		this.setState({ userPW: { ...this.state.userPW, confirmNewPw: e.target.value } })
    }
    onChangeSecurity2FA = (e) => {
        if (!this.state.userInfo.otpenable) {
            //Call API
            if(this.state.isAuth) {
                this.setState({ isLoadOTP: true })
                let data = this.state.userInfo
                data.otpenable = true;
                let headers = {
                    'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
                }
                sendRequest('put', '/user/info', data, headers).then((res) => {
                    console.log(res)
                    if(!res.isError) {
                        if(typeof(res.data.error) !== 'undefined') {
                            swal('Error!', `${res.data.error}`, 'error')
                            this.setState({userInfo: { ...this.state.userInfo, otpenable: false } });
                        } else if (typeof(res.data.urlQRCode) !== 'undefined') {
                            swal('Success!', `Bật chức năng xác minh 2 lớp thành công!`, 'success')
                            this.setState({ 
                                srcImageQRCode: `https://chart.googleapis.com/chart?chs=166x166&chld=L%7C0&cht=qr&chl=${res.data.urlQRCode}` 
                            })
                            this.setState({userInfo: { ...this.state.userInfo, otpenable: true } });
                        } else {
                            swal("Error!", `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                            this.setState({userInfo: { ...this.state.userInfo, otpenable: false } });
                        }
                    } else {
                        swal("Error!", `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                        this.setState({userInfo: { ...this.state.userInfo, otpenable: false } });
                    }
                    this.setState({ isLoadOTP: false })
                    
                    
                }).catch((err) => {
                    console.log(err)
                });
            }

        } else {
            if(this.state.isAuth) {

                let data = this.state.userInfo
                data.otpenable = false
                let headers = {
                    'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
                }
                this.setState({ isLoadOTP: true }) 
                sendRequest('put', '/user/info', data, headers).then((res) => {
                    if(!res.isError) {
                        if(typeof(res.data.error) !== 'undefined') {
                            swal('Error!', `${res.data.error}`, 'error')
                        } else {
                            swal('Success!', `Xác minh hai bước đã được vô hiệu hóa!`, 'success')
                        }
                    } else {
                        swal("Error!", `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                    }
                    this.setState({ isLoadOTP: false })
                    this.setState({userInfo: { ...this.state.userInfo, otpenable: false } });      
                }).catch((err) => {
                    console.log(err)
                });
            }
        }
    }
    inforChangeByUser = () => {
        const {userInfo, saveUserInfo} = this.state
        let data = {}
        if (userInfo.fullname !== saveUserInfo.fullname) {
            data.fullname = userInfo.fullname
        }
        if (userInfo.phone !== saveUserInfo.phone) {
            data.phone = userInfo.phone
        }
        return data
    }
    onSubmitInforUser = (e) => {
        e.preventDefault();
        if(this.state.isAuth) {
            this.setState({ isLoadInfo: true })
            let data = this.state.userInfo
            let headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            console.log(data)   
            sendRequest('put', '/user/info', data, headers).then((res) => {
                if(!res.isError) {
                    if(typeof(res.data.error) !== 'undefined') {
                        swal('Error!', `${res.data.error}`, 'error')
                    } else {
                        swal('Success!', `Update thông tin thành công`, 'success')
                    }
                } else {
                    swal("Error!", `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                }
                this.setState({ isLoadInfo: false })
                
            }).catch((err) => {
                console.log(err)
            });
        }
    }
    onSubmitChangePw = (e) => {
        e.preventDefault();
        if (this.state.userPW.newPw !== this.state.userPW.confirmNewPw) {
            swal('Error!', 'Nhập lại mật khẩu mới không khớp. Vui lòng thử lại', 'error')
            return
        }
        if (this.state.userPW.oldPw === this.state.userPW.newPw) {
            swal('Error!', 'Mật khẩu mới bị trùng. Vui lòng thử lại', 'error')
            return
        }
        if(this.state.isAuth) {
            this.setState({ isLoadChangePw: true })
            let data = {
                userid: AuthenticateService.getAuthenticateUser().id,
                oldpassword: this.state.userPW.oldPw,
                newpassword: this.state.userPW.newPw
            }
            let headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            sendRequest('put', '/user/password', data, headers).then((res) => {
                if(!res.isError) {
                    if(typeof(res.data.error) !== 'undefined') {
                        swal('Error!', `${res.data.error}`, 'error')
                    } else {
                        this.setState({
                            userPW: {
                                oldPw: '',
                                newPw: '',
                                confirmNewPw: ''
                            }
                        });
                        swal('Sucess!', 'Thay đổi mật khẩu thành công', 'success')
                    }
                } else {
                    swal('Error!', 'Có lỗi xảy ra, vui lòng thử lại', 'error')
                }
            })
        }
    }
    componentWillMount() {
        if(this.state.isAuth) {
            let headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            sendRequest('get', '/user/info', null,headers).then((res) => {
                console.log(res)
                if(!res.isError) {
                    if(typeof(res.data.error) !== 'undefined') {
                        swal('Error!', `${res.data.error}`, 'error')
                    } else {
                        this.setState({userInfo: res.data})
                        this.setState({saveUserInfo: res.data})
                        console.log(res.data)
                    }
                } else {
                    swal('Error!', 'Something went wrong. Please try again!', 'error')
                }
                this.setState({ isLoad: false })
            }).catch((err) => {
                console.log(err)
                this.setState({ isLoad: false })
            });
        }
    }
    formatMoney(money) {
        money = money.toLocaleString('it-IT', {style : 'currency', currency : 'VND'});
        return money;
    }
    render() {
        const {userInfo, isLoadInfo, isLoadChangePw, userPW, isLoad, isLoadOTP, srcImageQRCode} = this.state
        return (
            <div className="content-wrapper">
                <div className="wrapper-table"> 
                    <div className="table-header"><i className="icons glyphicon glyphicon-user"></i><span>Hồ sơ cá nhân</span></div>
                    <div className="table-body">
                        <div className="row">
                            <div className="col-xs-6">
                            <form onSubmit={this.onSubmitInforUser}>
                                <div className="table-responsive table-profile">
                                    <table className="table table-bordered ">
                                        <tbody>
                                            <tr>
                                                <td>ID</td>
                                                <td>User #{userInfo.userid}</td>
                                            </tr>
                                            <tr>
                                                <td>Full name</td>
                                                <td>
                                                    <input type="text" value={userInfo.fullname} onChange={this.onChangeName} className="form-control" placeholder="Enter name ..."/>
                                                </td>
                                            </tr>
                                            <tr>
                                                <td>Số điện thoại</td>
                                                <td><input type="text" value={userInfo.phone} onChange={this.onChangePhone} className="form-control" placeholder="Enter phonenumber ..."/></td>
                                            </tr>
                                            <tr>
                                                <td>Email</td>
                                                <td>{userInfo.email}</td>
                                            </tr>
                                            <tr>
                                                <td colSpan="2">
                                                    <button type="submit" disabled={isLoadInfo} className="btn btn-primary">Cập nhật thông tin</button>
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </form>
                            </div>
                            <div className="col-xs-6">
                                <div className="table-responsive table-profile">
                                <form onSubmit={this.onSubmitChangePw}>
                                    <table className="table table-bordered ">
                                        <tbody>
                                            <tr>
                                                <td colSpan="2">Thay đổi mật khẩu</td>
                                                {/* <td>User #</td> */}
                                            </tr>
                                            <tr>
                                                <td>Mật khẩu hiện tại</td>
                                                <td><input type="password" value={userPW.oldPw} onChange={this.onChangeOldPw} className="form-control" placeholder="Old password ..."/></td>
                                            </tr>
                                            <tr>
                                                <td>Mật khẩu mới</td>
                                                <td><input type="password" value={userPW.newPw} onChange={this.onChangeNewPw} className="form-control" placeholder="New password ..."/></td>
                                            </tr>
                                            <tr>
                                                <td>Nhập lại mật khẩu</td>
                                                <td><input type="password" value={userPW.confirmNewPw} onChange={this.onChangeConfirmPw} className="form-control" placeholder="Confirm new password ..."/></td>
                                            </tr>
                                            <tr>
                                                <td colSpan="2">
                                                    <button type="submit" disabled={isLoadChangePw} className="btn btn-primary">Thay đổi mật khẩu</button>
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </form>
                                </div>
                            </div>
                        </div>        
                    </div>
                    <div className="table-header">
                        <i className="icons glyphicon glyphicon-eye-open"></i>
                        <span>Xác minh 2 bước: { userInfo.otpenable ? "BẬT" : "TẮT" }</span>
                    </div>
                    <div className="table-body">
                        <div className="security-2fa">
                            <h4>Bảo vệ tài khoản của bạn bằng Xác minh 2 bước</h4>
                            <label className="custom-control custom-checkbox">
                                <input type="checkbox" className="custom-control-input"
                                value={userInfo.otpenable}
                                checked={userInfo.otpenable === true}
                                onChange={this.onChangeSecurity2FA} />
                                <span className="custom-control-indicator"></span>
                            </label> 
                        </div>
                        <h4>Thiết lập Authenticator</h4>
                        <ul>
                            <li>Tải Ứng dụng Authenticator từ 
                                <a target="_blank" rel="noopener noreferrer" href="https://itunes.apple.com/us/app/google-authenticator/id388497605"> App Store</a> hoặc
                                <a target="_blank" rel="noopener noreferrer" href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2"> CH Play</a>.
                            </li>
                            <li>Trong Ứng dụng chọn <b>Thiết lập tài khoản</b>.</li>
                            <li>Chọn <b>Quét mã vạch</b>.</li>
                        </ul>
                        <div className="qrCode">
                            
                            {
                                isLoadOTP &&
                                <img className="loader-otp" alt="" src="/images/giphy.gif" />
                            }
                            {
                                !isLoadOTP &&
                                <img className="qr-code-png" alt="" src={srcImageQRCode} />
                            }
                        </div>
                    </div> 
                </div>
                {
                    isLoad &&
                    <div className="spinner">
                        <div className="_loader">
                            <div className="circle"></div>
                            <div className="circle"></div>
                            <div className="circle"></div>
                            <div className="circle"></div>
                            <div className="circle"></div>
                        </div>
                    </div>
                }
            </div>
        );
    }
 }
 export default Profile;
 