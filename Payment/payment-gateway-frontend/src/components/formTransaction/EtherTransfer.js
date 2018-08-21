import React, {Component} from 'react';
import { sendRequest } from '../../services/Http.service'
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { CopyToClipboard } from 'react-copy-to-clipboard';
import swal from 'sweetalert2'

class EtherTransfer extends Component {
    constructor() {
        super();
        this.state = {
            isCreated: false,
            isAuth: AuthenticateService.isAuthenticate(),
            isLoading: false,
            txCode: '',
            hasError: false,
            valueCopy: '',
            copied: false,
            show: false
        }
    }
    hexEncode = function(str){
        var hex, i;
    
        var result = "";
        for (i=0; i<str.length; i++) {
            hex = str.charCodeAt(i).toString(16);
            result += hex
        }
        return "0x" + result
    }
    onSubmit = (e) => {
        e.preventDefault();
        // action="/transaction" method="POST"
        this.setState({isLoading: true});
        let data = {
            user_id: AuthenticateService.getAuthenticateUser().id,
            method_id: 4,
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
                    swal(`CODE${res.data.TxCode}`, `Thành công, mã hex đã được tạo bên dưới`, 'success')
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
    onCopy = (value) => {
        var react = this
        react.setState({show: true});
        react.setState({valueCopy: value});
        setTimeout(function(){ react.setState({show: false}); }, 4000);
        react.setState({copied: true});
    };
    render() {
       return (
        <form onSubmit={this.onSubmit} id="form-ether-transfer">
            <div className="resp-tabs-container">
                {/* <input type="userID" name="methodID" value="1" className="form-control hidden" placeholder="Input field" /> */}
                <h2>Ether Transfer</h2>
                <ul>
                    <li>Hướng dẫn: </li>
                    <li>Step 1: ...</li>
                    <li>Step 2: ...</li>
                    <li>Step 3: ...</li>
                    <li><span className="note">Lưu ý: Dữ liệu chuyển ether bao gồm mã code bên dưới ở dạng hex</span></li>
                </ul>
                
                <div className="wrapper-tx">
                    <button type="submit" 
                        disabled={this.state.isLoading || this.state.isCreated} 
                        className={this.state.isLoading? 'btn btn-success btn-create-tx sending':'btn btn-success btn-create-tx'}>
                        <div className="loader"></div>
                        Tạo giao dịch
                    </button>
                    <div className="tx-code">{this.state.txCode === '' ? 'YOUR CODE' : 'CODE' + this.state.txCode}</div>
                    
                    {
                        this.state.txCode === '' &&
                        <React.Fragment>
                        <h2>Hex Code</h2>  
                        <div className="tx-code hex-code">{this.state.txCode === '' ? '0x.....' : this.hexEncode('CODE' + this.state.txCode)}</div>
                        </React.Fragment>
                    }
                    {
                        this.state.txCode !== '' &&
                        <React.Fragment>
                            <h2>Hex Code</h2>     
                            <CopyToClipboard onCopy={this.onCopy.bind(this, this.hexEncode('CODE' + this.state.txCode))} text={this.hexEncode('CODE' + this.state.txCode)}>
                            <div className="tx-code hex-code">{this.state.txCode === '' ? '0x.....' : this.hexEncode('CODE' + this.state.txCode)}</div>
                            </CopyToClipboard>
                        </React.Fragment>
                    }
                </div>
                
                {
                    this.state.hasError &&
                    <div className="alert alert-danger">
                        <button type="button" className="close" data-dismiss="alert" aria-hidden="true">&times;</button>
                        <strong>Error!</strong> Có lỗi xảy ra, vui lòng thử lại ...
                    </div>
                }
            </div>
            <div className={this.state.show ? "tcl-alert open-alert" : "tcl-alert"}>
                Mã Hex {this.state.valueCopy} đã được copy vào clipboard
                <span className="tcl-closebtn">&times;</span> 
            </div>
        </form>
       );
    }
 }
 export default EtherTransfer;
 