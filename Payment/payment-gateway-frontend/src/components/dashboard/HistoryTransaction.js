import React, {Component} from 'react';
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { sendRequest } from '../../services/Http.service'
// import swal from 'sweetalert2'
import { Link } from "react-router-dom";

class HistoryTransaction extends Component {
    constructor() {
        super()
        console.log("Constructor Dashboard")
        this.state = {
            allData: null,
            isAuth: AuthenticateService.isAuthenticate(),
            intervalId: null,
            isLoad: true,
            filterValue: "*"
        }
    }
    componentDidMount() {
        // console.log('componentDidMount');
        var isAuth = this.state.isAuth;
        var elReact = this;
        var intervalId = setInterval(function(){
            if (isAuth) {
                let data = {
                    userid: AuthenticateService.getAuthenticateUser().id
                }
                let headers = {
                    'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
                }
                sendRequest('post', '/transactions', data, headers).then((res) => {
                    console.log(res)
                    if(!res.isError) {
                        if(!res.data.Error) {
                            elReact.setState({allData: res.data.data});
                        }
                    } else {
                        clearInterval(elReact.state.intervalId);
                    }
                    elReact.setState({isLoad: false});
                }).catch((err) => {
                    console.log(err);
                    elReact.setState({isLoad: false});
                });
            }
        }, 1000);
        elReact.setState({intervalId: intervalId});
            
    }
    componentWillUnmount() {
        // use intervalId from the state to clear the interval
        clearInterval(this.state.intervalId);
    }
    formatDate(str) {
        var d = new Date(str);
        var day = (d.getUTCDate() <= 9)? "0" + d.getUTCDate() : d.getUTCDate();
        var month = (d.getUTCMonth() + 1 <= 9)? "0" + (d.getUTCMonth() + 1): (d.getUTCMonth() + 1);
        var year = (d.getUTCFullYear() <= 9)? "0" + d.getUTCFullYear() : d.getUTCFullYear();
        var hour = (d.getUTCHours() <= 9)? "0" + d.getUTCHours() : d.getUTCHours();
        var minute = (d.getUTCMinutes() <= 9)? "0" + d.getUTCMinutes() : d.getUTCMinutes();
        return day + "/" + month + "/" + year + " " + hour + ":" + minute + '"'
    }
    formatMoney(money) {
        money = money.toLocaleString('it-IT', {style : 'currency', currency : 'VND'});
        return money.slice(0, money.length - 1);
    }
    // onCreateCoupon(event, tx_id) {
    //     if (this.state.isAuth) {
    //         let data = {
    //             user_id: AuthenticateService.getAuthenticateUser().id,
    //             tx_id: tx_id
    //         }
    //         let headers = {
    //             'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
    //         }
    //         event.target.disabled = true
    //         sendRequest('post', '/coupon/exist-transaction', data, headers).then((res) => {
    //             // console.log(res)
    //             if(!res.isError) {
    //                 if(res.data.Error !== "") {
    //                     swal('Error!', `${res.data.Error}`, 'error')
    //                 } else {
    //                     swal('Success!', `Đổi coupon thành công. Code: ${res.data.Coupon.code}`, 'success')
    //                 }
    //             } else {
    //                 swal('Error!', `Có lỗi xảy ra, vui lòng thử lại`, 'error')
    //             }

    //         }).catch((err) => {
    //             console.log(err);
    //         });
    //     }
    // }

