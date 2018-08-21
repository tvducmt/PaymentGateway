import React, {Component} from 'react';
import { Route, Switch, Redirect, NavLink } from 'react-router-dom'

import Header from './Header'

import Profile from './dashboard/Profile'
import HistoryTransaction from './dashboard/HistoryTransaction'
import DetailTransaction from './dashboard/DetailTransaction'

import ListCoupon from './dashboard/ListCoupon'

class Dashboard extends Component {
    
    render() {
       return (
            <div id="page-dashboard">
                <Header/>
                <ul className="navbar-nav navbar-sidenav">
                    <li className="nav-item">
                        <NavLink to="/dashboard/profile">
                            <i className="icons glyphicon glyphicon-user"></i>
                            <span>Hồ sơ cá nhân</span>
                        </NavLink>
                    </li>
                    <li className="nav-item">
                        <NavLink to="/dashboard/history-tx">
                            <i className="icons glyphicon glyphicon-calendar"></i>
                            <span>Lịch sử giao dịch</span>
                        </NavLink>
                    </li>
                    <li className="nav-item">
                        <NavLink to="/dashboard/coupon">
                            <i className="icons glyphicon glyphicon-credit-card"></i>
                            <span>Danh sách Coupon</span>
                        </NavLink>
                    </li>
                </ul>
                <Switch>
                    <Route exact path='/dashboard/profile' component={Profile} />
                    <Route exact path='/dashboard/history-tx' component={HistoryTransaction} />
                    <Route exact path='/dashboard/coupon' component={ListCoupon} />
                    <Route exact path='/dashboard/tx-detail/:id' component={DetailTransaction} />
                    <Route exact path='/dashboard' component={Profile} />
                    <Route render={() => (<Redirect to="/404"/>)}/>
                </Switch>
                {
                    // !this.state.isLoad &&
                    // <div className="spinner">
                    //     <div className="_loader">
                    //         <div className="circle"></div>
                    //         <div className="circle"></div>
                    //         <div className="circle"></div>
                    //         <div className="circle"></div>
                    //         <div className="circle"></div>
                    //     </div>
                    // </div>
                }
                

            </div>
       );
    }
 }
 export default Dashboard;
 