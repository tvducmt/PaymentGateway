import React, {Component} from 'react';
import { AuthenticateService } from '../services/AuthenticateService'
import { NavLink } from 'react-router-dom'

class Header extends Component {
    constructor() {
        super();
        this.state = {
            gmail: AuthenticateService.getAuthenticateGmail(),
			isAuth: AuthenticateService.isAuthenticate(),
		}
    }
    handleClick(event) {
        event.preventDefault();
        AuthenticateService.removeAuthenticate();
    }
    render() {
        const { isAuth } = this.state
       return (
            <nav className="navbar navbar-default">
                <div className="container-fluid">
                    <div className="navbar-header">
                        <button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
                            <span className="sr-only">Toggle navigation</span>
                            <span className="icon-bar"></span>
                            <span className="icon-bar"></span>
                            <span className="icon-bar"></span>
                        </button>
                        <a className="navbar-brand" href="/">Payment Gateway</a>
                    </div>

                    <div className="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
                        <ul className="nav navbar-nav">
                            <li>
                                <NavLink to='/dashboard'>Dashboard<span className="sr-only">(current)</span></NavLink>
                            </li>
                            <li>
                                <NavLink to='/transaction'>Giao dịch</NavLink>
                            </li>
                        </ul>
                        {
                            !isAuth &&
                            <ul className="nav navbar-nav navbar-right">
                                <li><a href="/login">Đăng nhập</a></li>
                            </ul>
                        }
                        {
                            isAuth &&
                            <ul className="nav navbar-nav navbar-right">
                                <li><a href="/" onClick={(e) => {this.handleClick(e)}}>Đăng xuất</a></li>
                                <li><a href="/">{this.state.gmail}</a></li>
                            </ul>
                        }
                            
                    </div>
                </div>
            </nav>
       );
    }
 }
 export default Header;
 