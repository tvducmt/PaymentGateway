from flask import Blueprint, jsonify
from authlib.flask.oauth2 import current_token

from oauth.protector import require_oauth

bp = Blueprint("index", __name__)

@bp.route("/")
def index():
    return "OK"

@bp.route("/verify")
@require_oauth()
def me():
    user = current_token.user
    return jsonify({
        "id": user.id,
        "username": user.username
    })
