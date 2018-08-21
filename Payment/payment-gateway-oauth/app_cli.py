import json

from sqlalchemy.exc import  SQLAlchemyError
from authlib.common.security import generate_token

from models.oauth import OAuth2Client
from models.base import db


def config_cli(app):
    @app.cli.command()
    def initdb():
        db.create_all()

    @app.cli.command()
    def addaclient():
        client = OAuth2Client(
            client_name="New Client",
            client_uri="http://localhost/",
            scope="profile",
            redirect_uri="http://localhost/",
            grant_type="password",
            response_type="code",
            token_endpoint_auth_method="client_secret_basic",
        )
        client.client_id = generate_token(24)
        client.client_secret = generate_token(48)
        try:
            db.session.add(client)
            db.session.commit()
        except SQLAlchemyError:
            print("DB commit failed")
            db.session.rollback()
        else:
            print("Client", client.client_id, client.client_secret)

        return app

