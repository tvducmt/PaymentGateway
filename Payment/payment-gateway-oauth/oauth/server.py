from authlib.flask.oauth2 import AuthorizationServer
from authlib.flask.oauth2.sqla import (
    create_query_client_func,
    create_save_token_func,
)

from models.oauth import db
from models.oauth import OAuth2Client, OAuth2Token
from oauth.grant import PasswordGrant


authorization_server = AuthorizationServer(
    query_client=create_query_client_func(db.session, OAuth2Client),
    save_token=create_save_token_func(db.session, OAuth2Token),
)

def config_oauth(app):
    authorization_server.init_app(app)
    # authorization_server.register_grant(AuthorizationCodeGrant)
    authorization_server.register_grant(PasswordGrant)
