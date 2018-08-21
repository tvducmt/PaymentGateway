import bcrypt

from models.oauth import db

class User(db.Model):
    __tablename__ = 'account'

    id = db.Column(db.Integer, primary_key=True)
    username = db.Column("email", db.String, unique=True, nullable=False)
    password = db.Column("passphrase", db.String, nullable=False)

    def __str__(self):
        return "<User {}>".format(self.username)

    def get_user_id(self):
        return self.id

    def check_password(self, password):
        return bcrypt.checkpw(password.encode('utf8'), self.password.encode('utf8'))
