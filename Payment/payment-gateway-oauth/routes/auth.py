from flask import Blueprint, render_template, request

from authlib.specs.rfc6749 import OAuth2Error

from oauth.server import authorization_server
from models.user import User


bp = Blueprint("auth", __name__)

"""
    /oauth/authorize?blah=blah&blah=blah%blah=blah
    ...
    <p>{{grant.client.client_name}} is requesting:
    <strong>{{ grant.request.scope }}</strong>
    ...
    'client_id', 'code', 'redirect_uri', 'scope', 'state',
                'response_type', 'grant_type'
"""
@bp.route('/authorize')
def authorize_get():
    try:
        grant = authorization_server.validate_consent_request(request=request)
    except OAuth2Error as error:
        return "OAuth2 Error: {}".format(error), 406

    return render_template('authorize.html', grant=grant)

# ...
@bp.route('/authorize', methods=["POST"])
def authorize_post():
    try:
        username = request.form["username"]
        password = request.form["password"]
    except KeyError:
        return "Username or Password not provived", 403

    user = User.query.filter_by(username=username).first()
    if user is None or not user.check_password(password):
        return "Wrong username or password", 403
    return authorization_server.create_authorization_response(grant_user=user)

"""
    curl -u HMlbWXVyYHjosQHhLy8jxeWs:i4m5YaCNLvZqZi6xMpF27YSsHhn8j0ml7mavBP5PJQDOvqEE -XPOST http://127.0.0.1:5000/auth/token -F grant_type=password -F username=stev -F password=hunter2 -F scope=profile
"""
@bp.route('/token', methods=['POST'])
def issue_token():
    return authorization_server.create_token_response()

"""
    https://stackoverflow.com/a/49880763
    curl -H "Content-type:application/x-www-form-urlencoded" https://accounts.google.com/o/oauth2/revoke?token={token}
"""
@bp.route('/revoke', methods=['POST'])
def revoke_token():
    return authorization_server.create_endpoint_response('revocation')
