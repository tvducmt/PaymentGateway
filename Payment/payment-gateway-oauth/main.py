import os

from flask import Flask
from dotenv import load_dotenv

from models.base import db
from routes.auth import bp as auth_bp
from routes.index import bp as index_bp

from oauth.server import config_oauth
from app_cli import config_cli

load_dotenv()

def create_app():
    app = Flask(__name__)
    app.config.update({
        'SECRET_KEY': os.getenv("FLASK_SECRET"),
        'OAUTH2_REFRESH_TOKEN_GENERATOR': True,
        'SQLALCHEMY_TRACK_MODIFICATIONS': False,
        'SQLALCHEMY_DATABASE_URI': os.getenv("SQLALCHEMY_DATABASE_URI"),
    })

    db.init_app(app)
    config_oauth(app)
    config_cli(app)
    app.register_blueprint(index_bp)
    app.register_blueprint(auth_bp, url_prefix="/auth")

    return app
