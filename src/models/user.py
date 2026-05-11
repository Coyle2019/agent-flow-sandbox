"""User model."""
from dataclasses import dataclass
from datetime import datetime
import uuid


@dataclass
class User:
    """User data model."""
    id: str
    username: str
    email: str
    password_hash: str
    created_at: datetime

    @classmethod
    def create(cls, username: str, email: str, password_hash: str) -> "User":
        """Factory method to create a new user."""
        return cls(
            id=str(uuid.uuid4()),
            username=username,
            email=email,
            password_hash=password_hash,
            created_at=datetime.utcnow()
        )

    def to_dict(self) -> dict:
        """Convert to dictionary (excluding password_hash)."""
        return {
            "id": self.id,
            "username": self.username,
            "email": self.email,
            "created_at": self.created_at.isoformat()
        }
