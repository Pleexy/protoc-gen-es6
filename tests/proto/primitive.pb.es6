// @flow
const jspb = require('google-protobuf');


export type PrimitiveTest$Object = {
  doubleField: number,  // double double_field = 1;
  floatField: number,  // float float_field = 2;
  int32Field: number,  // int32 int32_field = 3;
  int64Field: number,  // int64 int64_field = 4;
  uint32Field: number,  // uint32 uint32_field = 5;
  uint64Field: number,  // uint64 uint64_field = 6;
  sint32Field: number,  // sint32 sint32_field = 7;
  sint64Field: number,  // sint64 sint64_field = 8;
  fixed32Field: number,  // fixed32 fixed32_field = 9;
  fixed64Field: number,  // fixed64 fixed64_field = 10;
  sfixed32Field: number,  // sfixed32 sfixed32_field = 11;
  sfixed64Field: number,  // sfixed64 sfixed64_field = 12;
  boolField: boolean,  // bool bool_field = 13;
  stringField: string,  // string string_field = 14;
  bytesField: Uint8Array,  // bytes bytes_field = 15;
}

class PrimitiveTest {
  #doubleField: number = 0;  // double double_field = 1;
  #floatField: number = 0;  // float float_field = 2;
  #int32Field: number = 0;  // int32 int32_field = 3;
  #int64Field: number = 0;  // int64 int64_field = 4;
  #uint32Field: number = 0;  // uint32 uint32_field = 5;
  #uint64Field: number = 0;  // uint64 uint64_field = 6;
  #sint32Field: number = 0;  // sint32 sint32_field = 7;
  #sint64Field: number = 0;  // sint64 sint64_field = 8;
  #fixed32Field: number = 0;  // fixed32 fixed32_field = 9;
  #fixed64Field: number = 0;  // fixed64 fixed64_field = 10;
  #sfixed32Field: number = 0;  // sfixed32 sfixed32_field = 11;
  #sfixed64Field: number = 0;  // sfixed64 sfixed64_field = 12;
  #boolField: boolean = false;  // bool bool_field = 13;
  #stringField: string = '';  // string string_field = 14;
  #bytesField: Uint8Array;  // bytes bytes_field = 15;
  /**
  * optional double double_field = 1;
  * @return {number}
  */
  get doubleField():number {
    return this.#doubleField;
  }
  
  /** @param {number} val */
  set doubleField(val: number):void {
    if (typeof val === 'number') {
      this.#doubleField = val;
    } else {
      throw new Error('Expected type number for field #doubleField');
    }
  }
  /**
  * optional float float_field = 2;
  * @return {number}
  */
  get floatField():number {
    return this.#floatField;
  }
  
  /** @param {number} val */
  set floatField(val: number):void {
    if (typeof val === 'number') {
      this.#floatField = val;
    } else {
      throw new Error('Expected type number for field #floatField');
    }
  }
  /**
  * optional int32 int32_field = 3;
  * @return {number}
  */
  get int32Field():number {
    return this.#int32Field;
  }
  
  /** @param {number} val */
  set int32Field(val: number):void {
    if (typeof val === 'number') {
      this.#int32Field = val;
    } else {
      throw new Error('Expected type number for field #int32Field');
    }
  }
  /**
  * optional int64 int64_field = 4;
  * @return {number}
  */
  get int64Field():number {
    return this.#int64Field;
  }
  
  /** @param {number} val */
  set int64Field(val: number):void {
    if (typeof val === 'number') {
      this.#int64Field = val;
    } else {
      throw new Error('Expected type number for field #int64Field');
    }
  }
  /**
  * optional uint32 uint32_field = 5;
  * @return {number}
  */
  get uint32Field():number {
    return this.#uint32Field;
  }
  
  /** @param {number} val */
  set uint32Field(val: number):void {
    if (typeof val === 'number') {
      this.#uint32Field = val;
    } else {
      throw new Error('Expected type number for field #uint32Field');
    }
  }
  /**
  * optional uint64 uint64_field = 6;
  * @return {number}
  */
  get uint64Field():number {
    return this.#uint64Field;
  }
  
