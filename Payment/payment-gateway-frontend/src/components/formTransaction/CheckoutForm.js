import React, {Component} from 'react';
import { injectStripe, CardElement } from 'react-stripe-elements';
import { sendRequest } from '../../services/Http.service'
import { AuthenticateService } from '../../services/AuthenticateService'
import { LocalStorageService } from '../../services/LocalStorageService'
import swal from 'sweetalert2'

class CheckoutForm extends Component {
    state = {
        isLoading: false,
        name: AuthenticateService.getAuthenticateUser().username,
        amount: '',
        amountErr: false,
        selectedOption: 'new',
        selectedCreditOld: '',
        listCredits: null,
        isAuth: AuthenticateService.isAuthenticate(),
        remember_me: false,
        isLoadingHistoryCredit: false
    }
	onChangeName = (e) => {
		this.setState({ name: e.target.value })
    }
    onChangeAmount = (e) => {
        this.setState({ amount: e.target.value })
    }
    handleOptionChange = (e) => {
        if (e.target.value === 'old') {
            this.sendRequestGetAllCustomer()
        }
        this.setState({selectedOption: e.target.value});
    }
    onChangeCreditOld = (e) => {
        this.setState({selectedCreditOld: e.target.value});
    }
    onChangeRememberMe = (e) => {
        this.setState({remember_me: !this.state.remember_me});
    }
    sendRequestGetAllCustomer() {
        let headers = {
            'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
        }
        this.setState({isLoadingHistoryCredit: true })
        sendRequest('get', '/customer/cards', null, headers).then((res) => {
            // console.log(res.data)
            if (!res.isError) {
                if(typeof(res.data.error) !== 'undefined') {
                    swal('Error!', `${res.data.error}`, 'error')
                    this.setState({ listCredits: null })
                } else {
                    this.setState({ listCredits: res.data })
                }
            }
        })
        .then(() => {
            this.setState({isLoadingHistoryCredit: false })
        })
        .catch(() => {
            console.log("Error");
        });
    }
    componentWillMount() {
        // this.sendRequestGetAllCustomer()
    }
    renderListCreditHistory() {
        // console.log(this.state.listCredits)
        if (this.state.listCredits !== null){
            const listItems = this.state.listCredits.map((el) =>
                <div className="wrapper-credit-card" key={el.id}>
                    {/* Loại thẻ: {el.brand + " " + el.funding}  */}
                    <div className="credit-card radio-style-2 choose-credit">
                        <div className="radio">
                            <input type="radio" name="credit-card" value={el.id} 
                                checked={this.state.selectedCreditOld === el.id.toString()}
                                onChange={this.onChangeCreditOld}/>
                            <label className="label-check">
                                <svg width="18px" height="18px" viewBox="0 0 18 18">
                                    <path d="M1,9 L1,3.5 C1,2 2,1 3.5,1 L14.5,1 C16,1 17,2 17,3.5 L17,14.5 C17,16 16,17 14.5,17 L3.5,17 C2,17 1,16 1,14.5 L1,9 Z"></path>
                                    <polyline points="1 9 7 14 15 4"></polyline>
                                </svg>
                            </label>
                        </div>   
                        <div className="card-number">XXXX XXXX XXXX {el.last4}</div>
                        <div className="right-content">
                            <div className="card-exp">{el.exp_month < 10? '0' + el.exp_month : el.exp_month} / {el.exp_year.toString().slice(2,4)}</div>
                            <div className="card-cvc">CVC</div>
                            <div className="card-zip">{el.address_zip}</div>
                        </div>
                        <img className="card-type" alt="" src={"/images/" + el.brand + ".png"} />
                    </div>
                </div>
            );
            if (listItems.length !== 0) {
                return (
                    <React.Fragment>
                        <h4>Chọn một thẻ bên dưới để tiến hành giao dịch</h4>{listItems}
                    </React.Fragment>
                );
            }
        }
        return ( <div className="form-group list-card"><h4>Không có thẻ nào được lưu</h4></div>)
    }
	onSubmit = (e) => {
        e.preventDefault();
        const {name, amount, selectedOption, selectedCreditOld, remember_me} = this.state
        if (isNaN(parseInt(amount, 10)) || parseInt(amount, 10) <= 0) {
            this.setState({ amountErr: true })
            return
        }    
        this.setState({ amountErr: false })
        if (this.props.stripe) {
            this.setState({isLoading: true})
            var headers = {
                'Authorization': 'Bearer ' + LocalStorageService.get('accesstoken')
            }
            if (selectedOption === "new") {
                this.props.stripe.createToken().then(({token}) => {
                    // console.log('Received Stripe token:', token);
                    let data = {
                        name: name,
                        amount: parseInt(amount, 10),
                        token: token.id,
                        select_mode_card: selectedOption,
                        select_cus_id: 0,
                        remember_me: remember_me
                    }
                    // console.log(data)
                    sendRequest('post', '/credit-card', data, headers).then((res) => {
                        // console.log(res)
                        if(!res.isError) {
                            if(res.data.error !== '') {
                                swal('Error!', `${res.data.error}`, 'error')
                                this.setState({isLoading: false})
                            } else {
                                swal('Congrats! Giao dịch thành công!', '', 'success')
                                this.props.history.push("/dashboard/coupon")
                            }
                        } else {
                            swal('Error!', 'Có lỗi xảy ra, vui lòng thử lại!', 'error')
                            this.setState({isLoading: false})
                        }
                    })
                }).catch((error) => {
                    console.log("Error: " + error)
                });
            } else if (selectedOption === "old") {
                if (isNaN(parseInt(selectedCreditOld,10))) {
                    if(this.state.listCredits.length !== 0) 
                        swal('Error!', 'Vui lòng chọn một thẻ bất kỳ để tiếp tục giao dịch!', 'error')
                    this.setState({isLoading: false})
                } else {
                    let data = {
                        name: name,
                        amount: parseInt(amount, 10),
                        token: "",
                        select_mode_card: selectedOption,
                        select_cus_id: parseInt(selectedCreditOld, 10),
                        remember_me: false
                    }
                    // console.log(data)
                    sendRequest('post', '/credit-card', data, headers).then((res) => {
                        // console.log(res)
                        if(!res.isError) {
                            if(res.data.error !== '') {
                                swal('Error!', `${res.data.error}`, 'error')
                                this.setState({isLoading: false})
                            } else {
                                swal('Congrats! Giao dịch thành công!', '', 'success')
                                this.props.history.push("/dashboard/coupon")
                            }
                        } else {
                            swal('Error!', 'Có lỗi xảy ra, vui lòng thử lại!', 'error')
                            this.setState({isLoading: false})
                        }
                        // this.setState({isLoading: false})
                    })
                }
                    
            }
            
        } else {
            console.log('Stripe.js tải lên thất bại')
        }
	}
    render() {
        const { isLoading, name, amount, amountErr, selectedOption, remember_me, isLoadingHistoryCredit } = this.state
       return (
        <form onSubmit={this.onSubmit}>
            <h2>Credit hoặc Debit card</h2>
            <div className="form-group type-credit-card">
                <div className="radio">
                    <input type="radio" value="new" name="type_credit" 
                        checked={selectedOption === 'new'}
                        onChange={this.handleOptionChange} />
                    <label>Giao dịch bằng thẻ mới</label>
                </div>
                <div className="radio">
                    <input type="radio" value="old" name="type_credit" 
                        checked={selectedOption === 'old'}
                        onChange={this.handleOptionChange} />
                    <label>Sử dụng thẻ cũ</label>
                </div>
            </div>
            <div className={isLoadingHistoryCredit ? 'loader-history-credit' : ''}>
                
                <div className="form-group name-credit-card">
                    <input type="text" onChange={this.onChangeName} value={name} className="form-control" placeholder="Email ..."/>
                </div>
                <div className="form-group name-credit-card">
                    <input type="number" onChange={this.onChangeAmount} value={amount} min="1" className={amountErr? 'form-control amount-invalid':'form-control'} placeholder="Số tiền ..."/>
                    <div className="input-group-addon">VND</div>
                </div>
                {
                    this.state.selectedOption === 'new' &&
                    <div>
                        <CardElement
                            className='form-group card-wrapper'
                            style={{ base: { fontSize: '18px', color: "#555" } }}
                        />
                        <div className="remember-me radio-style-2">
                            <div className="radio">
                                <input type="checkbox" name="check-remember" value={remember_me}
                                checked={remember_me === true}
                                onChange={this.onChangeRememberMe} />
                                <label className="label-check">
                                    <svg width="18px" height="18px" viewBox="0 0 18 18">
                                        <path d="M1,9 L1,3.5 C1,2 2,1 3.5,1 L14.5,1 C16,1 17,2 17,3.5 L17,14.5 C17,16 16,17 14.5,17 L3.5,17 C2,17 1,16 1,14.5 L1,9 Z"></path>
                                        <polyline points="1 9 7 14 15 4"></polyline>
                                    </svg>
                                </label>
                            </div>
                            <label>Lưu thông tin thẻ</label>
                        </div>
                    </div>
                }
                {
                    this.state.selectedOption === 'old' && this.state.listCredits != null &&
                    <div className="form-group list-card">
                        {this.renderListCreditHistory()}
                    </div>
                }
                {
                    isLoadingHistoryCredit &&
                    <img className="loader-gif" alt="" src="/images/giphy.gif" />
                }
            </div>   
            <div className="form-group">
                <button type="submit" className="btn btn-success btn-create-tx btn-credit-card" disabled={isLoading}>
                    {isLoading ? 'Loading ...' : 'Submit Payment'}
                </button>
                <div className="noti-credit-card">
                </div>
            </div>
            
        </form>
       );
    }
 }
 export default injectStripe(CheckoutForm);
 