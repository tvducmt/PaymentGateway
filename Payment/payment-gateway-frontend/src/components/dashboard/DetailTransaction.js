import React, {Component} from 'react';
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { sendRequest } from '../../services/Http.service'
import swal from 'sweetalert2'

class DetailTransaction extends Component {
    state = {
        isAuth: AuthenticateService.isAuthenticate(),
        transaction: null,
        charge: null
    }
    componentWillMount() {
        if(this.state.isAuth) {
            let data = {
                tx_id: parseInt(this.props.match.params.id, 10)
            }
            let headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            sendRequest('get', '/transaction?tx_id=' + data.tx_id, null, headers)
            .then((res) => {
                console.log(res)
                if(!res.isError) {
                    if(typeof(res.data.error) === 'undefined') {
                        this.setState({ transaction: res.data.transaction })
                        if (res.data.charge != null) {
                            let data = {
                                amount: res.data.charge.amount,
                                currency: res.data.charge.currency,
                                address_zip: res.data.charge.source.address_zip,
                                brand: res.data.charge.source.brand,
                                exp_month: res.data.charge.source.exp_month,
                                exp_year: res.data.charge.source.exp_year,
                                funding: res.data.charge.source.funding,
                                last4: res.data.charge.source.last4,
                                status: res.data.charge.status,
                                refunded: res.data.charge.refunded,
                                amount_refunded: res.data.charge.amount_refunded
                            }
                            this.setState({charge: data})
                        }
                    } else {
                        swal('Error!', `${res.data.error}`, 'error')
                    }
                } else {
                    swal('Error!', `Có lỗi xảy ra, vui lòng thử lại`, 'error')
                }
            }).catch((err) => {
                console.log(err)
            });
        }
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
    formatMoney(money, method = "default") {
        if (method === "default") {
            money = money.toLocaleString('it-IT', {style : 'currency', currency : 'VND'});
            return money.slice(0, money.length - 1);
        } else {
            return money
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
    renderBanktransfer() {
        let data = this.state.transaction
        return (
            <React.Fragment>
                <tr>
                    <td  style={{width: 180}}>Loại giao dịch</td>
                    <td>{data.method_name.String}</td>
                </tr>
                <tr>
                    <td>Tài khoản</td>
                    <td>{data.parsed_account.String}</td>
                </tr>
                <tr>
                    <td>Số tiền</td>
                    {
                        data.method_name.String === "Ether Transfer" &&
                        <td>{this.formatMoney(data.coupon_value.Float64, "") + " " + data.coupon_currency.String.toUpperCase()}</td>
                    }
                    {
                        data.method_name.String !== "Ether Transfer" &&
                        <td>{this.formatMoney(data.parsed_amount.Float64) + " " + data.coupon_currency.String.toUpperCase()}</td>
                    } 
                </tr>
                <tr>
                    <td>Số điện thoại</td>
                    <td>{data.phonenumber.String}</td>
                </tr>
                <tr>
                    <td>Ngày tạo giao dịch</td>
                    <td>{this.formatDate(data.create_at)}</td>
                </tr>
                <tr>
                    <td>Mã giao dịch</td>
                    <td>CODE{data.code.String}</td>
                </tr>
                <tr>
                    <td>Ngày nhận biện nhận</td>
                    { data.method_name.String === "Ether Transfer" && <td>{this.formatDate(data.create_at)}</td> }
                    { data.method_name.String !== "Ether Transfer" && <td>{this.formatDate(data.sys_create_at.Time)}</td> }
                </tr>
                <tr>
                    <td>Biên nhận</td>
                    <td>{data.raw_receipt.String || this.hexEncode("CODE"+data.code.String)}</td>
                </tr>
                <tr>
                    <td>Trạng thái</td>
                    <td>{data.status}</td>
                </tr>
                <tr>
                    <td>Giá trị Coupon</td>
                    {
                        data.method_name.String === "Ether Transfer" &&
                        <td>{this.formatMoney(data.coupon_value.Float64, "")} {data.coupon_currency.String.toUpperCase()}</td>
                    }
                    {
                        data.method_name.String !== "Ether Transfer" &&
                        <td>{this.formatMoney(data.coupon_value.Float64)} {data.coupon_currency.String.toUpperCase()}</td>
                    } 
                </tr>
                <tr>
                    <td>Mã Coupon</td>
                    <td>{data.coupon_code.String}</td>
                </tr>
                <tr>
                    <td>Trạng thái Coupon</td>
                    <td>
                        {data.coupon_status.String === "unspend"?"Chưa sử dụng" : ""}
                        {data.coupon_status.String === "spend"?"Đã sử dụng" : ""}
                    </td>
                </tr>
            </React.Fragment>
        );
    }
    renderCreditCard() {
        let data = this.state.transaction
        let charge = this.state.charge
        return (
            <React.Fragment>
                <tr>
                    <td  style={{width: 180}}>Loại giao dịch</td>
                    <td>Credit/Debit Card</td>
                </tr>
                <tr>
                    <td>Ngày tạo giao dịch</td>
                    <td>{this.formatDate(data.create_at)}</td>
                </tr>
                <tr>
                    <td>Giá trị Coupon</td>
                    <td>{this.formatMoney(data.coupon_value.Float64)} {data.coupon_currency.String.toUpperCase()}</td>
                </tr>
                <tr>
                    <td>Mã Coupon</td>
                    <td>{data.coupon_code.String}</td>
                </tr>
                <tr>
                    <td>Trạng thái Coupon</td>
                    <td>
                        {data.coupon_status.String === "unspend"?"Chưa sử dụng" : ""}
                        {data.coupon_status.String === "spend"?"Đã sử dụng" : ""}
                    </td>
                </tr>
                <tr>
                    <td>Loại thẻ</td>
                    <td>{charge.brand} - {charge.funding}</td>
                </tr>
                <tr>
                    <td>Ngày hết hạn</td>
                    <td>{charge.exp_month}/{charge.exp_year}</td>
                </tr>
                <tr>
                    <td>Số cuối thẻ</td>
                    <td>{charge.last4}</td>
                </tr>
                <tr>
                    <td>Mã bưu điện</td>
                    <td>{charge.address_zip}</td>
                </tr>
                <tr>
                    <td>Hoàn tiền</td>
                    <td>{charge.refunded ? "Đã hoàn tiền" : "Chưa hoàn tiền"}</td>
                </tr>
                <tr>
                    <td>Số tiền đã hoàn lại</td>
                    <td>{charge.amount_refunded}</td>
                </tr>
                <tr>
                    <td>Trạng thái</td>
                    <td>{charge.status === "succeeded"? "Thành công" : "..." }</td>
                </tr>
            </React.Fragment>
        );
    }
    renderInfo() {
        console.log(this.state.transaction)
        return (
            <React.Fragment>
                {
                    this.state.transaction !== null && this.state.charge === null &&
                    this.renderBanktransfer()
                }
                {
                    this.state.transaction !== null && this.state.charge !== null &&
                    this.renderCreditCard()
                }
            </React.Fragment>
        );
    }
    render() {
       return (
        <div className="content-wrapper">
            <div className="wrapper-table">
                <div className="table-header">
                    <i className="icons glyphicon glyphicon-calendar"></i>
                    <span>Chi tiết giao dịch #{this.props.match.params.id}</span>
                </div>
                <div className="table-body">
                    <div className="clearfix"></div>
                    <div className="table-responsive">
                        <table className="table table-bordered table-transaction-detail" style={{width: 700}}>
                            <tbody>
                                {this.renderInfo()}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
       );
    }
 }
 export default DetailTransaction;
 