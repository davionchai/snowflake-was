from datetime import datetime


class WarehouseStructs:
    name: str
    state: str
    type: str
    size: str
    running: int
    queued: int
    is_default: str
    is_current: str
    auto_suspend: int
    auto_resume: str
    available: str
    provisioning: str
    quiescing: str
    other: str
    created_on: datetime
    resumed_on: datetime
    updated_on: datetime
    owner: str
    comment: str
    resource_monitor: str
    actives: int
    pendings: int
    failed: int
    suspended: int
    uuid: str
