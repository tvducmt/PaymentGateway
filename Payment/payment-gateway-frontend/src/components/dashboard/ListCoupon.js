import React, {Component} from 'react';
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import { sendRequest } from '../../services/Http.service'
import { CopyToClipboard } from 'react-copy-to-clipboard';
import swal from 'sweetalert2'


class ListCoupon extends Component {
    constructor() {
        super()
        this.state = {
            copied: false,
            value: '',
            show: false,
            listCoupons: null,
            isAuth: AuthenticateService.isAuthenticate(),
            couponValue: 0,
            isLoad: true
        }
    }
    componentWillMount() {
        if(this.state.isAuth) {
            // coupon/coupon-info?userid=1
            // let userid = AuthenticateService.getAuthenticateUser().id
            let headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            sendRequest('get', '/user/coupons', null, headers).then((res) => {
                // console.log(res);
                if(!res.isError) {
                    if(res.data.error !== '') {
                        swal("Error!", `${res.data.error}`, 'error')
                    } else {
                        this.setState({ listCoupons: res.data.coupons })
                    }
                } else {
                    swal("Error!", `Có lỗi xảy ra. Vui lòng kiểm tra lại`, 'error')
                }
                this.setState({isLoad: false})
            }).catch((err) => {
                console.log(err);
            });
        }
    }
    formatDate(str) {
        var d = new Date(str);
        var day = (d.getUTCDate() <= 9)? "0" + d.getUTCDate() : d.getUTCDate();
        var month = (d.getUTCMonth() + 1 <= 9)? "0" + (d.getUTCMonth() + 1): (d.getUTCMonth() + 1);
        var year = (d.getUTCFullYear() <= 9)? "0" + d.getUTCFullYear() : d.getUTCFullYear();
        return day + "/" + month + "/" + year
    }
    formatMoney(money) {
        let m = money.toLocaleString('it-IT', {style : 'currency', currency : 'VND'});
        if (parseInt(m.slice(0, m.length - 1),10) !== 0) {
            return m.slice(0, m.length - 1);
        }
        return money + " "
    }
    renderListCoupon() {
        if (this.state.listCoupons != null) {
            // console.log(this.state.listCoupons)
            const listItems = this.state.listCoupons.map((el, idx) => {
                return (
                    <CopyToClipboard key={idx.toString()} onCopy={this.onCopy.bind(this, el.code)} text={el.code}>
                        <div className="coupon">
                            <div className="coupon-header">Coupon Rockship</div>
                            <div className="coupon-value">Giá trị: {this.formatMoney(el.value)}{el.currency.toUpperCase()}</div>
                            <div className="coupon-footer">
                                <div className="coupon-code">Mã Coupon <b>{el.code}</b></div>
                                <div className="coupon-exp">{this.formatDate(el.create_at)}</div>
                                <div className="coupon-state">{el.status === "unspend"? 'Chưa sử dụng' : 'Đã sử dụng'}</div>
                            </div>
                        </div>
                    </CopyToClipboard>
                )
            });
            return (<div className="wrapper-coupon">{listItems}</div>)
        }
    }
    onCopy = (value) => {
        var react = this
        react.setState({show: true});
        react.setState({value: value});
        setTimeout(function(){ react.setState({show: false}); }, 4000);
        react.setState({copied: true});
    }
    OnChangeCouponValue = (e) => {
        this.setState({couponValue: e.target.value});
    }
    // onSubmitCoupon = (e) => {
    //     e.preventDefault();
    //     if(this.state.isAuth) {
    //         if (!isNaN(parseInt(this.state.couponValue, 10))) {
    //             let data = {
    //                 user_id: AuthenticateService.getAuthenticateUser().id,
    //                 value: parseInt(this.state.couponValue, 10)
    //             }
    //             console.log(data);
    //             let headers = {
    //                 'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
    //             }
    //             sendRequest('post', '/coupon/new-transaction', data, headers).then((res) => {
    //                 console.log(res)
    //                 if(!res.isError) {
    //                     if(res.data.error === "") {
    //                         let data = this.state.listCoupons
    //                         data.push(res.data.coupons)
    //                         this.setState({ listCoupons: data })
    //                         swal(`Thành công!`, `Mã coupon của bạn là: ${res.data.coupons.code}`, 'success')
    //                     } else {
    //                         swal(`Lỗi!!!`, res.data.error, 'error')
    //                     }
    //                 } else {
    //                     swal(`Lỗi!!!`, `Có lỗi xảy ra trong quá trình xử lí. Xin mời thử lại`, 'error')
    //                 }
    //             });
    //         } else {
    //             swal(`Lỗi!!!`, `Giá trị coupon không hợp lệ, xin mời nhập lại`, 'error')
    //         }
            
    //     }
    // }
    render() {
       return (
            <div className="content-wrapper">
                <div className="wrapper-table">
                    {/* <div className="table-header"><i className="icons glyphicon glyphicon-info-sign"></i><span>Yêu cầu tạo Coupon</span></div>  
                    <div className="table-body">
                    <form onSubmit={this.onSubmitCoupon}>
                        <div className="row">
                            <div className="col-xs-4">
                                <div className="form-group">
                                    <label>Nhập giá trị coupon</label>
                                    <input type="number" value={this.state.couponValue} onChange={this.OnChangeCouponValue} className="form-control"placeholder="Nhập số tiền ..."/>
                                </div>
                                <div className="form-group">
                                    <button type="submit" className="btn btn-primary">Gửi yêu cầu</button>
                                </div>
                            </div>
                        </div>
                    </form>
                    </div> */}
                    <div className="table-header"><i className="icons glyphicon glyphicon-list-alt"></i><span>Danh sách coupon</span></div>        
                    <div className="table-body">
                        {
                            this.state.listCoupons != null &&
                            this.renderListCoupon()
                        }
                    </div>
                </div>
                <div className={this.state.show ? "tcl-alert open-alert" : "tcl-alert"}>
                    Mã Coupon {this.state.value} đã được copy vào clipboard
                    <span className="tcl-closebtn">&times;</span> 
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
 export default ListCoupon;
 