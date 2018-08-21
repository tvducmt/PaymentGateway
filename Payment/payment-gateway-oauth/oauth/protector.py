from flask import jsonify

from authlib.flask.oauth2 import ResourceProtector, current_token
from authlib.flask.oauth2.sqla import create_bearer_token_validator
from authlib.specs.rfc6750 import BearerTokenValidator

from models.oauth import OAuth2Token


class MyBearerTokenValidator(BearerTokenValidator):
    def authenticate_token(self, token_string):
        return OAuth2Token.query.filter_by(access_token=token_string).first()

    def request_invalid(self, request):
        return False

    def token_revoked(self, token):
        return token.revoked

ResourceProtector.register_token_validator(MyBearerTokenValidator())

require_oauth = ResourceProtector()
