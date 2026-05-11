"""Registration schemas."""
from dataclasses import dataclass
import re


EMAIL_REGEX = re.compile(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$')


@dataclass
class RegisterRequest:
    """Registration request schema."""
    username: str
    email: str
    password: str

    def validate(self) -> tuple[bool, str | None]:
        """Validate request data."""
        if not self.username or not self.username.strip():
            return False, "请填写所有必填项"
        if not self.email or not EMAIL_REGEX.match(self.email):
            return False, "请输入有效的邮箱地址"
        if not self.password or len(self.password) < 8:
            return False, "请设置密码（至少8位）"
        return True, None


@dataclass
class RegisterResponse:
    """Registration response schema."""
    user_id: str
    username: str
    email: str
    access_token: str

    def to_dict(self) -> dict:
        """Convert to dictionary for JSON response."""
        return {
            "user_id": self.user_id,
            "username": self.username,
            "email": self.email,
            "access_token": self.access_token
        }