    renderBankTransfer(el, idx) {
        // console.log("renderBankTransfer")
        return (
            <React.Fragment>
                <td className="id_receipt">{idx + 1}</td>
                <td className="">Bank Transfer</td>
                {
                    el.status.toLowerCase() === "pending" &&
                    <td className="raw_receipt">{el.raw_receipt.String}</td>
                }
                {
                    el.status.toLowerCase() === "confirmed" && el.raw_receipt.String !== "" &&
                    <td className="raw_receipt">
                        <div>Tài khoản ngân hàng: <b>{el.parsed_account.String}</b></div>
                        <div>Số tiền: <b>{this.formatMoney(el.parsed_amount.Float64)} VND</b></div>
                        {/* <div>Biên nhận: <b>{el.raw_receipt.String.substring(0,40) + '...'}</b></div> */}
                        <div><Link to={`/dashboard/tx-detail/${el.id}`}>Chi tiết</Link></div>
                    </td>
                }
                <th className="tx_code">{el.code.Valid ? "CODE"+el.code.String : ""}</th>
                <th className="tx_code coupon_code">{el.coupon_code.String}</th>
                <td className="status"><span className={el.status.toLowerCase()}>{el.status}</span></td>
                <td className="create_at">{this.formatDate(el.create_at)}</td>
            </React.Fragment>
        ); 
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
    renderEtherTransfer(el, idx) {
        // console.log("renderBankTransfer")
        return (
            <React.Fragment>
                <td className="id_receipt">{idx + 1}</td>
                <td className="">Ether Transfer</td>
                <td className="raw_receipt"><b>
                    Hex code: {this.hexEncode("CODE"+el.code.String)}</b>
                    {
                        el.status.toLowerCase() === "confirmed" &&
                        <div><Link to={`/dashboard/tx-detail/${el.id}`}>Chi tiết</Link></div>
                    }
                </td>
                <th className="tx_code">{el.code.Valid ? "CODE"+el.code.String : ""}</th>
                <th className="tx_code coupon_code">{el.coupon_code.String}</th>
                <td className="status"><span className={el.status.toLowerCase()}>{el.status}</span></td>
                <td className="create_at">{this.formatDate(el.create_at)}</td>
            </React.Fragment>
        ); 
    }
    renderCredit(el, idx) {
        // console.log("renderCredit")
        return (
            <React.Fragment>
                <td className="id_receipt">{idx + 1}</td>
                <td className="">Credit/Debit Card</td>
                <td className="raw_receipt">
                    <div>Giá trị: <b>{this.formatMoney(el.coupon_value.Float64)} {el.coupon_currency.String.toUpperCase()}</b></div>
                    <div><Link to={`/dashboard/tx-detail/${el.id}`}>Chi tiết</Link></div>
                    {/* <div>Charge ID: {el.charge_id.String}</div> */}
                </td>
                <td className="tx_code"></td>
                <td className="tx_code coupon_code"><b>{el.coupon_code.String}</b></td>
                <td className="status"><span className={el.status.toLowerCase()}>{el.status}</span></td>
                <td className="create_at">{this.formatDate(el.create_at)}</td>
            </React.Fragment>
        )
    }
    handleFilter() {
        switch(this.state.filterValue) {
            case "bank-transfer":
                return this.state.allData.filter(el => el.method_name.String === "Bank Transfer")
            case "credit-card":
                return this.state.allData.filter(el => el.method_name.String === "Credit")
            case "ether-transfer":
                return this.state.allData.filter(el => el.method_name.String === "Ether Transfer")
            case "successed":
                return this.state.allData.filter(el => el.status === "successed")
            case "pending":
                return this.state.allData.filter(el => el.status === "pending")
            case "confirmed":
                return this.state.allData.filter(el => el.status === "confirmed")
            default:
                return this.state.allData

        }
        // console.log(datafilter)
    }
    renderListTransaction() {
        if(this.state.allData != null) {
            let datafilter = this.handleFilter()
            
            const listItems = datafilter.map((el, idx) => {
                return (
                    <tr key={idx.toString()} data_tx_id={el.id}>
                    {
                        el.method_name.String === "Bank Transfer" && this.renderBankTransfer(el, idx)
                    }
                    {
                        el.method_name.String === "Ether Transfer" && this.renderEtherTransfer(el, idx)
                    }
                    {
                        el.method_name.String === "Credit" && this.renderCredit(el, idx)
                    }
                    </tr>
                )
            })

            return (<tbody className="render-data">{listItems}</tbody>);
        }
    }
    FilterClick(e) {
        e.preventDefault();
        this.setState({filterValue: e.target.getAttribute('data-filter')})
        // console.log(e.target.getAttribute('data-filter'))
    }
    DefaultClick(e) {
        e.preventDefault();
    }
    render() {
       return (
            <div className="content-wrapper">
                <div className="wrapper-table">
                    <div className="table-header">
                        <i className="icons glyphicon glyphicon-calendar"></i>
                        <span>Lịch sử giao dịch</span>
                    </div>
                    <div className="table-body">
                        <div className="filter-data">
                            <span>Filter: </span>
                            <a href="" onClick={(e) => this.FilterClick(e)} data-filter="*">Tất cả</a>
                            <a href="" onClick={(e) => this.DefaultClick(e)} class="tcl-dropdown-wrapper">
                                <span>Hình thức <i className="icons glyphicon glyphicon-menu-down"></i></span>
                                <span className="tcl-dropdown">
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="bank-transfer">Bank Transfer</a>
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="ether-transfer">Ether Transfer</a>
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="credit-card">Credit Card</a>
                                </span>
                            </a>
                            <a href="" onClick={(e) => this.DefaultClick(e)} class="tcl-dropdown-wrapper">
                                <span>Trạng thái <i className="icons glyphicon glyphicon-menu-down"></i></span>
                                <span className="tcl-dropdown">
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="pending">Đang chờ</a>
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="successed">Thành công</a>
                                    <a href="" onClick={(e) => this.FilterClick(e)} data-filter="confirmed">Xác nhận</a>
                                </span>
                            </a>
                        </div>
                        <div className="clearfix"></div>
                        <div className="table-responsive">
                            <table className="table table-bordered table-hover table_receipt_log">
                                <thead>
                                    <tr>
                                        <th className="id_receipt">#</th>
                                        <th className="" style={{width: 140}}>Hình thức</th>
                                        <th className="raw_receipt">Biên nhận chi tiết</th>
                                        <th className="tx_code">Mã giao dịch</th>
                                        <th className="tx_code">Mã Coupon</th>
                                        <th className="status">Trạng thái</th>
                                        <th className="create_at">Ngày tạo</th>
                                        {/* <th className="coupon">Coupon</th> */}
                                    </tr>
                                </thead>
                                    {this.renderListTransaction()}
                            </table>
                        </div>
                    </div>
                </div>
                {
                    this.state.isLoad &&
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
 export default HistoryTransaction;
 