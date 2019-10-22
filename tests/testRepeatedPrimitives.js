const primitive = require('./proto/repeatedPrimitives.pb')
const googPrimitive = require('./proto/google/tests/proto/repeatedPrimitives_pb')
const expect = require('chai').expect;

const compareFields = (obj1, obj2, fields, obj1Suffix) => {
  for (let fld of fields) {
    expect(obj1[fld + obj1Suffix], fld).to.eql(obj2[fld])
  }
  return '';
};

const compareRepeatedPrimitiveProtos = (p1, p2, p1Suffix) => {
  compareFields(p1, p2, [
    "doubleField",/*"floatField"*/,"int32Field",
    "int64Field","uint32Field","uint64Field",
    "sint32Field","sint64Field","fixed32Field",
    "fixed64Field","sfixed32Field","sfixed64Field",
    "boolField","stringField","bytesField"], p1Suffix);
};

function genArray(elFunc, number) {
  let res = [];
  for (let i=0; i < number; i++) {
    let el = elFunc(i)
    res.push(el);
  }
  return res;
}

function setField(obj, field, value, isGoogleObj, isBytesField) {
  if (isGoogleObj) {
    let settetName = "set"+field[0].toUpperCase() + field.slice(1)+"List";
    obj[settetName].call(obj, value);
  } else {
    obj[field] = value;
  }
}

function generateObject(obj, isGoogleObj) {
  let i = 0;
  [
    "doubleField","floatField","int32Field",
    "int64Field","uint32Field","uint64Field",
    "sint32Field","sint64Field","fixed32Field",
    "fixed64Field","sfixed32Field","sfixed64Field",
    "boolField","stringField", "bytesField"].forEach(field => {
      let arr ;
      if (field === "boolField") {
        arr = genArray(i => i % 2 === 0, 10)
      } else if (field === "stringField") {
        arr = genArray(i => "test string" + i.toString(), 20)
      } else if (field === "bytesField") {
        const enc = new TextEncoder("utf-8");
        arr = genArray(i => enc.encode("test bytes string" + i.toString()), 30)
      } else  {
        arr = genArray(i => {
          return i * 2 + ((field === "doubleField" || field === "floatField") ? i / 10 : 0)
        }, 40 );
      }
      setField(obj, field, arr, isGoogleObj, field === "bytesField")
  })
}

describe('TestRepeatedPrimitives', () => {
  it('serialize and deserialize', () => {
    const pt = new primitive.RepeatedPrimitiveTest();
    generateObject(pt, false)

    const ser = pt.serializeBinary();

    const newPt = primitive.RepeatedPrimitiveTest.deserializeBinary(ser);
    compareRepeatedPrimitiveProtos(pt, newPt, "");
  });
  it('it compatible with default implementation', () => {
    const pt = new googPrimitive.RepeatedPrimitiveTest();
    generateObject(pt,true)
    const ptObj = pt.toObject();
    ptObj.bytesFieldList = pt.getBytesFieldList_asU8();

    const ser = pt.serializeBinary();
    const newPt = primitive.RepeatedPrimitiveTest.deserializeBinary(ser);
    compareRepeatedPrimitiveProtos(ptObj, newPt, "List");

    const reverseSer = newPt.serializeBinary();
    const googNewPt = googPrimitive.RepeatedPrimitiveTest.deserializeBinary(reverseSer);
    const googNewPtObj = googNewPt.toObject();
    googNewPtObj.bytesFieldList = googNewPt.getBytesFieldList_asU8();
    compareRepeatedPrimitiveProtos(googNewPtObj, newPt, "List");
  });
});
