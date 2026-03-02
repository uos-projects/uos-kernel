"""UOS Kernel Python SDK"""

# Open flags
O_RDONLY = 0x0000
O_WRONLY = 0x0001
O_RDWR = 0x0002
O_CREAT = 0x0200
O_EXCL = 0x0800

# Ioctl commands
CMD_GET_RESOURCE_INFO = 0x1000
CMD_LIST_CAPABILITIES = 0x1001
CMD_LIST_EVENTS = 0x1002
CMD_EXECUTE_CAPACITY = 0x1003
CMD_SYNC = 0x2001

# Event types
EVENT_STATE_CHANGE = "state_change"
EVENT_ATTRIBUTE_CHANGE = "attribute_change"
EVENT_CAPABILITY_CHANGE = "capability_change"
EVENT_CUSTOM = "custom"

from uos_kernel.client import KernelClient
from uos_kernel.process import Process
from uos_kernel.helpers import wait_for_property, wait_for_event, execute_capacity

__all__ = [
    "KernelClient",
    "Process",
    "wait_for_property",
    "wait_for_event",
    "execute_capacity",
    "O_RDONLY", "O_WRONLY", "O_RDWR", "O_CREAT", "O_EXCL",
    "CMD_GET_RESOURCE_INFO", "CMD_LIST_CAPABILITIES", "CMD_LIST_EVENTS",
    "CMD_EXECUTE_CAPACITY", "CMD_SYNC",
    "EVENT_STATE_CHANGE", "EVENT_ATTRIBUTE_CHANGE",
    "EVENT_CAPABILITY_CHANGE", "EVENT_CUSTOM",
]
