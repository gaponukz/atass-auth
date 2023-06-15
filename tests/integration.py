import time
import requests

api_url = 'http://localhost:8080'

with requests.Session() as session:
    response = session.post(f'{api_url}/signup', json={"gmail": "gaponukz@knu.ua"})
    assert response.status_code == 200

    time.sleep(1)
    response = session.post(f'{api_url}/confirmRegistration', json={
        "gmail": "gaponukz@knu.ua",
        "password": "somepass",
        "fullName": "Alex Yah",
        "phone": "380972748235",
        "key": "12345"
    })

    assert response.status_code == 200
    assert session.cookies.get("token") != None

    response = session.get(f'{api_url}/logout')
    assert response.status_code == 200
    assert session.cookies.get("token") == None

    response = session.post(f'{api_url}/signin', json={
        "gmail": "gaponukz@knu.ua",
        "password": "somepass",
        "rememberHim": True
    })

    assert response.status_code == 200
    assert session.cookies.get("token") != None

    response = session.get(f'{api_url}/getUserInfo')
    assert response.status_code == 200
    data = response.json()

    assert data['gmail'] == "gaponukz@knu.ua"
    assert data['phone'] == "380972748235"

    response = session.post(f'{api_url}/subscribeUserToTheRoute', json={
        "routeId": "g24g-h24hg2w-gh6j35w-w45g"
    })
    assert response.status_code == 200

    response = session.get(f'{api_url}/getUserInfo')
    data = response.json()
    assert "g24g-h24hg2w-gh6j35w-w45g" in data['purchasedRouteIds']

    