from __future__ import annotations
import dataclasses

@dataclasses.dataclass
class User:
    gmail: str
    password: str
    full_name: str
    phone: str
    allows_advertisement: bool
    purchased_route_ids: list[str] | None = None

    def from_json(self, data: dict) -> User:
        return User(
            gmail=data['gmail'],
            password=data['password'],
            full_name=data['fullName'],
            phone=data['phone'],
            allows_advertisement=data['allowsAdvertisement'],
            purchased_route_ids=data.get('purchasedRouteIds'),
        )

    def to_json(self) -> dict:
        return {
            'gmail': self.gmail,
            'password': self.password,
            'fullName': self.full_name,
            'phone': self.phone,
            'allowsAdvertisement': self.allows_advertisement,
            'purchasedRouteIds': self.purchased_route_ids
        }
