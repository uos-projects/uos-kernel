"""Helper functions for common Kernel API patterns."""

from __future__ import annotations

import time
from typing import Callable, TYPE_CHECKING

if TYPE_CHECKING:
    from uos_kernel.client import KernelClient

from uos_kernel import CMD_EXECUTE_CAPACITY, EVENT_ATTRIBUTE_CHANGE, EVENT_CUSTOM


def wait_for_property(
    client: KernelClient,
    fd: int,
    prop: str,
    expected,
    timeout: float = 30.0,
) -> None:
    """Wait until a resource property reaches the expected value.

    First checks current value via Read. If not matched, uses Watch
    to wait for attribute_change events until the property matches.

    Raises TimeoutError if the expected value is not reached within timeout.
    """
    # Check current value first
    state = client.read(fd)
    if state["properties"].get(prop) == expected:
        return

    # Watch for changes
    deadline = time.monotonic() + timeout
    for event in client.watch(fd, [EVENT_ATTRIBUTE_CHANGE]):
        if time.monotonic() > deadline:
            raise TimeoutError(f"Timed out waiting for {prop} == {expected}")
        data = event.get("data", {})
        if data.get("name") == prop and data.get("value") == expected:
            return


def wait_for_event(
    client: KernelClient,
    fd: int,
    match: Callable[[dict], bool],
    timeout: float = 30.0,
) -> dict:
    """Wait for a matching custom event.

    match: a callable that takes an event dict and returns True if it matches.
    Returns the matching event dict.

    Raises TimeoutError if no matching event arrives within timeout.
    """
    deadline = time.monotonic() + timeout
    for event in client.watch(fd, [EVENT_CUSTOM]):
        if time.monotonic() > deadline:
            raise TimeoutError("Timed out waiting for matching event")
        if match(event):
            return event
    raise RuntimeError("Watch stream closed before matching event")


def execute_capacity(
    client: KernelClient,
    fd: int,
    message_type: str,
    fields: dict | None = None,
    capacity: str | None = None,
) -> dict | None:
    """Execute a capacity via Ioctl CMD_EXECUTE_CAPACITY.

    Args:
        client: KernelClient instance
        fd: resource file descriptor
        message_type: the _type discriminator (e.g., "OpenBreakerCommand")
        fields: command fields as a dict
        capacity: optional capacity name

    Returns:
        Ioctl result dict, or None.
    """
    args = {"_type": message_type}
    if fields:
        args.update(fields)
    if capacity:
        args["capacity"] = capacity
    return client.ioctl(fd, CMD_EXECUTE_CAPACITY, args)
