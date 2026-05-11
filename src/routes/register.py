"""Registration endpoint."""
from flask import request, jsonify
from src.routes import api_bp
from src.models.user import User
from src.schemas.register import RegisterRequest, RegisterResponse
from src.utils.auth import hash_password, create_access_token

# In-memory user store (replace with database in production)
_users_db: dict[str, User] = {}
_email_index: set[str] = set()


@api_bp.route('/register', methods=['POST'])
def register_user():
    """Handle user registration."""
    data = request.get_json()
    if not data:
        return jsonify({"error": "请填写所有必填项"}), 400

    req = RegisterRequest(
        username=data.get('username', ''),
        email=data.get('email', ''),
        password=data.get('password', '')
    )

    valid, error_msg = req.validate()
    if not valid:
        return jsonify({"error": error_msg}), 400

    if req.email in _email_index:
        return jsonify({"error": "该邮箱已被注册"}), 409

    password_hash = hash_password(req.password)
    user = User.create(username=req.username, email=req.email, password_hash=password_hash)

    _users_db[user.id] = user
    _email_index.add(user.email)

    access_token = create_access_token(user.id)
    resp = RegisterResponse(
        user_id=user.id,
        username=user.username,
        email=user.email,
        access_token=access_token
    )

    return jsonify(resp.to_dict()), 201
