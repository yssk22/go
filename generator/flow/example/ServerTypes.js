import * as types from 'types';

export type Example = {
  BoolVal: boolean,
  BoolPtr: ?boolean,
  IntVaL: number,
  IntPtr: ?number,
  FloatVal: number,
  FloatPtr: ?number,
  StringVal: string,
  StringPtr: ?string,
  TimeVal: Date,
  TimePtr: ?Date,
  InnerVal: ?Inner,
  InnerPtr: ?Inner,
  Imported: types.RGB
};
export type Inner = {
  BoolVal: boolean,
  BoolPtr: ?boolean,
  IntVaL: number,
  IntPtr: ?number,
  FloatVal: number,
  FloatPtr: ?number,
  StringVal: string,
  StringPtr: ?string,
  TimeVal: Date,
  TimePtr: ?Date
};
