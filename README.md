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
- Description: Creates a new user account.
- Request Body:
```json
{
    "gmail": "user@example.com",
    "password": "somepass",
    "fullName": "Alex Yah",
    "phone": "380972748235",
    "rememberHim": true
}
```
### Confirm Registration

- URL: `/confirm`
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
    "uniqueKey": "906815"
}
```

### Sign In
- URL: `/signin`
- Method: `POST`
- Description: Signing in existing user.
- Request Body:
```json
{
    "gmail": "user@example.com",
    "password": "somepass"
}
```

### Refresh
- URL: `/refresh`
- Method: `GET`
- Description: Refresh token, it is desirable to call every 5 minutes.

### Logout
- URL: `/logout`
- Method: `GET`

### Welcome (for test)
- URL: `/welcome`
- Method: `GET`
- Description: returns some user info.

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
