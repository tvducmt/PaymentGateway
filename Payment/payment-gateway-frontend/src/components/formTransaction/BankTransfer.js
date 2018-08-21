import React, {Component} from 'react';
import { sendRequest } from '../../services/Http.service'
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { CopyToClipboard } from 'react-copy-to-clipboard';
import swal from 'sweetalert2'

class BankTransfer extends Component {
    constructor() {
        super();
        this.state = {
            isCreated: false,
            isAuth: AuthenticateService.isAuthenticate(),
            isLoading: false,
            txCode: '',
            hasError: false,
            valueCopy: '',
            copied: false
        }
    }
    onSubmit = (e) => {
        e.preventDefault();
        // action="/transaction" method="POST"
        this.setState({isLoading: true});
        let data = {
            user_id: AuthenticateService.getAuthenticateUser().id,
            method_id: 1,
            gateway_id: 1
        }
        let headers = {
            'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
        }
        sendRequest('post', '/bank-transfer', data, headers).then((res) => {
            // console.log(res)
            if(!res.isError) {
                if(res.data.Error === "") {
                    this.setState({isCreated: true});
                    this.setState({txCode: res.data.TxCode })
                    swal(`CODE${res.data.TxCode}`, `Mã CODE đã được copy vào clipboard`, 'success')
                } else if (res.data.Error === "UserNotExist"){
                    AuthenticateService.removeAuthenticate()
                } else {
                    this.setState({hasError: true});
                }
            } else {
                this.setState({hasError: true })
            }
            this.setState({isLoading: false}); 
        });
    }
    onCopy = () => {
        this.setState({copied: true});
    };

    render() {
        return (
            <form onSubmit={this.onSubmit} id="form-bank-transfer">
                <div className="resp-tabs-container">
                    {/* <input type="userID" name="methodID" value="1" className="form-control hidden" placeholder="Input field" /> */}
                    <h2>Bank Transfer</h2>
                    <ul>
                        <li>Hướng dẫn: </li>
                        <li>Step 1: ...</li>
                        <li>Step 2: ...</li>
                        <li>Step 3: ...</li>
                        <li><span className="note">Lưu ý: Nội dung chuyển khoản phải bao gồm mã Code bên dưới</span></li>
                    </ul>
                    <CopyToClipboard onCopy={this.onCopy} text={"CODE"+this.state.txCode}>
                    <div className="wrapper-tx">
                        <button type="submit" 
                            disabled={this.state.isLoading || this.state.isCreated} 
                            className={this.state.isLoading? 'btn btn-success btn-create-tx sending':'btn btn-success btn-create-tx'}>
                            <div className="loader"></div>
                            Tạo giao dịch
                        </button>
                        <div className="tx-code">{this.state.txCode === '' ? 'YOUR CODE' : 'CODE' + this.state.txCode}</div>
                    </div>
                    </CopyToClipboard>
                    {
                        this.state.hasError &&
                        <div className="alert alert-danger">
                            <button type="button" className="close" data-dismiss="alert" aria-hidden="true">&times;</button>
                            <strong>Error!</strong> Có lỗi xảy ra, vui lòng thử lại ...
                        </div>
                    }

                    <div className="wrapper-info">
                        <div className="image-bank">
                            <img src="/images/vietinbank.png" alt=""/>
                        </div>
                        <div className="info-bank">
                            <ul>
                                <li>Ngân hàng TMCP Công Tương Việt Nam - Vietinbank</li>
                                <li>Chủ tài khoản: Rockship</li>
                                <li>Số tài khoản: 123456789</li>
                                <li>Chi nhánh: TPHCM</li>
                            </ul>
                            {/* <input type="userID" className="hidden form-control" name="gatewayID" value="1" placeholder="Input field" /> */}
                        </div>
                    </div>
                    <div className="wrapper-info">
                        <div className="image-bank">
                            <img src="/images/vietcombank.png" alt=""/>
                        </div>
                        <div className="info-bank">
                            <ul>
                                <li>Ngân hàng TMCP Công Tương Việt Nam - Vietinbank</li>
                                <li>Chủ tài khoản: Rockship</li>
                                <li>Số tài khoản: 123456789</li>
                                <li>Chi nhánh: TPHCM</li>
                            </ul>
                        </div>
                    </div>
                </div>
            </form>
       );
    }
 }
 export default BankTransfer;
 