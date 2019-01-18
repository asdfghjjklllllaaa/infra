# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: testmessage.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='testmessage.proto',
  package='bigquery',
  syntax='proto3',
  serialized_pb=_b('\n\x11testmessage.proto\x12\x08\x62igquery\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\x82\x04\n\x0bTestMessage\x12\x0b\n\x03str\x18\x01 \x01(\t\x12\x0c\n\x04strs\x18\x02 \x03(\t\x12\x0b\n\x03num\x18\x03 \x01(\x03\x12\x0c\n\x04nums\x18\x04 \x03(\x03\x12\x16\n\x01\x65\x18\x05 \x01(\x0e\x32\x0b.bigquery.E\x12\x17\n\x02\x65s\x18\x06 \x03(\x0e\x32\x0b.bigquery.E\x12\'\n\x06nested\x18\x07 \x01(\x0b\x32\x17.bigquery.NestedMessage\x12(\n\x07nesteds\x18\x08 \x03(\x0b\x32\x17.bigquery.NestedMessage\x12%\n\x05\x65mpty\x18\t \x01(\x0b\x32\x16.google.protobuf.Empty\x12\'\n\x07\x65mpties\x18\n \x03(\x0b\x32\x16.google.protobuf.Empty\x12\'\n\x06struct\x18\x0b \x01(\x0b\x32\x17.google.protobuf.Struct\x12(\n\x07structs\x18\x0c \x03(\x0b\x32\x17.google.protobuf.Struct\x12-\n\ttimestamp\x18\r \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12.\n\ntimestamps\x18\x0e \x03(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x37\n\x12repeated_container\x18\x0f \x01(\x0b\x32\x1b.bigquery.RepeatedContainer\")\n\rNestedMessage\x12\x0b\n\x03num\x18\x01 \x01(\x03\x12\x0b\n\x03str\x18\x02 \x01(\t\"!\n\x11RepeatedContainer\x12\x0c\n\x04nums\x18\x01 \x03(\x03\"7\n\x0e\x45mptyContainer\x12%\n\x05\x65mpty\x18\x01 \x01(\x0b\x32\x16.google.protobuf.Empty*\x1b\n\x01\x45\x12\x06\n\x02\x45\x30\x10\x00\x12\x06\n\x02\x45\x31\x10\x01\x12\x06\n\x02\x45\x32\x10\x02\x62\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,google_dot_protobuf_dot_struct__pb2.DESCRIPTOR,google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

_E = _descriptor.EnumDescriptor(
  name='E',
  full_name='bigquery.E',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='E0', index=0, number=0,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='E1', index=1, number=1,
      options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='E2', index=2, number=2,
      options=None,
      type=None),
  ],
  containing_type=None,
  options=None,
  serialized_start=775,
  serialized_end=802,
)
_sym_db.RegisterEnumDescriptor(_E)

E = enum_type_wrapper.EnumTypeWrapper(_E)
E0 = 0
E1 = 1
E2 = 2



_TESTMESSAGE = _descriptor.Descriptor(
  name='TestMessage',
  full_name='bigquery.TestMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='str', full_name='bigquery.TestMessage.str', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='strs', full_name='bigquery.TestMessage.strs', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='num', full_name='bigquery.TestMessage.num', index=2,
      number=3, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='nums', full_name='bigquery.TestMessage.nums', index=3,
      number=4, type=3, cpp_type=2, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='e', full_name='bigquery.TestMessage.e', index=4,
      number=5, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='es', full_name='bigquery.TestMessage.es', index=5,
      number=6, type=14, cpp_type=8, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='nested', full_name='bigquery.TestMessage.nested', index=6,
      number=7, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='nesteds', full_name='bigquery.TestMessage.nesteds', index=7,
      number=8, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='empty', full_name='bigquery.TestMessage.empty', index=8,
      number=9, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='empties', full_name='bigquery.TestMessage.empties', index=9,
      number=10, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='struct', full_name='bigquery.TestMessage.struct', index=10,
      number=11, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='structs', full_name='bigquery.TestMessage.structs', index=11,
      number=12, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='bigquery.TestMessage.timestamp', index=12,
      number=13, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='timestamps', full_name='bigquery.TestMessage.timestamps', index=13,
      number=14, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='repeated_container', full_name='bigquery.TestMessage.repeated_container', index=14,
      number=15, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=124,
  serialized_end=638,
)


