import time
import requests
import pytest
from user import User

user = User(
    gmail="gaponukz@knu.ua",
    password="qweasdzxcQ1",
    full_name="Alex Yah",
    phone="380984516456",
    allows_advertisement=True,
    purchased_route_ids=None,
)

api_url = 'http://localhost:8080'

@pytest.fixture(scope='session')
def session():
    with requests.Session() as session:
        yield session

@pytest.mark.usefixtures('session')
class TestAPI:
    def test_signup(self, session: requests.Session):
        response = session.post(f'{api_url}/api/auth/signup', json={"gmail": user.gmail})
        assert response.status_code == 200

    def test_confirm_registration(self, session: requests.Session):
        time.sleep(1)
        json = user.to_json()
        
        json['key'] = "wrongcode"
        response = session.post(f'{api_url}/api/auth/confirmRegistration', json=json)
        assert response.status_code != 200

        json['key'] = "12345"
        response = session.post(f'{api_url}/api/auth/confirmRegistration', json=json)
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_logout(self, session: requests.Session):
        response = session.get(f'{api_url}/api/auth/logout')
        assert response.status_code == 200
        assert session.cookies.get("token") is None

    def test_signin_with_wrong_password(self, session: requests.Session):
        response = session.post(f'{api_url}/api/auth/signin', json={
            "gmail": user.gmail,
            "password": "wrongpassword",
            "rememberHim": True
        })
        assert response.status_code == 401

    def test_signin(self, session: requests.Session):
        response = session.post(f'{api_url}/api/auth/signin', json={
            "gmail": user.gmail,
            "password": user.password,
            "rememberHim": True
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_get_user_info(self, session: requests.Session):
        response = session.get(f'{api_url}/api/auth/getUserInfo')
        assert response.status_code == 200
        data = response.json()

        assert data['gmail'] == user.gmail
        assert data['phone'] == user.phone
        assert data['fullName'] == user.full_name
        assert data['allowsAdvertisement'] == user.allows_advertisement
    
    def test_reset_password(self, session: requests.Session):
        response = session.post(f'{api_url}/api/auth/resetPassword', json={"gmail": user.gmail})
        assert response.status_code == 200

        user.password = "P@ssw0rd123"
    
        response = session.post(f'{api_url}/api/auth/confirmResetPassword', json={
            "gmail": user.gmail,
            "password": user.password,
            "key": "wrongcode"
        })

        assert response.status_code != 200

        response = session.post(f'{api_url}/api/auth/confirmResetPassword', json={
            "gmail": user.gmail,
            "password": user.password,
            "key": "12345"
        })

        assert response.status_code == 200

    def test_login_with_new_password(self, session: requests.Session):
        response = session.post(f'{api_url}/api/auth/signin', json={
            "gmail": user.gmail,
            "password": user.password,
            "rememberHim": True
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_change_user_info(self, session: requests.Session):
        user.full_name = "Max Euq"
        user.phone = "0984516456"
        user.allows_advertisement = not user.allows_advertisement

        response = session.post(f'{api_url}/api/auth/updateUserInfo', json={
            "fullName": user.full_name,
            "phone": user.phone,
            "allowsAdvertisement": user.allows_advertisement
        })

        assert response.status_code == 200

        response = session.get(f'{api_url}/api/auth/getUserInfo')
        data = response.json()

        assert data['fullName'] == user.full_name
        assert data['phone'] == user.phone
        assert data['allowsAdvertisement'] == user.allows_advertisement
