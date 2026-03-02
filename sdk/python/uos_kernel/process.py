"""Process — Base class for business processes, mirroring Go app.Process."""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from uos_kernel.client import KernelClient


class Process(ABC):
    """Base class for business processes running on top of the Kernel API."""

    @abstractmethod
    def name(self) -> str:
        """Return the process name."""
        ...

    @abstractmethod
    def run(self, client: KernelClient) -> None:
        """Run the process. Blocks until complete or cancelled."""
        ...
