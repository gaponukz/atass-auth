import time
import requests
import pytest
from user import User

user = User(
    gmail="gaponukz@knu.ua",
    password="somepass",
    full_name="Alex Yah",
    phone="380972748235",
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
        response = session.post(f'{api_url}/signup', json={"gmail": user.gmail})
        assert response.status_code == 200

    def test_confirm_registration(self, session: requests.Session):
        time.sleep(1)
        json = user.to_json()
        json['key'] = "12345"

        response = session.post(f'{api_url}/confirmRegistration', json=json)
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_logout(self, session: requests.Session):
        response = session.get(f'{api_url}/logout')
        assert response.status_code == 200
        assert session.cookies.get("token") is None

    def test_signin_with_wrong_password(self, session: requests.Session):
        response = session.post(f'{api_url}/signin', json={
            "gmail": user.gmail,
            "password": "wrongpassword",
            "rememberHim": True
        })
        assert response.status_code == 401

    def test_signin(self, session: requests.Session):
        response = session.post(f'{api_url}/signin', json={
            "gmail": user.gmail,
            "password": user.password,
            "rememberHim": True
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_get_user_info(self, session: requests.Session):
        response = session.get(f'{api_url}/getUserInfo')
        assert response.status_code == 200
        data = response.json()

        assert data['gmail'] == user.gmail
        assert data['phone'] == user.phone
        assert data['fullName'] == user.full_name
        assert data['allowsAdvertisement'] == user.allows_advertisement

    def test_subscribe_user_to_route(self, session: requests.Session):
        response = session.post(f'{api_url}/subscribeUserToTheRoute', json={
            "routeId": "g24g-h24hg2w-gh6j35w-w45g"
        })
        
        if user.purchased_route_ids is None:
            user.purchased_route_ids = []
        
        user.purchased_route_ids.append("g24g-h24hg2w-gh6j35w-w45g")
        assert response.status_code == 200

    def test_get_user_info_after_subscription(self, session: requests.Session):
        response = session.get(f'{api_url}/getUserInfo')
        data = response.json()
        assert "g24g-h24hg2w-gh6j35w-w45g" in data['purchasedRouteIds']
    
    def test_reset_password(self, session: requests.Session):
        response = session.post(f'{api_url}/resetPassword', json={"gmail": user.gmail})
        assert response.status_code == 200

        user.password = "newpassword"
        response = session.post(f'{api_url}/confirmResetPassword', json={
            "gmail": user.gmail,
            "password": user.password,
            "key": "12345"
        })

        assert response.status_code == 200

    def test_login_with_new_password(self, session: requests.Session):
        response = session.post(f'{api_url}/signin', json={
            "gmail": user.gmail,
            "password": user.password,
            "rememberHim": True
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_change_name(self, session: requests.Session):
        user.full_name = "Max Euq"
        response = session.post(f'{api_url}/updateName', json={
            "fullName": user.full_name
        })

        assert response.status_code == 200

        response = session.get(f'{api_url}/getUserInfo')
        data = response.json()

        assert data['fullName'] == user.full_name
    
    def test_change_phone(self, session: requests.Session):
        user.phone = "3801234567"
        response = session.post(f'{api_url}/updatePhone', json={
            "phone": user.phone
        })

        assert response.status_code == 200

        response = session.get(f'{api_url}/getUserInfo')
        data = response.json()

        assert data['phone'] == user.phone
    
    def test_change_allows_advertisement(self, session: requests.Session):
        user.allows_advertisement = not user.allows_advertisement
        response = session.post(f'{api_url}/updateAllowsAdvertisement', json={
            "allowsAdvertisement": user.allows_advertisement
        })

        assert response.status_code == 200

        response = session.get(f'{api_url}/getUserInfo')
        data = response.json()

        assert data['allowsAdvertisement'] == user.allows_advertisement