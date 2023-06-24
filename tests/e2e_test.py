import time
import requests
import pytest

api_url = 'http://localhost:8080'

@pytest.fixture(scope='session')
def session():
    with requests.Session() as session:
        yield session

@pytest.mark.usefixtures('session')
class TestAPI:
    def test_signup(self, session: requests.Session):
        response = session.post(f'{api_url}/signup', json={"gmail": "gaponukz@knu.ua"})
        assert response.status_code == 200

    def test_confirm_registration(self, session: requests.Session):
        time.sleep(1)
        response = session.post(f'{api_url}/confirmRegistration', json={
            "gmail": "gaponukz@knu.ua",
            "password": "somepass",
            "fullName": "Alex Yah",
            "phone": "380972748235",
            "key": "12345"
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_logout(self, session: requests.Session):
        response = session.get(f'{api_url}/logout')
        assert response.status_code == 200
        assert session.cookies.get("token") is None

    def test_signin_with_wrong_password(self, session: requests.Session):
        response = session.post(f'{api_url}/signin', json={
            "gmail": "gaponukz@knu.ua",
            "password": "wrongpassword",
            "rememberHim": True
        })
        assert response.status_code == 401

    def test_signin(self, session: requests.Session):
        response = session.post(f'{api_url}/signin', json={
            "gmail": "gaponukz@knu.ua",
            "password": "somepass",
            "rememberHim": True
        })
        assert response.status_code == 200
        assert session.cookies.get("token") is not None

    def test_get_user_info(self, session: requests.Session):
        response = session.get(f'{api_url}/getUserInfo')
        assert response.status_code == 200
        data = response.json()

        assert data['gmail'] == "gaponukz@knu.ua"
        assert data['phone'] == "380972748235"

    def test_subscribe_user_to_route(self, session: requests.Session):
        response = session.post(f'{api_url}/subscribeUserToTheRoute', json={
            "routeId": "g24g-h24hg2w-gh6j35w-w45g"
        })
        assert response.status_code == 200

    def test_get_user_info_after_subscription(self, session: requests.Session):
        response = session.get(f'{api_url}/getUserInfo')
        data = response.json()
        assert "g24g-h24hg2w-gh6j35w-w45g" in data['purchasedRouteIds']
