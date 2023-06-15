# JWT auth for Atass with Golang
<p align="center" width="100%">
    <img width="25%" src="https://github.com/gaponukz/atass-auth/assets/49754258/69b6c02a-2358-4e7c-89cd-d590a891905e">
</p>

This document provides an overview of the API endpoints and request bodies for the simple JWT implementation. <br/>
<b>Important:</b> The token is stored in a cookie, so there is no need to pass it in the headers etc.

## Endpoints

### Sign Up

- URL: `/signup`
- Method: `POST`
- Description: Starting registration process.
- Request Body:
```json
{
    "gmail": "user@example.com"
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `400` | Request Body is not correct |
| `500` | Could not send letter with code to user |

### Confirm Registration

- URL: `/confirmRegistration`
- Method: `POST`
- Description: Confirm user registration.
- Request Body:
```json
{
    "gmail": "user@example.com",
    "password": "somepass",
    "fullName": "Alex Yah",
    "phone": "380972748235",
    "rememberHim": true,
    "key": "906815"
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `400` | Request Body is not correct or it's been a long time |
| `500` | something went wrong while generating token, try signin |

### Sign In
- URL: `/signin`
- Method: `POST`
- Description: Signing in existing user.
- Request Body:
```json
{
    "gmail": "user@example.com",
    "password": "somepass",
    "rememberHim": true
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `401` | Request Body is not correct or user not found |
| `500` | something went wrong while generating token (very bad) |

### Refresh
- URL: `/refresh`
- Method: `GET`
- Description: Refresh token, it is desirable to call every 5 minutes.
- Response:

| Code | Description |
| :--- | :--- |
| `401` | Lost token in cookie |
| `400` | Something wrong with cookie, or it is just not time to refresh |
| `500` | something went wrong while generating new token (very bad) |

### Logout
- URL: `/logout`
- Method: `GET`

### Reset password.
- URL: `/resetPassword`
- Method: `POST`
- Description: Reset password for existing user.
- Request Body:
```json
{
    "gmail": "user@example.com",
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `400` | Request Body is not correct |
| `500` | Could not send letter with code to user |

### Confirm reset password.
- URL: `/confirmResetPassword`
- Method: `POST`
- Description: Reset password for existing user.
- Request Body:
```json
{
    "gmail": "user@example.com",
    "password": "somenewpassword",
    "key": "539991"
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `400` | Request Body is not correct or it's been a long time |
| `500` | something went wrong while generating token, try signin |

### Get user.
- URL: `/getUserInfo`
- Method: `GET`
- Description: Get user information in json.
- Response:
```json
{
    "gmail": "user@example.com",
    "phone": "3809734275232",
    "fullName": "Alex Yah",
    "rememberHim": true,
    "purchasedRouteIds": null
}
```

### Add route.
- URL: `/subscribeUserToTheRoute`
- Method: `POST`
- Description: Subscribe user to the route by id.
- Request Body:
```json
{
    "routeId": "g24g-h24hg2w-gh6j35w-w45g"
}
```
- Response:

| Code | Description |
| :--- | :--- |
| `400` | Request Body is not correct |
| `500` | something went wrong, try signin |

## Before start
### Settings
Before usage you need to create `.env` file:
```env
gmail=user@gmail.com
gmailPassword=userpassowrf123
jwtSecret=secret
```
### Dependencies
* Redis client
```bash
docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest
```
* Golang packages
```bash
go mod download
```
