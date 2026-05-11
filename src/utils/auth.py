"""Authentication utilities."""
from datetime import datetime, timedelta
import hashlib
import hmac
import base64
import json


def hash_password(password: str) -> str:
    """Hash password using SHA-256 with salt."""
    import secrets
    salt = secrets.token_hex(16)
    pwd_hash = hashlib.pbkdf2_hmac('sha256', password.encode(), salt.encode(), 100000)
    return f"{salt}${base64.b64encode(pwd_hash).decode()}"


def verify_password(password: str, password_hash: str) -> bool:
    """Verify password against stored hash."""
    try:
        salt, pwd_hash_b64 = password_hash.split('$')
        pwd_hash = base64.b64decode(pwd_hash_b64)
        calc_hash = hashlib.pbkdf2_hmac('sha256', password.encode(), salt.encode(), 100000)
        return hmac.compare_digest(pwd_hash, calc_hash)
    except (ValueError, Exception):
        return False


def create_access_token(user_id: str, expires_delta: timedelta = timedelta(hours=24)) -> str:
    """Create JWT-like access token."""
    header = {"alg": "HS256", "typ": "JWT"}
    payload = {
        "sub": user_id,
        "iat": datetime.utcnow(),
        "exp": datetime.utcnow() + expires_delta
    }
    message = base64.urlsafe_b64encode(json.dumps(header).encode()).decode().rstrip('=')
    message += "." + base64.urlsafe_b64encode(json.dumps(payload).encode()).decode().rstrip('=')
    signature = base64.urlsafe_b64encode(
        hmac.new(b"secret_key", message.encode(), hashlib.sha256).digest()
    ).decode().rstrip('=')
    return f"{message}.{signature}"
