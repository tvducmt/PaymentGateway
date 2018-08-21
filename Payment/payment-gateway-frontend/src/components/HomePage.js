import React, {Component} from 'react';
import Slider from "react-slick";
import Header from './Header'
class HomePage extends Component {  
    componentWillMount() {
        console.log('componentWillMount');
    }
    
    componentDidMount() {
        console.log('componentDidMount');
    }
    
    componentWillReceiveProps() {
        console.log('componentWillReceiveProps');
    }
    
    shouldComponentUpdate() {
        console.log('shouldComponentUpdate');
    }
    
    componentWillUpdate() {
        console.log('componentWillUpdate');
    }
    
    componentDidUpdate() {
        console.log('componentDidUpdate');
    }
    
    componentWillUnmount() {
        console.log('componentWillUnmount');
    }
    render() {
        var settings = {
            dots: true,
            infinite: true,
            speed: 500,
            slidesToShow: 1,
            slidesToScroll: 1
        };
       return (
        <div>
            <Header />
            {/* <div className="container">
                
                <div className="row">
                    {
                        this.state.isAuth &&
                        <div className="col-xs-6 col-xs-push-3">
                            <div className="alert alert-success">
                                <button type="button" className="close" data-dismiss="alert" aria-hidden="true">&times;</button>
                                <strong>Thông báo!</strong> Chào mừng {AuthenticateService.getAuthenticateUser().email} đăng nhập thành công
                            </div>
                        </div>
                    }
                    
                </div>
            </div> */}
            <div className="wrapper-slider">
                <Slider {...settings}>
                    <div className="item">
                        <div className="tcl-table">
                            <div className="table-cell">
                                <img className="image-banner" alt="" src="/images/profile.png"/>
                                <h1 className="title">Welcome to Payment Gateway</h1>
                                <div className="wrapper-button">
                                    <a href="/login">Đăng nhập</a>
                                    <a href="/register">Đăng ký</a>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="item">
                        <div className="tcl-table">
                            <div className="table-cell">
                                <img className="image-banner" alt="" src="/images/profile.png"/>
                                <h1 className="title">Welcome to Payment Gateway</h1>
                                <div className="wrapper-button">
                                    <a href="/login">Đăng nhập</a>
                                    <a href="/register">Đăng ký</a>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="item">
                        <div className="tcl-table">
                            <div className="table-cell">
                                <img className="image-banner" alt="" src="/images/profile.png"/>
                                <h1 className="title">Welcome to Payment Gateway</h1>
                                <div className="wrapper-button">
                                    <a href="/login">Đăng nhập</a>
                                    <a href="/register">Đăng ký</a>
                                </div>
                            </div>
                        </div>
                    </div>
                </Slider>
            </div>
        </div>
       );
    }
 }
 export default HomePage;
 