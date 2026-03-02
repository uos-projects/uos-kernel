from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EventType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    EVENT_TYPE_UNSPECIFIED: _ClassVar[EventType]
    EVENT_TYPE_STATE_CHANGE: _ClassVar[EventType]
    EVENT_TYPE_ATTRIBUTE_CHANGE: _ClassVar[EventType]
    EVENT_TYPE_CAPABILITY_CHANGE: _ClassVar[EventType]
    EVENT_TYPE_CUSTOM: _ClassVar[EventType]
EVENT_TYPE_UNSPECIFIED: EventType
EVENT_TYPE_STATE_CHANGE: EventType
EVENT_TYPE_ATTRIBUTE_CHANGE: EventType
EVENT_TYPE_CAPABILITY_CHANGE: EventType
EVENT_TYPE_CUSTOM: EventType

class OpenRequest(_message.Message):
    __slots__ = ("resource_type", "resource_id", "flags")
    RESOURCE_TYPE_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    FLAGS_FIELD_NUMBER: _ClassVar[int]
    resource_type: str
    resource_id: str
    flags: int
    def __init__(self, resource_type: _Optional[str] = ..., resource_id: _Optional[str] = ..., flags: _Optional[int] = ...) -> None: ...

class OpenResponse(_message.Message):
    __slots__ = ("fd",)
    FD_FIELD_NUMBER: _ClassVar[int]
    fd: int
    def __init__(self, fd: _Optional[int] = ...) -> None: ...

class CloseRequest(_message.Message):
    __slots__ = ("fd",)
    FD_FIELD_NUMBER: _ClassVar[int]
    fd: int
    def __init__(self, fd: _Optional[int] = ...) -> None: ...

class CloseResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ReadRequest(_message.Message):
    __slots__ = ("fd",)
    FD_FIELD_NUMBER: _ClassVar[int]
    fd: int
    def __init__(self, fd: _Optional[int] = ...) -> None: ...

class ReadResponse(_message.Message):
    __slots__ = ("resource_id", "resource_type", "capabilities", "properties")
    RESOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_TYPE_FIELD_NUMBER: _ClassVar[int]
    CAPABILITIES_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    resource_id: str
    resource_type: str
    capabilities: _containers.RepeatedScalarFieldContainer[str]
    properties: _struct_pb2.Struct
    def __init__(self, resource_id: _Optional[str] = ..., resource_type: _Optional[str] = ..., capabilities: _Optional[_Iterable[str]] = ..., properties: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class WriteRequest(_message.Message):
    __slots__ = ("fd", "updates")
    FD_FIELD_NUMBER: _ClassVar[int]
    UPDATES_FIELD_NUMBER: _ClassVar[int]
    fd: int
    updates: _struct_pb2.Struct
    def __init__(self, fd: _Optional[int] = ..., updates: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class WriteResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class StatRequest(_message.Message):
    __slots__ = ("fd",)
    FD_FIELD_NUMBER: _ClassVar[int]
    fd: int
    def __init__(self, fd: _Optional[int] = ...) -> None: ...

class CapabilityInfo(_message.Message):
    __slots__ = ("name", "operations", "description")
    NAME_FIELD_NUMBER: _ClassVar[int]
    OPERATIONS_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    name: str
    operations: _containers.RepeatedScalarFieldContainer[str]
    description: str
    def __init__(self, name: _Optional[str] = ..., operations: _Optional[_Iterable[str]] = ..., description: _Optional[str] = ...) -> None: ...

class EventInfo(_message.Message):
    __slots__ = ("name", "event_type", "description")
    NAME_FIELD_NUMBER: _ClassVar[int]
    EVENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    name: str
    event_type: str
    description: str
    def __init__(self, name: _Optional[str] = ..., event_type: _Optional[str] = ..., description: _Optional[str] = ...) -> None: ...

class StatResponse(_message.Message):
    __slots__ = ("resource_id", "resource_type", "capabilities", "events")
    RESOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_TYPE_FIELD_NUMBER: _ClassVar[int]
    CAPABILITIES_FIELD_NUMBER: _ClassVar[int]
    EVENTS_FIELD_NUMBER: _ClassVar[int]
    resource_id: str
    resource_type: str
    capabilities: _containers.RepeatedCompositeFieldContainer[CapabilityInfo]
    events: _containers.RepeatedCompositeFieldContainer[EventInfo]
    def __init__(self, resource_id: _Optional[str] = ..., resource_type: _Optional[str] = ..., capabilities: _Optional[_Iterable[_Union[CapabilityInfo, _Mapping]]] = ..., events: _Optional[_Iterable[_Union[EventInfo, _Mapping]]] = ...) -> None: ...

class IoctlRequest(_message.Message):
    __slots__ = ("fd", "request", "args")
    FD_FIELD_NUMBER: _ClassVar[int]
    REQUEST_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    fd: int
    request: int
    args: _struct_pb2.Struct
    def __init__(self, fd: _Optional[int] = ..., request: _Optional[int] = ..., args: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class IoctlResponse(_message.Message):
    __slots__ = ("result",)
    RESULT_FIELD_NUMBER: _ClassVar[int]
    result: _struct_pb2.Struct
    def __init__(self, result: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class WatchRequest(_message.Message):
    __slots__ = ("fd", "event_types")
    FD_FIELD_NUMBER: _ClassVar[int]
    EVENT_TYPES_FIELD_NUMBER: _ClassVar[int]
    fd: int
    event_types: _containers.RepeatedScalarFieldContainer[EventType]
    def __init__(self, fd: _Optional[int] = ..., event_types: _Optional[_Iterable[_Union[EventType, str]]] = ...) -> None: ...

class WatchEvent(_message.Message):
    __slots__ = ("type", "resource_id", "data")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    type: EventType
    resource_id: str
    data: _struct_pb2.Struct
    def __init__(self, type: _Optional[_Union[EventType, str]] = ..., resource_id: _Optional[str] = ..., data: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...