  /** @param {number} val */
  set uint64Field(val: number):void {
    if (typeof val === 'number') {
      this.#uint64Field = val;
    } else {
      throw new Error('Expected type number for field #uint64Field');
    }
  }
  /**
  * optional sint32 sint32_field = 7;
  * @return {number}
  */
  get sint32Field():number {
    return this.#sint32Field;
  }
  
  /** @param {number} val */
  set sint32Field(val: number):void {
    if (typeof val === 'number') {
      this.#sint32Field = val;
    } else {
      throw new Error('Expected type number for field #sint32Field');
    }
  }
  /**
  * optional sint64 sint64_field = 8;
  * @return {number}
  */
  get sint64Field():number {
    return this.#sint64Field;
  }
  
  /** @param {number} val */
  set sint64Field(val: number):void {
    if (typeof val === 'number') {
      this.#sint64Field = val;
    } else {
      throw new Error('Expected type number for field #sint64Field');
    }
  }
  /**
  * optional fixed32 fixed32_field = 9;
  * @return {number}
  */
  get fixed32Field():number {
    return this.#fixed32Field;
  }
  
  /** @param {number} val */
  set fixed32Field(val: number):void {
    if (typeof val === 'number') {
      this.#fixed32Field = val;
    } else {
      throw new Error('Expected type number for field #fixed32Field');
    }
  }
  /**
  * optional fixed64 fixed64_field = 10;
  * @return {number}
  */
  get fixed64Field():number {
    return this.#fixed64Field;
  }
  
  /** @param {number} val */
  set fixed64Field(val: number):void {
    if (typeof val === 'number') {
      this.#fixed64Field = val;
    } else {
      throw new Error('Expected type number for field #fixed64Field');
    }
  }
  /**
  * optional sfixed32 sfixed32_field = 11;
  * @return {number}
  */
  get sfixed32Field():number {
    return this.#sfixed32Field;
  }
  
  /** @param {number} val */
  set sfixed32Field(val: number):void {
    if (typeof val === 'number') {
      this.#sfixed32Field = val;
    } else {
      throw new Error('Expected type number for field #sfixed32Field');
    }
  }
  /**
  * optional sfixed64 sfixed64_field = 12;
  * @return {number}
  */
  get sfixed64Field():number {
    return this.#sfixed64Field;
  }
  
  /** @param {number} val */
  set sfixed64Field(val: number):void {
    if (typeof val === 'number') {
      this.#sfixed64Field = val;
    } else {
      throw new Error('Expected type number for field #sfixed64Field');
    }
  }
  /**
  * optional bool bool_field = 13;
  * @return {boolean}
  */
  get boolField():boolean {
    return this.#boolField;
  }
  
  /** @param {boolean} val */
  set boolField(val: boolean):void {
    if (typeof val === 'boolean' || (typeof val === 'object' &&  val !== null && typeof val.valueOf() === 'boolean')) {
      this.#boolField = val;
    } else {
      throw new Error('Expected type boolean for field #boolField');
    }
  }
  /**
  * optional string string_field = 14;
  * @return {string}
  */
  get stringField():string {
    return this.#stringField;
  }
  
  /** @param {string} val */
  set stringField(val: string):void {
    if (typeof val === 'string' || val instanceof String) {
      this.#stringField = val;
    } else {
      throw new Error('Expected type string for field #stringField');
    }
  }
  /**
  * optional bytes bytes_field = 15;
  * @return {Uint8Array}
  */
  get bytesField():Uint8Array {
    return this.#bytesField;
  }
  
  /** @param {Uint8Array} val */
  set bytesField(val: Uint8Array):void {
    if (val instanceof Uint8Array) {
      this.#bytesField = val;
    } else {
      throw new Error('Expected type Uint8Array for field #bytesField');
    }
  }

  /**
  * Deserializes binary data (in protobuf wire format).
  * @param {Uint8Array} bytes The bytes to deserialize.
  * @return {!PrimitiveTest}
  */
  static deserializeBinary(bytes: Uint8Array):PrimitiveTest {
    const reader = new jspb.BinaryReader(bytes);
    const msg = new PrimitiveTest();
    return PrimitiveTest.deserializeBinaryFromReader(msg, reader);
  }
  
