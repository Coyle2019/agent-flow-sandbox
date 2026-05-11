"""Tests for registration endpoint."""
import pytest
from src.app import create_app
from src.models.user import User
from src.utils.auth import hash_password, verify_password, create_access_token


@pytest.fixture
def client():
    """Create test client."""
    app = create_app()
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client


class TestRegisterEndpoint:
    """Test POST /api/v1/register."""

    def test_register_success(self, client):
        """Test successful registration returns 201."""
        response = client.post('/api/v1/register', json={
            'username': 'testuser',
            'email': 'test@example.com',
            'password': 'password123'
        })
        assert response.status_code == 201
        data = response.get_json()
        assert 'user_id' in data
        assert 'access_token' in data
        assert data['username'] == 'testuser'
        assert data['email'] == 'test@example.com'

    def test_register_duplicate_email(self, client):
        """Test duplicate email returns 409."""
        client.post('/api/v1/register', json={
            'username': 'user1',
            'email': 'duplicate@example.com',
            'password': 'password123'
        })
        response = client.post('/api/v1/register', json={
            'username': 'user2',
            'email': 'duplicate@example.com',
            'password': 'password456'
        })
        assert response.status_code == 409
        assert response.get_json()['error'] == '该邮箱已被注册'

    def test_register_invalid_email(self, client):
        """Test invalid email format returns 400."""
        response = client.post('/api/v1/register', json={
            'username': 'testuser',
            'email': 'invalid-email',
            'password': 'password123'
        })
        assert response.status_code == 400
        assert response.get_json()['error'] == '请输入有效的邮箱地址'

    def test_register_short_password(self, client):
        """Test password less than 8 chars returns 400."""
        response = client.post('/api/v1/register', json={
            'username': 'testuser',
            'email': 'test@example.com',
            'password': 'short'
        })
        assert response.status_code == 400


class TestAuthUtils:
    """Test authentication utilities."""

    def test_password_hash_and_verify(self):
        """Test password hashing and verification."""
        password = 'testpassword123'
        hashed = hash_password(password)
        assert hashed != password
        assert verify_password(password, hashed)
        assert not verify_password('wrongpassword', hashed)

    def test_create_access_token(self):
        """Test JWT token creation."""
        token = create_access_token('user123')
        assert isinstance(token, str)
        assert token.count('.') == 2
