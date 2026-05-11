"""Flask application."""
from flask import Flask
from src.routes import api_bp


def create_app() -> Flask:
    """Application factory."""
    app = Flask(__name__)
    app.register_blueprint(api_bp)
    return app


if __name__ == '__main__':
    app = create_app()
    app.run(debug=True, port=5000)
