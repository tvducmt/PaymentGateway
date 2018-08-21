import React, {Component} from 'react';
import { Elements, StripeProvider } from 'react-stripe-elements';
import InjectedCheckoutForm from './CheckoutForm'
class CreditCard extends Component {

    render() {
       return (
        <StripeProvider apiKey="pk_test_3KZy6O2l5hseQy9hglRwswYT">
            <div className="resp-tabs-container">
                <Elements>
                    <InjectedCheckoutForm history={this.props.history}/>
                </Elements>
            </div>
        </StripeProvider>
       );
    }
 }
 export default CreditCard;
 