  /**
  * Deserializes binary data (in protobuf wire format) from the
  * given reader into the given message object.
  * @param {!PrimitiveTest} msg The message object to deserialize into.
  * @param {!jspb.BinaryReader} reader The BinaryReader to use.
  * @return {!PrimitiveTest}
  */
  static deserializeBinaryFromReader (msg: PrimitiveTest, reader: BinaryReader) {
    while (reader.nextField()) {
      if (reader.isEndGroup()) {
        break;
      }
      let field = reader.getFieldNumber();
      switch (field) {
        case 1:
          msg.#doubleField = reader.readDouble();
          break;
        case 2:
          msg.#floatField = reader.readFloat();
          break;
        case 3:
          msg.#int32Field = reader.readInt32();
          break;
        case 4:
          msg.#int64Field = reader.readInt64();
          break;
        case 5:
          msg.#uint32Field = reader.readUint32();
          break;
        case 6:
          msg.#uint64Field = reader.readUint64();
          break;
        case 7:
          msg.#sint32Field = reader.readSint32();
          break;
        case 8:
          msg.#sint64Field = reader.readSint64();
          break;
        case 9:
          msg.#fixed32Field = reader.readFixed32();
          break;
        case 10:
          msg.#fixed64Field = reader.readFixed64();
          break;
        case 11:
          msg.#sfixed32Field = reader.readSfixed32();
          break;
        case 12:
          msg.#sfixed64Field = reader.readSfixed64();
          break;
        case 13:
          msg.#boolField = reader.readBool();
          break;
        case 14:
          msg.#stringField = reader.readString();
          break;
        case 15:
          msg.#bytesField = reader.readBytes();
          break;
        default:
          reader.skipField();
          break;
        }
      }
      return msg;
    };

  /**
  * Serializes the message to binary data (in protobuf wire format).
  * @return {!Uint8Array}
  */
  serializeBinary(): Uint8Array {
    const writer = new jspb.BinaryWriter();
    PrimitiveTest.serializeBinaryToWriter(this, writer);
    return writer.getResultBuffer();
  }
  
  /**
  * Serializes the given message to binary data (in protobuf wire
  * format), writing to the given BinaryWriter.
  * @param {!PrimitiveTest} msg
  * @param {!jspb.BinaryWriter} writer
  * @suppress {unusedLocalVariables} f is only used for nested messages
  */
  static serializeBinaryToWriter(msg: PrimitiveTest, writer: BinaryWriter): Uint8Array {
    if (msg.#doubleField !== 0){
      writer.writeDouble(1, msg.#doubleField); 
    }
    if (msg.#floatField !== 0){
      writer.writeFloat(2, msg.#floatField); 
    }
    if (msg.#int32Field !== 0){
      writer.writeInt32(3, msg.#int32Field); 
    }
    if (msg.#int64Field !== 0){
      writer.writeInt64(4, msg.#int64Field); 
    }
    if (msg.#uint32Field !== 0){
      writer.writeUint32(5, msg.#uint32Field); 
    }
    if (msg.#uint64Field !== 0){
      writer.writeUint64(6, msg.#uint64Field); 
    }
    if (msg.#sint32Field !== 0){
      writer.writeSint32(7, msg.#sint32Field); 
    }
    if (msg.#sint64Field !== 0){
      writer.writeSint64(8, msg.#sint64Field); 
    }
    if (msg.#fixed32Field !== 0){
      writer.writeFixed32(9, msg.#fixed32Field); 
    }
    if (msg.#fixed64Field !== 0){
      writer.writeFixed64(10, msg.#fixed64Field); 
    }
    if (msg.#sfixed32Field !== 0){
      writer.writeSfixed32(11, msg.#sfixed32Field); 
    }
    if (msg.#sfixed64Field !== 0){
      writer.writeSfixed64(12, msg.#sfixed64Field); 
    }
    if (msg.#boolField){
      writer.writeBool(13, msg.#boolField); 
    }
    if (msg.#stringField.length > 0){
      writer.writeString(14, msg.#stringField); 
    }
    if (msg.#bytesField && msg.#bytesField.length > 0){
      writer.writeBytes(15, msg.#bytesField); 
    }
  }
}
module.exports.PrimitiveTest = PrimitiveTest;
