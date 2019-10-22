const primitive = require('./proto/primitive.pb')
const googPrimitive = require('./proto/google/tests/proto/primitive_pb')
const expect = require('chai').expect;

const compareFields = (obj1, obj2, fields) => {
  for (let fld of fields) {
    expect(obj1[fld], fld).to.eql(obj2[fld])
  }
  return '';
};

const comparePrimitiveProtos = (p1, p2) => {
  compareFields(p1, p2, [
    "doubleField",/*"floatField"*/,"int32Field",
    "int64Field","uint32Field","uint64Field",
    "sint32Field","sint64Field","fixed32Field",
    "fixed64Field","sfixed32Field","sfixed64Field",
    "boolField","stringField","bytesField"]);
  //expect(p1.bytesField).to.be.eql(p2.bytesField);
};

describe('TestPrimitives', () => {
  it('serialize and deserialize', () => {
    const pt = new primitive.PrimitiveTest();
    pt.doubleField = 1.1;  // double double_field = 1;
    pt.floatField = 2.2;  // float float_field = 2;
    pt.int32Field = 3;  // int32 int32_field = 3;
    pt.int64Field = 4;  // int64 int64_field = 4;
    pt.uint32Field = 5;  // uint32 uint32_field = 5;
    pt.uint64Field = 6;  // uint64 uint64_field = 6;
    pt.sint32Field = 7;  // sint32 sint32_field = 7;
    pt.sint64Field = 8;  // sint64 sint64_field = 8;
    pt.fixed32Field = 9;  // fixed32 fixed32_field = 9;
    pt.fixed64Field = 10;  // fixed64 fixed64_field = 10;
    pt.sfixed32Field = 11;  // sfixed32 sfixed32_field = 11;
    pt.sfixed64Field = 12;  // sfixed64 sfixed64_field = 12;
    pt.boolField = true;  // bool bool_field = 13;
    pt.stringField = 'test string';  // string string_field = 14;
    pt.bytesField = new TextEncoder("utf-8").encode("test buffer");  // bytes bytes_field = 15;

    const ser = pt.serializeBinary();

    const newPt = primitive.PrimitiveTest.deserializeBinary(ser);
    comparePrimitiveProtos(pt, newPt);
  });
  it('it compatible with default implementation', () => {
    const pt = new googPrimitive.PrimitiveTest();
    pt.setDoubleField(1.1);  // double double_field = 1;
    pt.setFloatField(2.2);  // float float_field = 2;
    pt.setInt32Field(3);  // int32 int32_field = 3;
    pt.setInt64Field(4);  // int64 int64_field = 4;
    pt.setUint32Field(5);  // uint32 uint32_field = 5;
    pt.setUint64Field(6);  // uint64 uint64_field = 6;
    pt.setSint32Field(7);  // sint32 sint32_field = 7;
    pt.setSint64Field(8);  // sint64 sint64_field = 8;
    pt.setFixed32Field(9);  // fixed32 fixed32_field = 9;
    pt.setFixed64Field(10);  // fixed64 fixed64_field = 10;
    pt.setSfixed32Field(11);  // sfixed32 sfixed32_field = 11;
    pt.setSfixed64Field(12);  // sfixed64 sfixed64_field = 12;
    pt.setBoolField(true);  // bool bool_field = 13;
    pt.setStringField('test string');  // string string_field = 14;
    pt.setBytesField(new TextEncoder("utf-8").encode("test buffer"));  // bytes bytes_field = 15;
    const ptObj = pt.toObject();
    ptObj.bytesField = pt.getBytesField_asU8();

    const ser = pt.serializeBinary();
    const newPt = primitive.PrimitiveTest.deserializeBinary(ser);
    comparePrimitiveProtos(ptObj, newPt);

    const reverseSer = newPt.serializeBinary();
    const googNewPt = googPrimitive.PrimitiveTest.deserializeBinary(reverseSer);
    const googNewPtObj = googNewPt.toObject();
    googNewPtObj.bytesField = googNewPt.getBytesField_asU8();
    comparePrimitiveProtos(googNewPtObj, newPt);
  });
});
