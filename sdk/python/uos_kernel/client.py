"""KernelClient — Pythonic wrapper for the gRPC KernelService."""

from __future__ import annotations

import sys
import os
from typing import Iterator

import grpc
from google.protobuf.json_format import MessageToDict, ParseDict
from google.protobuf.struct_pb2 import Struct

# Add generated code to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "gen"))

from uos.kernel.v1 import kernel_pb2 as pb
from uos.kernel.v1 import kernel_pb2_grpc as pb_grpc

# Proto EventType enum → string mapping
_PROTO_TO_EVENT_TYPE = {
    pb.EVENT_TYPE_STATE_CHANGE: "state_change",
    pb.EVENT_TYPE_ATTRIBUTE_CHANGE: "attribute_change",
    pb.EVENT_TYPE_CAPABILITY_CHANGE: "capability_change",
    pb.EVENT_TYPE_CUSTOM: "custom",
}

_EVENT_TYPE_TO_PROTO = {v: k for k, v in _PROTO_TO_EVENT_TYPE.items()}


class KernelClient:
    """Wraps the gRPC KernelService stub with a Pythonic API."""

    def __init__(self, address: str = "localhost:50051"):
        self._channel = grpc.insecure_channel(address)
        self._stub = pb_grpc.KernelServiceStub(self._channel)

    def open(self, resource_type: str, resource_id: str, flags: int = 0) -> int:
        """Open a resource. Returns file descriptor (int)."""
        resp = self._stub.Open(pb.OpenRequest(
            resource_type=resource_type,
            resource_id=resource_id,
            flags=flags,
        ))
        return resp.fd

    def close(self, fd: int) -> None:
        """Close a resource descriptor."""
        self._stub.Close(pb.CloseRequest(fd=fd))

    def read(self, fd: int) -> dict:
        """Read resource state. Returns dict with resource_id, resource_type, capabilities, properties."""
        resp = self._stub.Read(pb.ReadRequest(fd=fd))
        props = {}
        if resp.HasField("properties"):
            props = MessageToDict(resp.properties)
        return {
            "resource_id": resp.resource_id,
            "resource_type": resp.resource_type,
            "capabilities": list(resp.capabilities),
            "properties": props,
        }

    def write(self, fd: int, updates: dict) -> None:
        """Write property updates to a resource."""
        s = Struct()
        s.update(updates)
        self._stub.Write(pb.WriteRequest(fd=fd, updates=s))

    def stat(self, fd: int) -> dict:
        """Get resource stat information."""
        resp = self._stub.Stat(pb.StatRequest(fd=fd))
        return {
            "resource_id": resp.resource_id,
            "resource_type": resp.resource_type,
            "capabilities": [
                {"name": c.name, "operations": list(c.operations), "description": c.description}
                for c in resp.capabilities
            ],
            "events": [
                {"name": e.name, "event_type": e.event_type, "description": e.description}
                for e in resp.events
            ],
        }

    def ioctl(self, fd: int, request: int, args: dict | None = None) -> dict | None:
        """Execute ioctl command. Args dict is sent as protobuf Struct."""
        req = pb.IoctlRequest(fd=fd, request=request)
        if args is not None:
            s = Struct()
            s.update(args)
            req.args.CopyFrom(s)
        resp = self._stub.Ioctl(req)
        if resp.HasField("result"):
            return MessageToDict(resp.result)
        return None

    def watch(self, fd: int, event_types: list[str] | None = None) -> Iterator[dict]:
        """Watch for events. Returns an iterator of event dicts.

        Each event dict has keys: type (str), resource_id (str), data (dict).
        Custom events include a '_type' field in data for type discrimination.
        """
        proto_types = []
        for et in (event_types or []):
            if et in _EVENT_TYPE_TO_PROTO:
                proto_types.append(_EVENT_TYPE_TO_PROTO[et])

        stream = self._stub.Watch(pb.WatchRequest(fd=fd, event_types=proto_types))
        for event in stream:
            data = {}
            if event.HasField("data"):
                data = MessageToDict(event.data)
            yield {
                "type": _PROTO_TO_EVENT_TYPE.get(event.type, "unknown"),
                "resource_id": event.resource_id,
                "data": data,
            }

    def __enter__(self):
        return self

    def __exit__(self, *args):
        self._channel.close()