_NESTEDMESSAGE = _descriptor.Descriptor(
  name='NestedMessage',
  full_name='bigquery.NestedMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='num', full_name='bigquery.NestedMessage.num', index=0,
      number=1, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='str', full_name='bigquery.NestedMessage.str', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=640,
  serialized_end=681,
)


_REPEATEDCONTAINER = _descriptor.Descriptor(
  name='RepeatedContainer',
  full_name='bigquery.RepeatedContainer',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='nums', full_name='bigquery.RepeatedContainer.nums', index=0,
      number=1, type=3, cpp_type=2, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=683,
  serialized_end=716,
)


_EMPTYCONTAINER = _descriptor.Descriptor(
  name='EmptyContainer',
  full_name='bigquery.EmptyContainer',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='empty', full_name='bigquery.EmptyContainer.empty', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=718,
  serialized_end=773,
)

_TESTMESSAGE.fields_by_name['e'].enum_type = _E
_TESTMESSAGE.fields_by_name['es'].enum_type = _E
_TESTMESSAGE.fields_by_name['nested'].message_type = _NESTEDMESSAGE
_TESTMESSAGE.fields_by_name['nesteds'].message_type = _NESTEDMESSAGE
_TESTMESSAGE.fields_by_name['empty'].message_type = google_dot_protobuf_dot_empty__pb2._EMPTY
_TESTMESSAGE.fields_by_name['empties'].message_type = google_dot_protobuf_dot_empty__pb2._EMPTY
_TESTMESSAGE.fields_by_name['struct'].message_type = google_dot_protobuf_dot_struct__pb2._STRUCT
_TESTMESSAGE.fields_by_name['structs'].message_type = google_dot_protobuf_dot_struct__pb2._STRUCT
_TESTMESSAGE.fields_by_name['timestamp'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_TESTMESSAGE.fields_by_name['timestamps'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_TESTMESSAGE.fields_by_name['repeated_container'].message_type = _REPEATEDCONTAINER
_EMPTYCONTAINER.fields_by_name['empty'].message_type = google_dot_protobuf_dot_empty__pb2._EMPTY
DESCRIPTOR.message_types_by_name['TestMessage'] = _TESTMESSAGE
DESCRIPTOR.message_types_by_name['NestedMessage'] = _NESTEDMESSAGE
DESCRIPTOR.message_types_by_name['RepeatedContainer'] = _REPEATEDCONTAINER
DESCRIPTOR.message_types_by_name['EmptyContainer'] = _EMPTYCONTAINER
DESCRIPTOR.enum_types_by_name['E'] = _E

TestMessage = _reflection.GeneratedProtocolMessageType('TestMessage', (_message.Message,), dict(
  DESCRIPTOR = _TESTMESSAGE,
  __module__ = 'testmessage_pb2'
  # @@protoc_insertion_point(class_scope:bigquery.TestMessage)
  ))
_sym_db.RegisterMessage(TestMessage)

NestedMessage = _reflection.GeneratedProtocolMessageType('NestedMessage', (_message.Message,), dict(
  DESCRIPTOR = _NESTEDMESSAGE,
  __module__ = 'testmessage_pb2'
  # @@protoc_insertion_point(class_scope:bigquery.NestedMessage)
  ))
_sym_db.RegisterMessage(NestedMessage)

RepeatedContainer = _reflection.GeneratedProtocolMessageType('RepeatedContainer', (_message.Message,), dict(
  DESCRIPTOR = _REPEATEDCONTAINER,
  __module__ = 'testmessage_pb2'
  # @@protoc_insertion_point(class_scope:bigquery.RepeatedContainer)
  ))
_sym_db.RegisterMessage(RepeatedContainer)

EmptyContainer = _reflection.GeneratedProtocolMessageType('EmptyContainer', (_message.Message,), dict(
  DESCRIPTOR = _EMPTYCONTAINER,
  __module__ = 'testmessage_pb2'
  # @@protoc_insertion_point(class_scope:bigquery.EmptyContainer)
  ))
_sym_db.RegisterMessage(EmptyContainer)


# @@protoc_insertion_point(module_scope)