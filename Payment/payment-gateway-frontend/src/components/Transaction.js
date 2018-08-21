import React, {Component} from 'react';
import { Switch, Route, NavLink, Redirect } from 'react-router-dom';
import CreditCard from './formTransaction/CreditCard';
import BankTransfer from './formTransaction/BankTransfer';
import EtherTransfer from './formTransaction/EtherTransfer';
import Header from './Header'
class Transaction extends Component {
    render() {
       return (
        <div id="page-transaction">
            <Header/>
            <div className="container">
                <div className="content">
                    <div className="sap_tabs">
                        <div className="pay-tabs">
                            <h2>Phương thức thanh toán</h2>
                            <ul className="resp-tabs-list">
                                <li className="resp-tab-item" aria-controls="tab_item-0" role="tab">
                                    <NavLink to='/transaction/credit-card' >
                                        <span><label className="pic1"></label>Credit Card</span>
                                    </NavLink>
                                </li>
                                <li className="resp-tab-item resp-tab-active" aria-controls="tab_item-1" role="tab">
                                    <NavLink to='/transaction/bank-transfer' >    
                                        <span><label className="pic3"></label>Bank Transfer</span>
                                    </NavLink>
                                </li>
                                <li className="resp-tab-item" aria-controls="tab_item-2" role="tab">
                                    <NavLink to='/transaction/ether-transfer' >    
                                        <span><label className="pic4"></label>Ether Transfer</span>
                                    </NavLink>
                                    {/* <span><label className="pic4"></label>PayPal</span> */}
                                </li>
                                <div className="clear"></div>
                            </ul>
                        </div>
                        <Switch>
                            <Route exact path='/transaction/credit-card' component={CreditCard} />
                            <Route exact path='/transaction/bank-transfer' component={BankTransfer} />
                            <Route exact path='/transaction/ether-transfer' component={EtherTransfer} />
                            <Route render={() => (<Redirect to="/transaction/bank-transfer"/>)}/>
                        </Switch>
                    </div>       
                </div>
            </div>
        </div>
       );
    }
}
export default Transaction;
 